package sim

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"mcp-rug-pull-proxy/internal/store"
)

type Admin struct {
	URL    string
	Token  string
	Client *http.Client
}

func NewAdmin(url, token string) *Admin {
	return &Admin{URL: url, Token: token, Client: &http.Client{Timeout: 5 * time.Second}}
}

func (a *Admin) Pins(server string) (map[string]store.Pin, error) {
	var out struct {
		Pins []store.Pin `json:"pins"`
	}
	if err := a.do(http.MethodGet, "/api/tools", nil, &out); err != nil {
		return nil, err
	}
	pins := map[string]store.Pin{}
	for _, p := range out.Pins {
		if p.Server == server {
			pins[p.Tool] = p
		}
	}
	return pins, nil
}

func (a *Admin) Alerts(server string) ([]store.Alert, error) {
	var out struct {
		Alerts []store.Alert `json:"alerts"`
	}
	if err := a.do(http.MethodGet, "/api/alerts", nil, &out); err != nil {
		return nil, err
	}
	var alerts []store.Alert
	for _, al := range out.Alerts {
		if al.Server == server {
			alerts = append(alerts, al)
		}
	}
	return alerts, nil
}

func (a *Admin) Approve(server, tool, hash string) error {
	return a.do(http.MethodPost, "/api/tools/approve", map[string]string{"server": server, "tool": tool, "hash": hash}, nil)
}

func (a *Admin) Forget(server string) error {
	return a.do(http.MethodPost, "/api/servers/forget", map[string]string{"server": server}, nil)
}

func (a *Admin) do(method, path string, body, out any) error {
	var reader io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reader = bytes.NewReader(raw)
	}
	req, err := http.NewRequest(method, a.URL+path, reader)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Admin-Token", a.Token)
	resp, err := a.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		msg, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("%s %s: %s %s", method, path, resp.Status, bytes.TrimSpace(msg))
	}
	if out == nil {
		return nil
	}
	return json.NewDecoder(resp.Body).Decode(out)
}
