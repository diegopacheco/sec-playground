package store

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"mcp-rug-pull-proxy/internal/mcp"
	"mcp-rug-pull-proxy/internal/pin"
	"mcp-rug-pull-proxy/internal/scan"
)

const (
	StatusPending  = "pending"
	StatusBlocked  = "blocked"
	StatusApproved = "approved"
	StatusRejected = "rejected"

	DecisionAllow   = "allow"
	DecisionPending = "pending"
	DecisionBlocked = "blocked"

	AlertOpen      = "open"
	AlertAccepted  = "accepted"
	AlertRejected  = "rejected"
	AlertForgotten = "forgotten"
)

var (
	ErrNotFound     = errors.New("tool is not pinned")
	ErrHashMismatch = errors.New("hash does not match the definition under review, approve only what you reviewed")
)

type Pin struct {
	Server     string      `json:"server"`
	Tool       string      `json:"tool"`
	Hash       string      `json:"hash"`
	Status     string      `json:"status"`
	Definition mcp.Tool    `json:"definition"`
	Scan       scan.Report `json:"scan"`
	FirstSeen  time.Time   `json:"firstSeen"`
	ApprovedAt *time.Time  `json:"approvedAt,omitempty"`
	AlertID    string      `json:"alertId,omitempty"`
}

type Alert struct {
	ID         string      `json:"id"`
	Server     string      `json:"server"`
	Tool       string      `json:"tool"`
	State      string      `json:"state"`
	OldHash    string      `json:"oldHash"`
	NewHash    string      `json:"newHash"`
	Old        mcp.Tool    `json:"old"`
	New        mcp.Tool    `json:"new"`
	OldScan    scan.Report `json:"oldScan"`
	NewScan    scan.Report `json:"newScan"`
	DetectedAt time.Time   `json:"detectedAt"`
	UpdatedAt  time.Time   `json:"updatedAt"`
}

type Decision struct {
	Tool     string `json:"tool"`
	Hash     string `json:"hash"`
	Decision string `json:"decision"`
	Status   string `json:"status"`
	Severity string `json:"severity"`
	Reason   string `json:"reason"`
}

type lockFile struct {
	Version int      `json:"version"`
	Pins    []*Pin   `json:"pins"`
	Alerts  []*Alert `json:"alerts"`
}

type Store struct {
	mu     sync.Mutex
	dir    string
	lock   lockFile
	events []Event
	now    func() time.Time
}

func Open(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, err
	}
	s := &Store{dir: dir, lock: lockFile{Version: 1}, now: func() time.Time { return time.Now().UTC() }}
	raw, err := os.ReadFile(s.lockPath())
	if err == nil {
		if err := json.Unmarshal(raw, &s.lock); err != nil {
			return nil, fmt.Errorf("read lockfile: %w", err)
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	s.events, err = loadEvents(s.eventsPath())
	if err != nil {
		return nil, fmt.Errorf("read events: %w", err)
	}
	return s, nil
}

func (s *Store) lockPath() string {
	return filepath.Join(s.dir, "tools.lock.json")
}

func (s *Store) eventsPath() string {
	return filepath.Join(s.dir, "events.jsonl")
}

func (s *Store) Check(server string, tools []mcp.Tool) ([]Decision, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	decisions := make([]Decision, 0, len(tools))
	for _, t := range tools {
		decisions = append(decisions, s.decide(server, t))
	}
	return decisions, s.save()
}

func (s *Store) decide(server string, t mcp.Tool) Decision {
	hash := pin.Hash(t)
	p := s.find(server, t.Name)
	if p == nil {
		report := scan.Scan(t, s.toolNamesExcept(server))
		p = &Pin{Server: server, Tool: t.Name, Hash: hash, Status: statusFor(report), Definition: t, Scan: report, FirstSeen: s.now()}
		s.lock.Pins = append(s.lock.Pins, p)
		if p.Status == StatusBlocked {
			s.record("tool-blocked-by-scanner", server, t.Name, hash, fmt.Sprintf("new tool quarantined, %d findings, severity %s", len(report.Findings), report.Severity))
		} else {
			s.record("tool-discovered", server, t.Name, hash, fmt.Sprintf("new tool waiting for review, severity %s", report.Severity))
		}
		return decision(p, hash, DecisionFor(p.Status), "new tool, not exposed until approved")
	}
	switch p.Status {
	case StatusApproved:
		if hash == p.Hash {
			return decision(p, hash, DecisionAllow, "definition matches the approved pin")
		}
		d := decision(p, hash, DecisionBlocked, "definition changed after approval")
		d.Severity = s.rugPull(p, t, hash).Severity
		return d
	case StatusRejected:
		return decision(p, hash, DecisionBlocked, "tool was rejected by a reviewer")
	}
	if hash != p.Hash {
		report := scan.Scan(t, s.toolNamesExcept(server))
		p.Hash, p.Definition, p.Scan, p.Status = hash, t, report, statusFor(report)
		s.record("definition-changed-before-approval", server, t.Name, hash, fmt.Sprintf("tool changed while waiting for review, severity %s", report.Severity))
	}
	return decision(p, hash, DecisionFor(p.Status), "waiting for review")
}

func (s *Store) rugPull(p *Pin, t mcp.Tool, hash string) scan.Report {
	report := scan.Scan(t, s.toolNamesExcept(p.Server))
	if a := s.alert(p.AlertID); a != nil && a.State == AlertOpen {
		if a.NewHash != hash {
			a.NewHash, a.New, a.NewScan, a.UpdatedAt = hash, t, report, s.now()
			s.record("definition-changed-again", p.Server, p.Tool, hash, "open rug pull alert updated with the newest definition")
		}
		return report
	}
	a := &Alert{
		ID: newID(), Server: p.Server, Tool: p.Tool, State: AlertOpen,
		OldHash: p.Hash, NewHash: hash, Old: p.Definition, New: t,
		OldScan: p.Scan, NewScan: report, DetectedAt: s.now(), UpdatedAt: s.now(),
	}
	s.lock.Alerts = append(s.lock.Alerts, a)
	p.AlertID = a.ID
	s.record("rug-pull-detected", p.Server, p.Tool, hash, fmt.Sprintf("approved %s changed to %s, tool blocked, severity %s", pin.Short(p.Hash), pin.Short(hash), report.Severity))
	return report
}

func (s *Store) Approve(server, tool, hash string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	p := s.find(server, tool)
	if p == nil {
		return ErrNotFound
	}
	if a := s.alert(p.AlertID); a != nil && a.State == AlertOpen && a.NewHash == hash {
		now := s.now()
		p.Hash, p.Definition, p.Scan, p.Status, p.ApprovedAt, p.AlertID = hash, a.New, a.NewScan, StatusApproved, &now, ""
		a.State, a.UpdatedAt = AlertAccepted, now
		s.record("change-accepted", server, tool, hash, "reviewer accepted the changed definition")
		return s.save()
	}
	if hash != p.Hash {
		return ErrHashMismatch
	}
	if p.Status == StatusApproved {
		return nil
	}
	now := s.now()
	p.Status, p.ApprovedAt = StatusApproved, &now
	s.record("tool-approved", server, tool, hash, "reviewer pinned this exact definition")
	return s.save()
}

func (s *Store) Reject(server, tool string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	p := s.find(server, tool)
	if p == nil {
		return ErrNotFound
	}
	if a := s.alert(p.AlertID); a != nil && a.State == AlertOpen {
		a.State, a.UpdatedAt = AlertRejected, s.now()
	}
	p.Status, p.AlertID = StatusRejected, ""
	s.record("tool-rejected", server, tool, p.Hash, "reviewer rejected the tool, it stays blocked")
	return s.save()
}

func (s *Store) Forget(server string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	kept := s.lock.Pins[:0]
	removed := 0
	for _, p := range s.lock.Pins {
		if p.Server == server {
			removed++
			continue
		}
		kept = append(kept, p)
	}
	s.lock.Pins = kept
	for _, a := range s.lock.Alerts {
		if a.Server == server && a.State == AlertOpen {
			a.State, a.UpdatedAt = AlertForgotten, s.now()
		}
	}
	s.record("server-forgotten", server, "", "", fmt.Sprintf("%d pins removed, tools will be reviewed again", removed))
	return s.save()
}

func (s *Store) RecordDenied(server, tool, reason string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.record("call-denied", server, tool, "", reason)
	return nil
}

func (s *Store) Pins() []Pin {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Pin, 0, len(s.lock.Pins))
	for _, p := range s.lock.Pins {
		out = append(out, *p)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Server != out[j].Server {
			return out[i].Server < out[j].Server
		}
		return out[i].Tool < out[j].Tool
	})
	return out
}

