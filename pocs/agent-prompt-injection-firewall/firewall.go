package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

type FirewallConfig struct {
	ReadRoot      string
	AllowedHosts  []string
	Canary        string
	ApprovalToken string
	LocalURL      string
	AuditPath     string
}

type Firewall struct {
	config  FirewallConfig
	hosts   map[string]bool
	client  *http.Client
	auditMu sync.Mutex
}

type Candidate struct {
	Tool   string `json:"tool"`
	Method string `json:"method,omitempty"`
	Target string `json:"target"`
	Body   string `json:"body,omitempty"`
}

type Decision struct {
	Allowed bool   `json:"allowed"`
	Tool    string `json:"tool"`
	Target  string `json:"target"`
	Reason  string `json:"reason"`
	Output  string `json:"output,omitempty"`
}

type AuditEntry struct {
	Time     time.Time `json:"time"`
	Source   string    `json:"source"`
	Decision Decision  `json:"decision"`
}

func NewFirewall(config FirewallConfig) *Firewall {
	hosts := map[string]bool{}
	for _, host := range config.AllowedHosts {
		hosts[host] = true
	}
	return &Firewall{
		config: config,
		hosts:  hosts,
		client: &http.Client{
			Timeout: 3 * time.Second,
			CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
				return errors.New("redirect blocked")
			},
		},
	}
}

func ParseCandidates(content string) ([]Candidate, error) {
	candidates := []Candidate{}
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "ACTION ") {
			continue
		}
		var candidate Candidate
		if err := json.Unmarshal([]byte(strings.TrimPrefix(line, "ACTION ")), &candidate); err != nil {
			return nil, fmt.Errorf("invalid action: %w", err)
		}
		candidates = append(candidates, candidate)
	}
	return candidates, scanner.Err()
}

func (f *Firewall) Execute(source string, candidate Candidate, approval string) Decision {
	candidate.Target = f.expand(candidate.Target)
	candidate.Body = f.expand(candidate.Body)
	decision := f.authorize(candidate, approval)
	if decision.Allowed {
		output, err := f.invoke(candidate)
		if err != nil {
			decision.Allowed = false
			decision.Reason = "tool execution failed"
			decision.Output = f.redact(err.Error())
		} else {
			decision.Output = f.redact(output)
		}
	}
	decision.Target = f.redact(decision.Target)
	f.record(AuditEntry{Time: time.Now().UTC(), Source: source, Decision: decision})
	return decision
}

func (f *Firewall) authorize(candidate Candidate, approval string) Decision {
	decision := Decision{Tool: candidate.Tool, Target: candidate.Target}
	if f.containsCanary(candidate.Target) || f.containsCanary(candidate.Body) {
		decision.Reason = "sensitive canary detected"
		return decision
	}
	switch candidate.Tool {
	case "read":
		if _, err := safeJoin(f.config.ReadRoot, candidate.Target); err != nil {
			decision.Reason = "read target outside capability root"
			return decision
		}
		decision.Allowed = true
		decision.Reason = "read capability allowed"
		return decision
	case "http":
		parsed, err := url.Parse(candidate.Target)
		if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || !f.hosts[parsed.Host] || parsed.User != nil {
			decision.Reason = "HTTP destination not allowed"
			return decision
		}
		method := strings.ToUpper(candidate.Method)
		if method == "" {
			method = http.MethodGet
		}
		if method != http.MethodGet && method != http.MethodHead && method != http.MethodPost {
			decision.Reason = "HTTP method outside capability set"
			return decision
		}
		if method == http.MethodPost && approval != f.config.ApprovalToken {
			decision.Reason = "explicit approval required"
			return decision
		}
		decision.Allowed = true
		if method == http.MethodPost {
			decision.Reason = "approved HTTP write allowed"
		} else {
			decision.Reason = "HTTP read allowed"
		}
		return decision
	default:
		decision.Reason = "tool outside capability set"
		return decision
	}
}

func (f *Firewall) invoke(candidate Candidate) (string, error) {
	switch candidate.Tool {
	case "read":
		path, err := safeJoin(f.config.ReadRoot, candidate.Target)
		if err != nil {
			return "", err
		}
		content, err := os.ReadFile(path)
		return string(content), err
	case "http":
		method := strings.ToUpper(candidate.Method)
		if method == "" {
			method = http.MethodGet
		}
		request, err := http.NewRequest(method, candidate.Target, strings.NewReader(candidate.Body))
		if err != nil {
			return "", err
		}
		request.Header.Set("Content-Type", "application/json")
		response, err := f.client.Do(request)
		if err != nil {
			return "", err
		}
		defer response.Body.Close()
		body, err := io.ReadAll(io.LimitReader(response.Body, 4096))
		if err != nil {
			return "", err
		}
		if response.StatusCode < 200 || response.StatusCode > 299 {
			return "", fmt.Errorf("HTTP status %d", response.StatusCode)
		}
		return string(body), nil
	default:
		return "", errors.New("unknown tool")
	}
}

func (f *Firewall) expand(value string) string {
	value = strings.ReplaceAll(value, "${CANARY}", f.config.Canary)
	return strings.ReplaceAll(value, "${LOCAL_URL}", f.config.LocalURL)
}

func (f *Firewall) containsCanary(value string) bool {
	return f.config.Canary != "" && strings.Contains(value, f.config.Canary)
}

func (f *Firewall) redact(value string) string {
	if f.config.Canary == "" {
		return value
	}
	return strings.ReplaceAll(value, f.config.Canary, "[REDACTED]")
}

func (f *Firewall) record(entry AuditEntry) {
	f.auditMu.Lock()
	defer f.auditMu.Unlock()
	file, err := os.OpenFile(f.config.AuditPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return
	}
	defer file.Close()
	json.NewEncoder(file).Encode(entry)
}

func (f *Firewall) AuditEntries() ([]AuditEntry, error) {
	f.auditMu.Lock()
	defer f.auditMu.Unlock()
	file, err := os.Open(f.config.AuditPath)
	if errors.Is(err, os.ErrNotExist) {
		return []AuditEntry{}, nil
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()
	entries := []AuditEntry{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var entry AuditEntry
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, scanner.Err()
}
