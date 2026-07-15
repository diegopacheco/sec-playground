package main

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestMintAndVerify(t *testing.T) {
	service := newMacaroonService([]byte("root-key"))
	token, _, err := service.mint("/records/*", "read", "us-west", time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	result := service.verify(token, AccessRequest{Resource: "/records/payroll", Operation: "read", Location: "us-west"})
	if !result.Allowed {
		t.Fatalf("expected access to pass: %s", result.Reason)
	}
}

func TestAttenuationNarrowsAuthority(t *testing.T) {
	service := newMacaroonService([]byte("root-key"))
	token, _, err := service.mint("/records/*", "*", "*", time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	token, _, err = service.attenuate(token, Caveat{Type: "resource", Value: "/records/payroll"})
	if err != nil {
		t.Fatal(err)
	}
	token, _, err = service.attenuate(token, Caveat{Type: "operation", Value: "read"})
	if err != nil {
		t.Fatal(err)
	}
	allowed := service.verify(token, AccessRequest{Resource: "/records/payroll", Operation: "read", Location: "eu-central"})
	deniedResource := service.verify(token, AccessRequest{Resource: "/records/budget", Operation: "read", Location: "eu-central"})
	deniedOperation := service.verify(token, AccessRequest{Resource: "/records/payroll", Operation: "write", Location: "eu-central"})
	if !allowed.Allowed || deniedResource.Allowed || deniedOperation.Allowed {
		t.Fatal("attenuated authority was not enforced")
	}
}

func TestLaterWildcardCannotRemoveRestriction(t *testing.T) {
	service := newMacaroonService([]byte("root-key"))
	token, _, err := service.mint("/records/payroll", "read", "us-west", time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	token, _, err = service.attenuate(token, Caveat{Type: "resource", Value: "*"})
	if err != nil {
		t.Fatal(err)
	}
	result := service.verify(token, AccessRequest{Resource: "/records/budget", Operation: "read", Location: "us-west"})
	if result.Allowed {
		t.Fatal("a later wildcard expanded authority")
	}
}

func TestTamperingBreaksSignature(t *testing.T) {
	service := newMacaroonService([]byte("root-key"))
	token, macaroon, err := service.mint("/records/*", "read", "us-west", time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	macaroon.Caveats[1].Value = "write"
	data, err := json.Marshal(macaroon)
	if err != nil {
		t.Fatal(err)
	}
	tampered := base64.RawURLEncoding.EncodeToString(data)
	result := service.verify(tampered, AccessRequest{Resource: "/records/payroll", Operation: "write", Location: "us-west"})
	if result.Allowed || !strings.Contains(result.Reason, "signature") || token == tampered {
		t.Fatal("tampered token was not rejected")
	}
}

func TestExpiration(t *testing.T) {
	service := newMacaroonService([]byte("root-key"))
	current := time.Date(2026, 7, 14, 12, 0, 0, 0, time.UTC)
	service.now = func() time.Time { return current }
	token, _, err := service.mint("*", "*", "*", time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	service.now = func() time.Time { return current.Add(2 * time.Minute) }
	result := service.verify(token, AccessRequest{Resource: "/records/payroll", Operation: "read", Location: "us-west"})
	if result.Allowed {
		t.Fatal("expired token was accepted")
	}
}
