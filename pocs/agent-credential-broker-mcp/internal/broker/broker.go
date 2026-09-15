package broker

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"agent-credential-broker-mcp/internal/audit"
	"agent-credential-broker-mcp/internal/egress"
	"agent-credential-broker-mcp/internal/intent"
	"agent-credential-broker-mcp/internal/policy"
	"agent-credential-broker-mcp/internal/provider"
	"agent-credential-broker-mcp/internal/redact"
)

const (
	StatusPending  = "pending"
	StatusApproved = "approved"
	StatusDenied   = "denied"
	StatusUsed     = "used"
	StatusExpired  = "expired"
	maxBody        = 64 << 10
	maxResponse    = 1 << 20
)

var (
	ErrNotFound    = errors.New("intent not found")
	ErrNotPending  = errors.New("intent is not pending")
	ErrUnknownCred = errors.New("unknown credential")
	ErrBodyTooBig  = errors.New("request body too large")
)

type Intent struct {
	ID         string    `json:"id"`
	Credential string    `json:"credential"`
	Method     string    `json:"method"`
	URL        string    `json:"url"`
	Body       string    `json:"body"`
	BodyHash   string    `json:"body_sha256"`
	Reason     string    `json:"reason"`
	Approval   string    `json:"approval"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	DecidedAt  time.Time `json:"decided_at,omitzero"`
	ExpiresAt  time.Time `json:"expires_at,omitzero"`
	token      string
}

type AccessRequest struct {
	Credential string `json:"credential"`
	Method     string `json:"method"`
	URL        string `json:"url"`
	Body       string `json:"body"`
	Reason     string `json:"reason"`
}

type CallRequest struct {
	Token  string `json:"token"`
	Method string `json:"method"`
	URL    string `json:"url"`
	Body   string `json:"body"`
}

type CallResult struct {
	Status      int    `json:"status"`
	ContentType string `json:"content_type,omitempty"`
	Location    string `json:"location,omitempty"`
	Body        string `json:"body"`
	Redactions  int    `json:"redactions"`
}

type CredentialView struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Provider    string        `json:"provider"`
	Approval    string        `json:"approval"`
	Header      string        `json:"inject_header"`
	Allow       []policy.Rule `json:"allow"`
}

type Broker struct {
	policy     *policy.Policy
	source     provider.Source
	signer     *intent.Signer
	audit      *audit.Log
	client     *http.Client
	now        func() time.Time
	TokenTTL   time.Duration
	PendingTTL time.Duration
	mu         sync.Mutex
	intents    map[string]*Intent
}

func New(p *policy.Policy, source provider.Source, log *audit.Log, now func() time.Time) *Broker {
	return &Broker{
		policy:     p,
		source:     source,
		signer:     intent.NewSigner(now),
		audit:      log,
		now:        now,
		TokenTTL:   2 * time.Minute,
		PendingTTL: 10 * time.Minute,
		intents:    map[string]*Intent{},
		client: &http.Client{
			Timeout: 10 * time.Second,
			CheckRedirect: func(*http.Request, []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
}

func (b *Broker) Audit() *audit.Log { return b.audit }

func (b *Broker) Credentials() []CredentialView {
	out := []CredentialView{}
	for _, c := range b.policy.Credentials {
		out = append(out, CredentialView{c.Name, c.Description, c.Provider, c.Approval, c.Inject.Header, c.Allow})
	}
	return out
}

func (b *Broker) record(e audit.Entry) { b.audit.Append(e) }

func (b *Broker) RequestAccess(req AccessRequest) (Intent, error) {
	method := strings.ToUpper(req.Method)
	entry := audit.Entry{Event: "request", Credential: req.Credential, Method: method, URL: req.URL}
	c, ok := b.policy.Find(req.Credential)
	if !ok {
		entry.Outcome, entry.Detail = "denied", ErrUnknownCred.Error()
		b.record(entry)
		return Intent{}, ErrUnknownCred
	}
	if len(req.Body) > maxBody {
		entry.Outcome, entry.Detail = "denied", ErrBodyTooBig.Error()
		b.record(entry)
		return Intent{}, ErrBodyTooBig
	}
	u, err := egress.Check(c, method, req.URL)
	if err != nil {
		entry.Outcome, entry.Detail = "denied", err.Error()
		b.record(entry)
		return Intent{}, err
	}
	now := b.now()
	in := &Intent{
		ID:         intent.RandomID("int_"),
		Credential: c.Name,
		Method:     method,
		URL:        u.String(),
		Body:       req.Body,
		BodyHash:   intent.BodyHash(req.Body),
		Reason:     req.Reason,
		Approval:   c.Approval,
		Status:     StatusPending,
		CreatedAt:  now,
		ExpiresAt:  now.Add(b.PendingTTL),
	}
	b.mu.Lock()
	b.intents[in.ID] = in
	if c.Approval == "auto" {
		b.approveLocked(in)
	}
	view := *in
	b.mu.Unlock()
	entry.IntentID, entry.URL, entry.Outcome = in.ID, in.URL, in.Status
	entry.Detail = "approval=" + c.Approval
	b.record(entry)
	return view, nil
}

func (b *Broker) approveLocked(in *Intent) {
	now := b.now()
	in.Status = StatusApproved
	in.DecidedAt = now
	in.ExpiresAt = now.Add(b.TokenTTL)
	in.token = b.signer.Mint(intent.Claims{
		ID:         intent.RandomID("jti_"),
		IntentID:   in.ID,
		Credential: in.Credential,
		Method:     in.Method,
		URL:        in.URL,
		BodyHash:   in.BodyHash,
		ExpiresAt:  in.ExpiresAt.Unix(),
	})
}

func (b *Broker) decide(id string, approve bool) (Intent, error) {
	b.mu.Lock()
	in, ok := b.intents[id]
	if !ok {
		b.mu.Unlock()
		return Intent{}, ErrNotFound
	}
	b.refreshLocked(in)
	if in.Status != StatusPending {
		b.mu.Unlock()
		return *in, ErrNotPending
	}
	event := "deny"
	if approve {
		event = "approve"
		b.approveLocked(in)
	} else {
		in.Status = StatusDenied
		in.DecidedAt = b.now()
	}
	view := *in
	b.mu.Unlock()
	b.record(audit.Entry{Event: event, Outcome: view.Status, IntentID: id, Credential: view.Credential, Method: view.Method, URL: view.URL, Detail: "decided by human"})
	return view, nil
}

func (b *Broker) Approve(id string) (Intent, error) { return b.decide(id, true) }

func (b *Broker) Deny(id string) (Intent, error) { return b.decide(id, false) }

func (b *Broker) refreshLocked(in *Intent) {
	if (in.Status == StatusPending || in.Status == StatusApproved) && !b.now().Before(in.ExpiresAt) {
		in.Status = StatusExpired
		in.token = ""
	}
}

func (b *Broker) Check(id string) (Intent, string, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	in, ok := b.intents[id]
	if !ok {
		return Intent{}, "", ErrNotFound
	}
	b.refreshLocked(in)
	if in.Status != StatusApproved {
		return *in, "", nil
	}
	return *in, in.token, nil
}

func (b *Broker) Intents() []Intent {
	b.mu.Lock()
	defer b.mu.Unlock()
	out := make([]Intent, 0, len(b.intents))
	for _, in := range b.intents {
		b.refreshLocked(in)
		out = append(out, *in)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.After(out[j].CreatedAt) })
	return out
}

func (b *Broker) markUsed(id string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if in, ok := b.intents[id]; ok {
		in.Status = StatusUsed
		in.token = ""
	}
}

func (b *Broker) Call(ctx context.Context, req CallRequest) (CallResult, error) {
	method := strings.ToUpper(req.Method)
	entry := audit.Entry{Event: "call", Method: method, URL: req.URL}
	fail := func(err error) (CallResult, error) {
		entry.Outcome, entry.Detail = "denied", err.Error()
		b.record(entry)
		return CallResult{}, err
	}
	u, err := egress.Normalize(req.URL)
	if err != nil {
		return fail(err)
	}
	entry.URL = u.String()
	claims, err := b.signer.Redeem(req.Token, method, u.String(), req.Body)
	entry.IntentID, entry.Credential = claims.IntentID, claims.Credential
	if claims.IntentID != "" && err != nil && !errors.Is(err, intent.ErrSignature) {
		b.markUsed(claims.IntentID)
	}
	if err != nil {
		return fail(err)
	}
	b.markUsed(claims.IntentID)
	c, ok := b.policy.Find(claims.Credential)
	if !ok {
		return fail(ErrUnknownCred)
	}
	if _, err := egress.Check(c, method, u.String()); err != nil {
		return fail(err)
	}
	value, err := b.source.Fetch(ctx, c)
	if err != nil {
		entry.Outcome, entry.Detail = "error", err.Error()
		b.record(entry)
		return CallResult{}, err
	}
	res, err := b.send(ctx, c, method, u.String(), req.Body, value)
	if err != nil {
		entry.Outcome, entry.Detail = "error", scrub(err.Error(), value, c.Name)
		b.record(entry)
		return CallResult{}, errors.New(entry.Detail)
	}
	entry.Outcome = "allowed"
	entry.Detail = fmt.Sprintf("upstream status %d, %d redaction(s)", res.Status, res.Redactions)
	b.record(entry)
	return res, nil
}

func scrub(text, value, label string) string {
	out, _ := redact.Apply(text, value, label)
	return out
}

func (b *Broker) send(ctx context.Context, c policy.Credential, method, url, body, value string) (CallResult, error) {
	var reader io.Reader
	if body != "" {
		reader = strings.NewReader(body)
	}
	hreq, err := http.NewRequestWithContext(ctx, method, url, reader)
	if err != nil {
		return CallResult{}, err
	}
	if body != "" {
		hreq.Header.Set("Content-Type", "application/json")
	}
	hreq.Header.Set("User-Agent", "agent-credential-broker")
	hreq.Header.Set(c.Inject.Header, strings.ReplaceAll(c.Inject.Template, "{{credential}}", value))
	hres, err := b.client.Do(hreq)
	if err != nil {
		return CallResult{}, err
	}
	defer hres.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(hres.Body, maxResponse))
	if err != nil {
		return CallResult{}, err
	}
	text, n := redact.Apply(string(raw), value, c.Name)
	location, m := redact.Apply(hres.Header.Get("Location"), value, c.Name)
	return CallResult{
		Status:      hres.StatusCode,
		ContentType: hres.Header.Get("Content-Type"),
		Location:    location,
		Body:        text,
		Redactions:  n + m,
	}, nil
}
