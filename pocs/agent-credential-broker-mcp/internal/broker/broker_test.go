package broker

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"agent-credential-broker-mcp/internal/audit"
	"agent-credential-broker-mcp/internal/intent"
	"agent-credential-broker-mcp/internal/policy"
)

const key = "sk_live_7f3a9c1e5b"

type fixedSource struct {
	value string
	err   error
	hits  int
}

func (f *fixedSource) Fetch(context.Context, policy.Credential) (string, error) {
	f.hits++
	return f.value, f.err
}

type env struct {
	b         *Broker
	src       *fixedSource
	upstream  *httptest.Server
	collector *atomic.Int32
}

func setup(t *testing.T, approval string) *env {
	t.Helper()
	collected := &atomic.Int32{}
	collector := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { collected.Add(1) }))
	t.Cleanup(collector.Close)
	up := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer "+key {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		switch r.URL.Path {
		case "/v1/echo":
			json.NewEncoder(w).Encode(map[string]string{"auth": r.Header.Get("Authorization")})
		case "/v1/redirect":
			http.Redirect(w, r, collector.URL+"/grab", http.StatusFound)
		default:
			w.Write([]byte(`{"ok":true}`))
		}
	}))
	t.Cleanup(up.Close)
	host := strings.TrimPrefix(up.URL, "http://")
	p := &policy.Policy{Credentials: []policy.Credential{{
		Name: "vendor-api", Provider: "vault", Approval: approval,
		Inject: policy.Inject{Header: "Authorization", Template: "Bearer {{credential}}"},
		Allow:  []policy.Rule{{Scheme: "http", Host: host, Methods: []string{"GET", "POST"}, PathPrefix: "/v1/"}},
	}}}
	log, _ := audit.Open("", time.Now)
	src := &fixedSource{value: key}
	return &env{b: New(p, src, log, time.Now), src: src, upstream: up, collector: collected}
}

func (e *env) approved(t *testing.T, method, path, body string) (Intent, string) {
	t.Helper()
	in, err := e.b.RequestAccess(AccessRequest{Credential: "vendor-api", Method: method, URL: e.upstream.URL + path, Body: body, Reason: "test"})
	if err != nil {
		t.Fatal(err)
	}
	if in.Status == StatusPending {
		if _, err := e.b.Approve(in.ID); err != nil {
			t.Fatal(err)
		}
	}
	_, token, _ := e.b.Check(in.ID)
	if token == "" {
		t.Fatal("approved intent must expose a token")
	}
	return in, token
}

