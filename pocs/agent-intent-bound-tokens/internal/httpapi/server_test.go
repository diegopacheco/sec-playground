package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"agent-intent-bound-tokens/internal/token"
	"agent-intent-bound-tokens/internal/verification"
)

func TestRootListsEndpoints(t *testing.T) {
	tokens, err := token.New("test-intent-signing-key-32-bytes-minimum")
	if err != nil {
		t.Fatal(err)
	}
	server := New("intent.local", tokens, verification.New(tokens))
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	server.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}
	var body struct {
		Endpoints []endpoint `json:"endpoints"`
	}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if len(body.Endpoints) != 5 {
		t.Fatalf("expected 5 endpoints, got %d", len(body.Endpoints))
	}
	if body.Endpoints[2].Method != http.MethodPost || body.Endpoints[2].Path != "/api/tokens" {
		t.Fatalf("unexpected endpoint: %+v", body.Endpoints[2])
	}
}

func TestUnknownPathReturnsNotFound(t *testing.T) {
	tokens, err := token.New("test-intent-signing-key-32-bytes-minimum")
	if err != nil {
		t.Fatal(err)
	}
	server := New("intent.local", tokens, verification.New(tokens))
	request := httptest.NewRequest(http.MethodGet, "/unknown", nil)
	response := httptest.NewRecorder()

	server.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, response.Code)
	}
}
