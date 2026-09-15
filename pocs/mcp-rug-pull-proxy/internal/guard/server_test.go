package guard

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"mcp-rug-pull-proxy/internal/fakeserver"
	"mcp-rug-pull-proxy/internal/mcp"
	"mcp-rug-pull-proxy/internal/store"
)

var tokens = Tokens{Admin: "admin-token-0123456789abcdef0123456789", Proxy: "proxy-token-0123456789abcdef0123456789"}

func newServer(t *testing.T) http.Handler {
	t.Helper()
	st, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return New(st, tokens).Handler()
}

func call(h http.Handler, method, path, header, token string, body any) *httptest.ResponseRecorder {
	raw, _ := json.Marshal(body)
	req := httptest.NewRequest(method, "http://localhost:8741"+path, bytes.NewReader(raw))
	if header != "" {
		req.Header.Set(header, token)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func checkBody(tools ...mcp.Tool) map[string]any {
	return map[string]any{"server": "weather", "tools": tools}
}

func TestCheckRequiresTheProxyTokenSoAnAgentCannotPreRegisterTools(t *testing.T) {
	h := newServer(t)
	if rec := call(h, "POST", "/api/check", "", "", checkBody(fakeserver.WeatherTool())); rec.Code != http.StatusUnauthorized {
		t.Fatalf("check without token returned %d", rec.Code)
	}
	if rec := call(h, "POST", "/api/check", "X-Proxy-Token", tokens.Admin, checkBody(fakeserver.WeatherTool())); rec.Code != http.StatusUnauthorized {
		t.Fatalf("check with the admin token returned %d", rec.Code)
	}
}

func TestApproveRequiresAdminTokenAndTheExactHash(t *testing.T) {
	h := newServer(t)
	rec := call(h, "POST", "/api/check", "X-Proxy-Token", tokens.Proxy, checkBody(fakeserver.WeatherTool()))
	var out struct {
		Decisions []store.Decision `json:"decisions"`
	}
	_ = json.Unmarshal(rec.Body.Bytes(), &out)
	hash := out.Decisions[0].Hash
	approve := map[string]string{"server": "weather", "tool": "get_weather", "hash": hash}
	if rec := call(h, "POST", "/api/tools/approve", "X-Proxy-Token", tokens.Proxy, approve); rec.Code != http.StatusUnauthorized {
		t.Fatalf("proxy token approved a tool: %d", rec.Code)
	}
	wrong := map[string]string{"server": "weather", "tool": "get_weather", "hash": "sha256:deadbeef"}
	if rec := call(h, "POST", "/api/tools/approve", "X-Admin-Token", tokens.Admin, wrong); rec.Code != http.StatusConflict {
		t.Fatalf("wrong hash approve returned %d", rec.Code)
	}
	if rec := call(h, "POST", "/api/tools/approve", "X-Admin-Token", tokens.Admin, approve); rec.Code != http.StatusOK {
		t.Fatalf("exact hash approve returned %d %s", rec.Code, rec.Body.String())
	}
}

func TestNonLocalHostHeaderIsRefusedToBlockDnsRebinding(t *testing.T) {
	h := newServer(t)
	req := httptest.NewRequest("GET", "http://attacker.test/api/tools", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("foreign host header returned %d", rec.Code)
	}
}

func TestScanEndpointReportsPoisonedTool(t *testing.T) {
	h := newServer(t)
	rec := call(h, "POST", "/api/scan", "", "", fakeserver.PoisonedAddTool())
	var report struct {
		Severity string `json:"severity"`
	}
	_ = json.Unmarshal(rec.Body.Bytes(), &report)
	if rec.Code != http.StatusOK || report.Severity != "high" {
		t.Fatalf("scan returned %d severity %s", rec.Code, report.Severity)
	}
}
