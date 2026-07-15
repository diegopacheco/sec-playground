package verification

import (
	"sync"
	"time"

	"agent-intent-bound-tokens/internal/intent"
	"agent-intent-bound-tokens/internal/token"
)

type AuditRecord struct {
	Time     time.Time `json:"time"`
	TokenID  string    `json:"token_id,omitempty"`
	Subject  string    `json:"subject,omitempty"`
	Action   string    `json:"action"`
	Resource string    `json:"resource"`
	Allowed  bool      `json:"allowed"`
	Code     string    `json:"code"`
	Reason   string    `json:"reason"`
}

type Service struct {
	tokens *token.Service
	mu     sync.Mutex
	used   map[string]struct{}
	audit  []AuditRecord
}

func New(tokens *token.Service) *Service {
	return &Service{tokens: tokens, used: make(map[string]struct{}), audit: make([]AuditRecord, 0)}
}

func (s *Service) Verify(request intent.Request) intent.Decision {
	claims, err := s.tokens.Parse(request.Token)
	if err != nil {
		return s.record(request, nil, intent.Decision{Code: "invalid_token", Reason: err.Error()})
	}
	decision := evaluate(request, claims)
	if !decision.Allowed {
		return s.record(request, &claims, decision)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.used[claims.ID]; exists {
		decision = intent.Decision{Code: "already_used", Reason: "intent token has already completed its authorized action", Claims: &claims}
		s.appendAudit(request, &claims, decision)
		return decision
	}
	s.used[claims.ID] = struct{}{}
	s.appendAudit(request, &claims, decision)
	return decision
}

func (s *Service) Audit() []AuditRecord {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]AuditRecord, len(s.audit))
	copy(result, s.audit)
	return result
}

func evaluate(request intent.Request, claims intent.Claims) intent.Decision {
	if request.AmountCents < 0 {
		return intent.Decision{Code: "invalid_amount", Reason: "requested amount cannot be negative", Claims: &claims}
	}
	if request.Audience != claims.Audience {
		return intent.Decision{Code: "audience_mismatch", Reason: "request audience is outside the signed intent", Claims: &claims}
	}
	if request.Action != claims.Action {
		return intent.Decision{Code: "action_mismatch", Reason: "request action is outside the signed intent", Claims: &claims}
	}
	if request.Resource != claims.Resource {
		return intent.Decision{Code: "resource_mismatch", Reason: "request resource is outside the signed intent", Claims: &claims}
	}
	if request.AmountCents > claims.MaxAmountCents {
		return intent.Decision{Code: "amount_exceeded", Reason: "request amount exceeds the signed limit", Claims: &claims}
	}
	return intent.Decision{Allowed: true, Code: "allowed", Reason: "request matches every signed intent constraint", Claims: &claims}
}

func (s *Service) record(request intent.Request, claims *intent.Claims, decision intent.Decision) intent.Decision {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.appendAudit(request, claims, decision)
	return decision
}

func (s *Service) appendAudit(request intent.Request, claims *intent.Claims, decision intent.Decision) {
	record := AuditRecord{Time: time.Now().UTC(), Action: request.Action, Resource: request.Resource, Allowed: decision.Allowed, Code: decision.Code, Reason: decision.Reason}
	if claims != nil {
		record.TokenID = claims.ID
		record.Subject = claims.Subject
	}
	s.audit = append(s.audit, record)
}
