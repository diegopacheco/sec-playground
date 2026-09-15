package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"agent-credential-broker-mcp/internal/audit"
	"agent-credential-broker-mcp/internal/broker"
	"agent-credential-broker-mcp/internal/httpapi"
	"agent-credential-broker-mcp/internal/policy"
)

const key = "sk_live_mcp_91b2"

type source struct{}

func (source) Fetch(context.Context, policy.Credential) (string, error) { return key, nil }

type session struct {
	t      *testing.T
	in     *io.PipeWriter
	out    *bufio.Scanner
	nextID int
}

func (s *session) call(method string, params any) map[string]any {
	s.t.Helper()
	s.nextID++
	raw, _ := json.Marshal(map[string]any{"jsonrpc": "2.0", "id": s.nextID, "method": method, "params": params})
	s.in.Write(append(raw, '\n'))
	if !s.out.Scan() {
		s.t.Fatal("no response")
	}
	var resp map[string]any
	json.Unmarshal(s.out.Bytes(), &resp)
	return resp
}

func (s *session) tool(name string, args map[string]any) (string, bool) {
	resp := s.call("tools/call", map[string]any{"name": name, "arguments": args})
	result := resp["result"].(map[string]any)
	text := result["content"].([]any)[0].(map[string]any)["text"].(string)
	return text, result["isError"].(bool)
}

func start(t *testing.T) (*session, *httpapi.Server, string) {
	up := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"seen": r.Header.Get("Authorization")})
	}))
	t.Cleanup(up.Close)
	p := &policy.Policy{Credentials: []policy.Credential{{
		Name: "vendor-api", Provider: "vault", Approval: "human",
		Inject: policy.Inject{Header: "Authorization", Template: "Bearer {{credential}}"},
		Allow:  []policy.Rule{{Scheme: "http", Host: strings.TrimPrefix(up.URL, "http://"), Methods: []string{"GET"}, PathPrefix: "/v1/"}},
	}}}
	log, _ := audit.Open("", time.Now)
	api := &httpapi.Server{Broker: broker.New(p, source{}, log, time.Now), AdminToken: "admin-test"}
	srv := httptest.NewServer(api.Handler())
	t.Cleanup(srv.Close)
	inR, inW := io.Pipe()
	outR, outW := io.Pipe()
	m := &Server{Backend: Backend{BaseURL: srv.URL}, Name: "t", Version: "1"}
	go m.Serve(context.Background(), inR, outW)
	t.Cleanup(func() { inW.Close() })
	return &session{t: t, in: inW, out: bufio.NewScanner(outR)}, api, srv.URL
}

func TestHandshakeAndToolListFollowMCP(t *testing.T) {
	s, _, _ := start(t)
	init := s.call("initialize", map[string]any{"protocolVersion": "2025-06-18", "capabilities": map[string]any{}, "clientInfo": map[string]string{"name": "t", "version": "1"}})
	if init["result"].(map[string]any)["protocolVersion"] != "2025-06-18" {
		t.Fatalf("bad initialize: %v", init)
	}
	s.in.Write([]byte(`{"jsonrpc":"2.0","method":"notifications/initialized"}` + "\n"))
	tools := s.call("tools/list", nil)["result"].(map[string]any)["tools"].([]any)
	if len(tools) != 4 {
		t.Fatalf("expected 4 tools, got %d", len(tools))
	}
	if s.call("nope", nil)["error"] == nil {
		t.Fatal("unknown method must return an error")
	}
}

func TestAgentFlowOverMCPNeverExposesTheCredential(t *testing.T) {
	s, api, base := start(t)
	upstream := api.Broker.Credentials()[0].Allow[0].Host
	target := "http://" + upstream + "/v1/whoami"
	text, isErr := s.tool("list_credentials", nil)
	if isErr || strings.Contains(text, key) {
		t.Fatalf("list_credentials leaked or failed: %s", text)
	}
	text, isErr = s.tool("request_access", map[string]any{"credential": "vendor-api", "method": "GET", "url": target, "reason": "who am i"})
	if isErr {
		t.Fatal(text)
	}
	var in broker.Intent
	json.Unmarshal([]byte(text), &in)
	text, _ = s.tool("check_access", map[string]any{"intent_id": in.ID})
	if strings.Contains(text, `"token":"ey`) || !strings.Contains(text, `"token":""`) {
		t.Fatalf("token must be empty before approval: %s", text)
	}
	req, _ := http.NewRequest("POST", base+"/api/intents/"+in.ID+"/approve", nil)
	if res, _ := http.DefaultClient.Do(req); res.StatusCode != http.StatusUnauthorized {
		t.Fatal("approval without admin token must be refused, the agent must not approve itself")
	}
	req.Header.Set("X-Admin-Token", "admin-test")
	if res, _ := http.DefaultClient.Do(req); res.StatusCode != http.StatusOK {
		t.Fatal("admin approval failed")
	}
	text, _ = s.tool("check_access", map[string]any{"intent_id": in.ID})
	var checked struct{ Token string }
	json.Unmarshal([]byte(text), &checked)
	text, isErr = s.tool("call_api", map[string]any{"token": checked.Token, "method": "GET", "url": target})
	if isErr || strings.Contains(text, key) || !strings.Contains(text, "REDACTED") {
		t.Fatalf("call leaked or failed: %s", text)
	}
	text, isErr = s.tool("call_api", map[string]any{"token": checked.Token, "method": "GET", "url": target})
	if !isErr || !strings.Contains(text, "already used") {
		t.Fatalf("replay must be an MCP tool error: %s", text)
	}
}

func TestCheckAccessRejectsPathInjection(t *testing.T) {
	s, _, _ := start(t)
	text, isErr := s.tool("check_access", map[string]any{"intent_id": "../../api/audit"})
	if !isErr {
		t.Fatalf("path injection must be rejected: %s", text)
	}
}
