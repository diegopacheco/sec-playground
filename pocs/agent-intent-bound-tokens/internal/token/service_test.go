package token

import (
	"errors"
	"strings"
	"testing"
	"time"

	"agent-intent-bound-tokens/internal/intent"
)

func TestIssueAndParse(t *testing.T) {
	service, err := New(strings.Repeat("s", 32))
	if err != nil {
		t.Fatal(err)
	}
	now := time.Unix(2_000_000_000, 0)
	service.now = func() time.Time { return now }
	issued, expected, err := service.Issue(intent.Claims{
		Issuer: "intent.local", Subject: "agent:buyer", Audience: "payments-api",
		Action: "transfer", Resource: "account:operations", MaxAmountCents: 50_000,
		NotBefore: now.Unix(), ExpiresAt: now.Add(time.Minute).Unix(),
	})
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := service.Parse(issued)
	if err != nil {
		t.Fatal(err)
	}
	if parsed != expected || parsed.ID == "" {
		t.Fatalf("unexpected claims: %+v", parsed)
	}
}

func TestTamperedTokenFails(t *testing.T) {
	service, _ := New(strings.Repeat("s", 32))
	now := time.Now()
	issued, _, err := service.Issue(intent.Claims{
		Issuer: "intent.local", Subject: "agent:buyer", Audience: "payments-api",
		Action: "transfer", Resource: "account:operations", MaxAmountCents: 50_000,
		NotBefore: now.Unix(), ExpiresAt: now.Add(time.Minute).Unix(),
	})
	if err != nil {
		t.Fatal(err)
	}
	parts := strings.Split(issued, ".")
	parts[1] = parts[1] + "A"
	if _, err := service.Parse(strings.Join(parts, ".")); !errors.Is(err, ErrSignature) {
		t.Fatalf("expected signature error, received %v", err)
	}
}

func TestExpiredTokenFails(t *testing.T) {
	service, _ := New(strings.Repeat("s", 32))
	now := time.Unix(2_000_000_000, 0)
	service.now = func() time.Time { return now }
	issued, _, err := service.Issue(intent.Claims{
		Issuer: "intent.local", Subject: "agent:buyer", Audience: "payments-api",
		Action: "transfer", Resource: "account:operations", MaxAmountCents: 50_000,
		NotBefore: now.Add(-time.Minute).Unix(), ExpiresAt: now.Add(-time.Second).Unix(),
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.Parse(issued); !errors.Is(err, ErrExpired) {
		t.Fatalf("expected expiration error, received %v", err)
	}
}
