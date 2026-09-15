package proxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"mcp-rug-pull-proxy/internal/mcp"
	"mcp-rug-pull-proxy/internal/store"
)

type Guard interface {
	Check(server string, tools []mcp.Tool) ([]store.Decision, error)
	Denied(server, tool, reason string)
}

type HTTPGuard struct {
	URL    string
	Token  string
	Client *http.Client
}

func NewHTTPGuard(url, token string) *HTTPGuard {
	return &HTTPGuard{URL: url, Token: token, Client: &http.Client{Timeout: 3 * time.Second}}
}

func (g *HTTPGuard) Check(server string, tools []mcp.Tool) ([]store.Decision, error) {
	var out struct {
		Decisions []store.Decision `json:"decisions"`
	}
	if err := g.post("/api/check", map[string]any{"server": server, "tools": tools}, &out); err != nil {
		return nil, err
	}
	return out.Decisions, nil
}

func (g *HTTPGuard) Denied(server, tool, reason string) {
	_ = g.post("/api/denied", map[string]string{"server": server, "tool": tool, "reason": reason}, nil)
}

func (g *HTTPGuard) post(path string, body, out any) error {
	raw, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, g.URL+path, bytes.NewReader(raw))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Proxy-Token", g.Token)
	resp, err := g.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("guard answered %s", resp.Status)
	}
	if out == nil {
		return nil
	}
	return json.NewDecoder(resp.Body).Decode(out)
}