func TestPendingIntentHasNoTokenSoTheAgentCannotActBeforeAHumanDecides(t *testing.T) {
	e := setup(t, "human")
	in, err := e.b.RequestAccess(AccessRequest{Credential: "vendor-api", Method: "GET", URL: e.upstream.URL + "/v1/profile"})
	if err != nil {
		t.Fatal(err)
	}
	if _, token, _ := e.b.Check(in.ID); token != "" || in.Status != StatusPending {
		t.Fatalf("pending intent leaked a token: status=%s", in.Status)
	}
	if _, err := e.b.Deny(in.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := e.b.Approve(in.ID); !errors.Is(err, ErrNotPending) {
		t.Fatalf("denied intent must not be approvable later, got %v", err)
	}
}

func TestAutoApprovalIssuesTokenImmediately(t *testing.T) {
	e := setup(t, "auto")
	in, err := e.b.RequestAccess(AccessRequest{Credential: "vendor-api", Method: "GET", URL: e.upstream.URL + "/v1/profile"})
	if err != nil || in.Status != StatusApproved {
		t.Fatalf("auto policy must approve, got %s %v", in.Status, err)
	}
}

func TestCredentialNeverAppearsInAnythingTheAgentReceivesEvenWhenTheVendorEchoesIt(t *testing.T) {
	e := setup(t, "human")
	_, token := e.approved(t, "GET", "/v1/echo", "")
	res, err := e.b.Call(context.Background(), CallRequest{Token: token, Method: "GET", URL: e.upstream.URL + "/v1/echo"})
	if err != nil {
		t.Fatal(err)
	}
	raw, _ := json.Marshal(res)
	if strings.Contains(string(raw), key) || res.Redactions == 0 {
		t.Fatalf("credential leaked to agent: %s", raw)
	}
	for _, entry := range e.b.Audit().Entries() {
		line, _ := json.Marshal(entry)
		if strings.Contains(string(line), key) {
			t.Fatalf("credential leaked to audit: %s", line)
		}
	}
}

func TestRedirectIsNotFollowedSoUpstreamCannotBounceTheCallToAnotherHost(t *testing.T) {
	e := setup(t, "human")
	_, token := e.approved(t, "GET", "/v1/redirect", "")
	res, err := e.b.Call(context.Background(), CallRequest{Token: token, Method: "GET", URL: e.upstream.URL + "/v1/redirect"})
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != http.StatusFound || e.collector.Load() != 0 {
		t.Fatalf("redirect followed: status=%d collector hits=%d", res.Status, e.collector.Load())
	}
}

func TestHostOutsideAllowlistIsDeniedBeforeAnyTokenOrCredentialFetch(t *testing.T) {
	e := setup(t, "auto")
	collector := strings.Replace(e.upstream.URL, "127.0.0.1", "localhost", 1)
	_, err := e.b.RequestAccess(AccessRequest{Credential: "vendor-api", Method: "POST", URL: collector + "/v1/upload"})
	if err == nil {
		t.Fatal("different host must be denied")
	}
	if e.src.hits != 0 {
		t.Fatal("credential must not be fetched for a denied request")
	}
	last := e.b.Audit().Entries()
	if last[len(last)-1].Outcome != "denied" {
		t.Fatal("denial must be audited")
	}
}

func TestReplayAndTamperAreBlockedAndMarkTheIntentUsed(t *testing.T) {
	e := setup(t, "human")
	in, token := e.approved(t, "GET", "/v1/profile", "")
	if _, err := e.b.Call(context.Background(), CallRequest{Token: token, Method: "POST", URL: e.upstream.URL + "/v1/profile", Body: "{}"}); !errors.Is(err, intent.ErrMismatch) {
		t.Fatalf("tampered call must mismatch, got %v", err)
	}
	if _, err := e.b.Call(context.Background(), CallRequest{Token: token, Method: "GET", URL: e.upstream.URL + "/v1/profile"}); !errors.Is(err, intent.ErrReplay) {
		t.Fatalf("token must be burned after tamper, got %v", err)
	}
	if got, _, _ := e.b.Check(in.ID); got.Status != StatusUsed {
		t.Fatalf("intent must be used, got %s", got.Status)
	}
	if e.src.hits != 0 {
		t.Fatal("credential must not be fetched for rejected calls")
	}
}

func TestPolicyIsRecheckedAtCallTime(t *testing.T) {
	e := setup(t, "human")
	_, token := e.approved(t, "GET", "/v1/profile", "")
	e.b.policy.Credentials[0].Allow[0].Methods = []string{"POST"}
	if _, err := e.b.Call(context.Background(), CallRequest{Token: token, Method: "GET", URL: e.upstream.URL + "/v1/profile"}); err == nil {
		t.Fatal("revoked method must be blocked even with a valid token")
	}
}

func TestProviderFailureStopsTheCall(t *testing.T) {
	e := setup(t, "human")
	e.src.err = errors.New("vault sealed")
	_, token := e.approved(t, "GET", "/v1/profile", "")
	if _, err := e.b.Call(context.Background(), CallRequest{Token: token, Method: "GET", URL: e.upstream.URL + "/v1/profile"}); err == nil {
		t.Fatal("call must fail when the credential cannot be fetched")
	}
}

func TestPendingIntentExpires(t *testing.T) {
	e := setup(t, "human")
	now := time.Now()
	e.b.now = func() time.Time { return now }
	in, _ := e.b.RequestAccess(AccessRequest{Credential: "vendor-api", Method: "GET", URL: e.upstream.URL + "/v1/profile"})
	now = now.Add(e.b.PendingTTL + time.Second)
	if _, err := e.b.Approve(in.ID); !errors.Is(err, ErrNotPending) {
		t.Fatalf("expired intent must not be approvable, got %v", err)
	}
}

func TestUnknownCredentialIsDenied(t *testing.T) {
	e := setup(t, "auto")
	u, _ := url.JoinPath(e.upstream.URL, "/v1/profile")
	if _, err := e.b.RequestAccess(AccessRequest{Credential: "prod-db", Method: "GET", URL: u}); !errors.Is(err, ErrUnknownCred) {
		t.Fatalf("got %v", err)
	}
}
