package lab

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"agent-credential-broker-mcp/internal/audit"
	"agent-credential-broker-mcp/internal/broker"
	"agent-credential-broker-mcp/internal/policy"
	"agent-credential-broker-mcp/internal/vendor"
)

const key = "sk_live_lab_44d1"

type source struct{}

func (source) Fetch(context.Context, policy.Credential) (string, error) { return key, nil }

func outcomes(steps []Step) string {
	var parts []string
	for _, s := range steps {
		parts = append(parts, s.Actor+":"+s.Outcome)
	}
	return strings.Join(parts, " ")
}

func TestEveryScenarioEndsWithTheExpectedControlFiring(t *testing.T) {
	up := httptest.NewServer(vendor.Handler(key))
	defer up.Close()
	p := &policy.Policy{Credentials: []policy.Credential{{
		Name: credential, Provider: "vault", Approval: "human",
		Inject: policy.Inject{Header: "Authorization", Template: "Bearer {{credential}}"},
		Allow:  []policy.Rule{{Scheme: "http", Host: strings.TrimPrefix(up.URL, "http://"), Methods: []string{"GET", "POST"}, PathPrefix: "/v1/"}},
	}}}
	want := map[string]string{
		"approved-call":    "agent:ok human:ok broker:ok",
		"host-not-allowed": "agent:blocked",
		"path-traversal":   "agent:blocked",
		"token-replay":     "agent:ok human:ok broker:ok broker:blocked",
		"tampered-request": "agent:ok human:ok broker:blocked broker:blocked",
		"echo-leak":        "agent:ok human:ok broker:ok",
		"redirect-exfil":   "agent:ok human:ok broker:ok",
	}
	if len(want) != len(Scenarios) {
		t.Fatalf("every scenario needs an expectation")
	}
	for _, sc := range Scenarios {
		log, _ := audit.Open("", time.Now)
		b := broker.New(p, source{}, log, time.Now)
		steps, err := Run(context.Background(), b, up.URL, sc.ID)
		if err != nil {
			t.Fatal(err)
		}
		if got := outcomes(steps); got != want[sc.ID] {
			t.Errorf("%s: got %q want %q", sc.ID, got, want[sc.ID])
		}
		raw, _ := json.Marshal(steps)
		if strings.Contains(string(raw), key) {
			t.Errorf("%s leaked the credential into the transcript", sc.ID)
		}
		if sc.ID == "echo-leak" && !strings.Contains(string(raw), "REDACTED") {
			t.Errorf("echo-leak must show redaction")
		}
		if sc.ID == "redirect-exfil" && !strings.Contains(string(raw), "status 302") {
			t.Errorf("redirect-exfil must return the 302 without following it")
		}
	}
}
