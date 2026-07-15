package verification

import (
	"strings"
	"testing"
	"time"

	"agent-intent-bound-tokens/internal/intent"
	"agent-intent-bound-tokens/internal/token"
)

func TestIntentConstraintsAndOneTimeUse(t *testing.T) {
	tokens, err := token.New(strings.Repeat("s", 32))
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	value, _, err := tokens.Issue(intent.Claims{
		Issuer: "intent.local", Subject: "agent:buyer", Audience: "payments-api",
		Action: "transfer", Resource: "account:operations", MaxAmountCents: 50_000,
		NotBefore: now.Unix(), ExpiresAt: now.Add(time.Minute).Unix(),
	})
	if err != nil {
		t.Fatal(err)
	}
	service := New(tokens)
	base := intent.Request{Token: value, Audience: "payments-api", Action: "transfer", Resource: "account:operations", AmountCents: 49_999}

	wrongAction := base
	wrongAction.Action = "delete"
	if decision := service.Verify(wrongAction); decision.Allowed || decision.Code != "action_mismatch" {
		t.Fatalf("unexpected action decision: %+v", decision)
	}

	overLimit := base
	overLimit.AmountCents = 50_001
	if decision := service.Verify(overLimit); decision.Allowed || decision.Code != "amount_exceeded" {
		t.Fatalf("unexpected amount decision: %+v", decision)
	}

	if decision := service.Verify(base); !decision.Allowed {
		t.Fatalf("expected allowed decision: %+v", decision)
	}
	if decision := service.Verify(base); decision.Allowed || decision.Code != "already_used" {
		t.Fatalf("unexpected reuse decision: %+v", decision)
	}
	if len(service.Audit()) != 4 {
		t.Fatalf("expected four audit records")
	}
}
