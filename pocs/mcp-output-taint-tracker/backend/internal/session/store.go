package session

import (
	"crypto/rand"
	"encoding/hex"
	"slices"
	"sync"
	"time"

	"mcp-output-taint-tracker/internal/policy"
	"mcp-output-taint-tracker/internal/taint"
)

type Decision struct {
	Allow   bool         `json:"allow"`
	Checked []string     `json:"checked"`
	Flows   []taint.Flow `json:"flows"`
}

type Event struct {
	Seq      int                   `json:"seq"`
	Kind     string                `json:"kind"`
	Tool     string                `json:"tool,omitempty"`
	Args     map[string]any        `json:"args,omitempty"`
	Text     string                `json:"text,omitempty"`
	Labels   []taint.Label         `json:"labels,omitempty"`
	Entities []taint.TaintedEntity `json:"entities,omitempty"`
	Decision *Decision             `json:"decision,omitempty"`
	At       time.Time             `json:"at"`
}

type Session struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	Created time.Time `json:"created"`
	Events  []Event   `json:"events"`
	tracker *taint.Tracker
}

type Summary struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	Created time.Time `json:"created"`
	Events  int       `json:"events"`
	Blocked int       `json:"blocked"`
}

type Store struct {
	mu       sync.Mutex
	policy   *policy.Policy
	sessions map[string]*Session
	order    []string
}

func NewStore(p *policy.Policy) *Store {
	return &Store{policy: p, sessions: map[string]*Session{}}
}

func (s *Store) Policy() *policy.Policy {
	return s.policy
}

func (s *Store) Create(name string) Session {
	s.mu.Lock()
	defer s.mu.Unlock()
	return snapshot(s.ensure(newID(), name))
}

func (s *Store) List() []Summary {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Summary, 0, len(s.order))
	for _, id := range slices.Backward(s.order) {
		ss := s.sessions[id]
		blocked := 0
		for _, e := range ss.Events {
			if e.Decision != nil && !e.Decision.Allow {
				blocked++
			}
		}
		out = append(out, Summary{ss.ID, ss.Name, ss.Created, len(ss.Events), blocked})
	}
	return out
}

func (s *Store) Get(id string) (Session, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	ss, ok := s.sessions[id]
	if !ok {
		return Session{}, false
	}
	return snapshot(ss), true
}

func (s *Store) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.sessions[id]; !ok {
		return false
	}
	delete(s.sessions, id)
	s.order = slices.DeleteFunc(s.order, func(x string) bool { return x == id })
	return true
}

func (s *Store) Prompt(id, text string) Event {
	s.mu.Lock()
	defer s.mu.Unlock()
	ss := s.ensure(id, "")
	ss.tracker.AddTrusted(text)
	return ss.add(Event{Kind: "prompt", Text: text})
}

func (s *Store) Result(id, tool string, args map[string]any, output string) Event {
	s.mu.Lock()
	defer s.mu.Unlock()
	ss := s.ensure(id, "")
	labels := s.policy.SourceLabels(tool, args)
	seq := len(ss.Events) + 1
	entities := ss.tracker.AddSource(seq, tool, labels, output)
	return ss.add(Event{Kind: "result", Tool: tool, Args: args, Text: output, Labels: labels, Entities: entities})
}

func (s *Store) Call(id, tool string, args map[string]any) Event {
	s.mu.Lock()
	defer s.mu.Unlock()
	ss := s.ensure(id, "")
	sinks := s.policy.SinkArgs(tool, args)
	d := &Decision{Allow: true, Checked: []string{}, Flows: ss.tracker.Check(sinks)}
	for _, a := range sinks {
		d.Checked = append(d.Checked, a.Arg)
	}
	d.Allow = len(d.Flows) == 0
	return ss.add(Event{Kind: "call", Tool: tool, Args: args, Decision: d})
}

func (s *Store) ensure(id, name string) *Session {
	if ss, ok := s.sessions[id]; ok {
		return ss
	}
	if name == "" {
		name = id
	}
	ss := &Session{ID: id, Name: name, Created: time.Now().UTC(), Events: []Event{}, tracker: taint.NewTracker()}
	s.sessions[id] = ss
	s.order = append(s.order, id)
	return ss
}

func (ss *Session) add(e Event) Event {
	e.Seq = len(ss.Events) + 1
	e.At = time.Now().UTC()
	ss.Events = append(ss.Events, e)
	return e
}

func snapshot(ss *Session) Session {
	return Session{ID: ss.ID, Name: ss.Name, Created: ss.Created, Events: slices.Clone(ss.Events)}
}

func newID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}
