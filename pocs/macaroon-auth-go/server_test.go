package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCapabilityAPIFlow(t *testing.T) {
	service := newServer([]byte("root-key"))
	handler := service.routes()
	minted := requestToken(t, handler, "/api/macaroon", map[string]any{
		"resource": "/vault/*", "operation": "*", "location": "*", "expires_in_seconds": 300,
	})
	attenuated := requestToken(t, handler, "/api/attenuate", map[string]any{
		"token": minted.Token, "caveat": map[string]string{"type": "operation", "value": "read"},
	})
	allowed := requestVerification(t, handler, attenuated.Token, "/vault/report", "read", "us-west")
	denied := requestVerification(t, handler, attenuated.Token, "/vault/report", "delete", "us-west")
	if !allowed.Allowed || denied.Allowed {
		t.Fatal("API did not enforce the capability")
	}
}

func TestResourceRestrictionRejectsOtherResource(t *testing.T) {
	service := newServer([]byte("root-key"))
	handler := service.routes()
	minted := requestToken(t, handler, "/api/macaroon", map[string]any{
		"resource": "/records/*", "operation": "*", "location": "*", "expires_in_seconds": 300,
	})
	attenuated := requestToken(t, handler, "/api/attenuate", map[string]any{
		"token": minted.Token, "caveat": map[string]string{"type": "resource", "value": "/records/payroll"},
	})
	payroll := requestVerification(t, handler, attenuated.Token, "/records/payroll", "read", "us-west")
	budget := requestVerification(t, handler, attenuated.Token, "/records/budget", "read", "us-west")
	if !payroll.Allowed || budget.Allowed {
		t.Fatal("resource restriction was not enforced")
	}
}

func requestToken(t *testing.T, handler http.Handler, path string, body any) tokenResponse {
	t.Helper()
	data, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(data))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code < 200 || response.Code >= 300 {
		t.Fatalf("request failed with %d: %s", response.Code, response.Body.String())
	}
	var token tokenResponse
	if err := json.NewDecoder(response.Body).Decode(&token); err != nil {
		t.Fatal(err)
	}
	return token
}

func requestVerification(t *testing.T, handler http.Handler, token, resource, operation, location string) Verification {
	t.Helper()
	data, err := json.Marshal(map[string]string{
		"token": token, "resource": resource, "operation": operation, "location": location,
	})
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodPost, "/api/verify", bytes.NewReader(data))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	var result Verification
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	return result
}
