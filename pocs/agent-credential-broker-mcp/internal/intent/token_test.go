package intent

import (
	"errors"
	"strings"
	"testing"
	"time"
)

type clock struct{ t time.Time }

func (c *clock) now() time.Time { return c.t }

func mint(s *Signer, c *clock) string {
	return s.Mint(Claims{ID: RandomID("jti_"), IntentID: "int_1", Credential: "vendor-api", Method: "GET", URL: "http://h/v1/profile", BodyHash: BodyHash(""), ExpiresAt: c.t.Add(time.Minute).Unix()})
}

func TestApprovedTokenWorksExactlyOnceSoAnAgentCannotRepeatAWrite(t *testing.T) {
	c := &clock{time.Unix(1000, 0)}
	s := NewSigner(c.now)
	tok := mint(s, c)
	if _, err := s.Redeem(tok, "GET", "http://h/v1/profile", ""); err != nil {
		t.Fatalf("first redemption must succeed: %v", err)
	}
	if _, err := s.Redeem(tok, "GET", "http://h/v1/profile", ""); !errors.Is(err, ErrReplay) {
		t.Fatalf("second redemption must be a replay, got %v", err)
	}
}

func TestTokenIsBoundToTheApprovedRequestSoApprovalForAReadCannotBecomeAWrite(t *testing.T) {
	cases := []struct{ method, url, body string }{
		{"POST", "http://h/v1/profile", ""},
		{"GET", "http://h/v1/tickets", ""},
		{"GET", "http://h/v1/profile", "{}"},
	}
	for _, tc := range cases {
		c := &clock{time.Unix(1000, 0)}
		s := NewSigner(c.now)
		if _, err := s.Redeem(mint(s, c), tc.method, tc.url, tc.body); !errors.Is(err, ErrMismatch) {
			t.Fatalf("%v must be rejected as mismatch, got %v", tc, err)
		}
	}
}

func TestFailedRedemptionBurnsTheTokenSoAnAgentCannotProbeVariations(t *testing.T) {
	c := &clock{time.Unix(1000, 0)}
	s := NewSigner(c.now)
	tok := mint(s, c)
	s.Redeem(tok, "POST", "http://h/v1/profile", "")
	if _, err := s.Redeem(tok, "GET", "http://h/v1/profile", ""); !errors.Is(err, ErrReplay) {
		t.Fatalf("token must be burned after a mismatched attempt, got %v", err)
	}
}

func TestExpiredTokenIsRejectedSoStaleApprovalsDoNotLinger(t *testing.T) {
	c := &clock{time.Unix(1000, 0)}
	s := NewSigner(c.now)
	tok := mint(s, c)
	c.t = c.t.Add(2 * time.Minute)
	if _, err := s.Redeem(tok, "GET", "http://h/v1/profile", ""); !errors.Is(err, ErrExpired) {
		t.Fatalf("expected expiry, got %v", err)
	}
}

func TestForgedOrForeignTokensAreRejectedBecauseOnlyTheBrokerHoldsTheKey(t *testing.T) {
	c := &clock{time.Unix(1000, 0)}
	s := NewSigner(c.now)
	other := NewSigner(c.now)
	if _, err := s.Redeem(mint(other, c), "GET", "http://h/v1/profile", ""); !errors.Is(err, ErrSignature) {
		t.Fatalf("token from another key must fail signature, got %v", err)
	}
	tok := mint(s, c)
	parts := strings.Split(tok, ".")
	forged := parts[0][:len(parts[0])-2] + "AA." + parts[1]
	if _, err := s.Redeem(forged, "GET", "http://h/v1/profile", ""); !errors.Is(err, ErrSignature) {
		t.Fatalf("edited payload must fail signature, got %v", err)
	}
	if _, err := s.Redeem("garbage", "GET", "http://h/v1/profile", ""); !errors.Is(err, ErrMalformed) {
		t.Fatalf("garbage must be malformed, got %v", err)
	}
}