func (s *Store) Alerts() []Alert {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Alert, 0, len(s.lock.Alerts))
	for i := len(s.lock.Alerts) - 1; i >= 0; i-- {
		out = append(out, *s.lock.Alerts[i])
	}
	return out
}

func (s *Store) Events() ([]Event, ChainStatus) {
	s.mu.Lock()
	defer s.mu.Unlock()
	events, err := loadEvents(s.eventsPath())
	if err != nil {
		return s.events, ChainStatus{Count: len(s.events), Reason: "event log unreadable: " + err.Error()}
	}
	return events, Verify(events)
}

func (s *Store) record(kind, server, tool, hash, detail string) {
	prev := genesis
	if n := len(s.events); n > 0 {
		prev = s.events[n-1].Sum
	}
	e := Event{Seq: len(s.events) + 1, Time: s.now(), Type: kind, Server: server, Tool: tool, Hash: hash, Detail: detail, Prev: prev}
	e.Sum = seal(e)
	s.events = append(s.events, e)
	if err := appendEvent(s.eventsPath(), e); err != nil {
		fmt.Fprintln(os.Stderr, "append event:", err)
	}
}

func (s *Store) save() error {
	raw, err := json.MarshalIndent(s.lock, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.lockPath() + ".tmp"
	if err := os.WriteFile(tmp, raw, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, s.lockPath())
}

func (s *Store) find(server, tool string) *Pin {
	for _, p := range s.lock.Pins {
		if p.Server == server && p.Tool == tool {
			return p
		}
	}
	return nil
}

func (s *Store) alert(id string) *Alert {
	if id == "" {
		return nil
	}
	for _, a := range s.lock.Alerts {
		if a.ID == id {
			return a
		}
	}
	return nil
}

func (s *Store) ToolNamesExcept(server, tool string) []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	var names []string
	for _, n := range s.toolNamesExcept(server) {
		if n != tool {
			names = append(names, n)
		}
	}
	return names
}

func (s *Store) toolNamesExcept(server string) []string {
	var names []string
	for _, p := range s.lock.Pins {
		if p.Server != server {
			names = append(names, p.Tool)
		}
	}
	return names
}

func statusFor(report scan.Report) string {
	if report.Severity == scan.SeverityHigh {
		return StatusBlocked
	}
	return StatusPending
}

func DecisionFor(status string) string {
	switch status {
	case StatusApproved:
		return DecisionAllow
	case StatusPending:
		return DecisionPending
	}
	return DecisionBlocked
}

func decision(p *Pin, hash, kind, reason string) Decision {
	return Decision{Tool: p.Tool, Hash: hash, Decision: kind, Status: p.Status, Severity: p.Scan.Severity, Reason: reason}
}

func newID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
