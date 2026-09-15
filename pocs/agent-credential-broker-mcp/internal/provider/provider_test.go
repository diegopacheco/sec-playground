package provider

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"agent-credential-broker-mcp/internal/policy"
)

func TestVaultReadsKVv2FieldWithTheBrokerToken(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Vault-Token") != "root-token" || r.URL.Path != "/v1/secret/data/agent/vendor-api" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.Write([]byte(`{"data":{"data":{"api_key":"sk_from_vault"}}}`))
	}))
	defer srv.Close()
	v := Vault{Addr: srv.URL, Token: "root-token"}
	got, err := v.Fetch(context.Background(), policy.Credential{Path: "agent/vendor-api", Field: "api_key"})
	if err != nil || got != "sk_from_vault" {
		t.Fatalf("got %q %v", got, err)
	}
	if _, err := v.Fetch(context.Background(), policy.Credential{Path: "agent/vendor-api", Field: "missing"}); err == nil {
		t.Fatal("missing field must fail instead of injecting an empty credential")
	}
	if _, err := (Vault{Addr: srv.URL, Token: "wrong"}).Fetch(context.Background(), policy.Credential{Path: "agent/vendor-api", Field: "api_key"}); err == nil {
		t.Fatal("wrong vault token must fail")
	}
}

func TestKeychainUsesSecurityCliWithServiceAndAccount(t *testing.T) {
	var got []string
	k := Keychain{Run: func(ctx context.Context, name string, args ...string) ([]byte, error) {
		got = append([]string{name}, args...)
		return []byte("ghp_value\n"), nil
	}}
	value, err := k.Fetch(context.Background(), policy.Credential{Service: "agent-credential-broker", Account: "github-readonly"})
	if err != nil || value != "ghp_value" {
		t.Fatalf("got %q %v", value, err)
	}
	want := []string{"security", "find-generic-password", "-s", "agent-credential-broker", "-a", "github-readonly", "-w"}
	if !slices.Equal(got, want) {
		t.Fatalf("args %v", got)
	}
	failing := Keychain{Run: func(context.Context, string, ...string) ([]byte, error) { return nil, errors.New("exit 44") }}
	if _, err := failing.Fetch(context.Background(), policy.Credential{}); err == nil {
		t.Fatal("missing keychain item must fail")
	}
}
