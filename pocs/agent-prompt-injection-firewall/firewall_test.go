package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func testFirewall(t *testing.T) *Firewall {
	t.Helper()
	return NewFirewall(FirewallConfig{
		ReadRoot:      "fixtures/data",
		AllowedHosts:  []string{"127.0.0.1:8091"},
		Canary:        "FIREWALL_CANARY_TEST",
		ApprovalToken: "approved",
		LocalURL:      "http://127.0.0.1:8091",
		AuditPath:     filepath.Join(t.TempDir(), "audit.jsonl"),
	})
}

func TestReadCapabilityBoundary(t *testing.T) {
	firewall := testFirewall(t)
	decision := firewall.authorize(Candidate{Tool: "read", Target: "../../../../etc/passwd"}, "")
	if decision.Allowed || decision.Reason != "read target outside capability root" {
		t.Fatalf("unexpected decision: %+v", decision)
	}
}

func TestHTTPDestinationAllowlist(t *testing.T) {
	firewall := testFirewall(t)
	decision := firewall.authorize(Candidate{Tool: "http", Method: "GET", Target: "http://attacker.invalid/collect"}, "")
	if decision.Allowed || decision.Reason != "HTTP destination not allowed" {
		t.Fatalf("unexpected decision: %+v", decision)
	}
}

func TestSensitiveCanaryBlockedAndRedacted(t *testing.T) {
	firewall := testFirewall(t)
	decision := firewall.Execute("poisoned.txt", Candidate{Tool: "http", Method: "GET", Target: "${LOCAL_URL}/receiver?secret=${CANARY}"}, "")
	if decision.Allowed || decision.Reason != "sensitive canary detected" {
		t.Fatalf("unexpected decision: %+v", decision)
	}
	if decision.Target == "http://127.0.0.1:8091/receiver?secret=FIREWALL_CANARY_TEST" {
		t.Fatal("canary was not redacted")
	}
}

func TestHTTPWriteRequiresApproval(t *testing.T) {
	firewall := testFirewall(t)
	candidate := Candidate{Tool: "http", Method: "POST", Target: "http://127.0.0.1:8091/receiver"}
	denied := firewall.authorize(candidate, "")
	if denied.Allowed || denied.Reason != "explicit approval required" {
		t.Fatalf("unexpected denial: %+v", denied)
	}
	allowed := firewall.authorize(candidate, "approved")
	if !allowed.Allowed || allowed.Reason != "approved HTTP write allowed" {
		t.Fatalf("unexpected approval: %+v", allowed)
	}
}

func TestIndexListsEndpoints(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()
	newMux(&Server{}).ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", response.Code)
	}
	var body map[string][]string
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	want := []string{"GET /", "GET /health", "POST /run", "GET /audit", "GET /receiver", "POST /receiver"}
	if len(body["endpoints"]) != len(want) {
		t.Fatalf("unexpected endpoints: %v", body["endpoints"])
	}
	for i := range want {
		if body["endpoints"][i] != want[i] {
			t.Fatalf("unexpected endpoints: %v", body["endpoints"])
		}
	}
}

func TestUnknownEndpointReturnsNotFound(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/missing", nil)
	response := httptest.NewRecorder()
	newMux(&Server{}).ServeHTTP(response, request)
	if response.Code != http.StatusNotFound {
		t.Fatalf("unexpected status: %d", response.Code)
	}
}
