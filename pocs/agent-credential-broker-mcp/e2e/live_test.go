//go:build e2e

package e2e

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"
)

type client struct {
	t   *testing.T
	in  io.WriteCloser
	out *bufio.Scanner
	id  int
}

func (c *client) tool(name string, args map[string]any) (string, bool) {
	c.t.Helper()
	c.id++
	raw, _ := json.Marshal(map[string]any{"jsonrpc": "2.0", "id": c.id, "method": "tools/call", "params": map[string]any{"name": name, "arguments": args}})
	fmt.Fprintf(c.in, "%s\n", raw)
	if !c.out.Scan() {
		c.t.Fatal("mcp bridge closed")
	}
	var resp struct {
		Result struct {
			Content []struct{ Text string }
			IsError bool
		}
	}
	json.Unmarshal(c.out.Bytes(), &resp)
	return resp.Result.Content[0].Text, resp.Result.IsError
}

func TestLiveStackThroughTheRealMCPBinary(t *testing.T) {
	brokerURL := os.Getenv("BROKER_URL")
	vendorURL := os.Getenv("VENDOR_URL")
	admin, err := os.ReadFile(os.Getenv("ADMIN_TOKEN_FILE"))
	if err != nil || brokerURL == "" || vendorURL == "" {
		t.Fatalf("BROKER_URL, VENDOR_URL and ADMIN_TOKEN_FILE are required: %v", err)
	}
	key, err := os.ReadFile(os.Getenv("VENDOR_KEY_FILE"))
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(os.Getenv("MCP_BIN"))
	cmd.Env = append(os.Environ(), "BROKER_URL="+brokerURL)
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	defer cmd.Process.Kill()
	c := &client{t: t, in: stdin, out: bufio.NewScanner(stdout)}

	text, isErr := c.tool("request_access", map[string]any{"credential": "vendor-api", "method": "GET", "url": vendorURL + "/v1/debug/echo-auth", "reason": "live e2e"})
	if isErr {
		t.Fatal(text)
	}
	var in struct{ ID string }
	json.Unmarshal([]byte(text), &in)

	req, _ := http.NewRequest("POST", brokerURL+"/api/intents/"+in.ID+"/approve", nil)
	req.Header.Set("X-Admin-Token", strings.TrimSpace(string(admin)))
	res, err := http.DefaultClient.Do(req)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("approve failed: %v", err)
	}

	text, _ = c.tool("check_access", map[string]any{"intent_id": in.ID})
	var checked struct{ Token string }
	json.Unmarshal([]byte(text), &checked)

	text, isErr = c.tool("call_api", map[string]any{"token": checked.Token, "method": "GET", "url": vendorURL + "/v1/debug/echo-auth"})
	if isErr {
		t.Fatalf("call failed: %s", text)
	}
	if strings.Contains(text, strings.TrimSpace(string(key))) || !strings.Contains(text, "REDACTED:vendor-api") {
		t.Fatalf("vault credential leaked through the live stack: %s", text)
	}

	text, isErr = c.tool("request_access", map[string]any{"credential": "vendor-api", "method": "POST", "url": "https://collector.invalid/upload", "reason": "exfil"})
	if !isErr {
		t.Fatalf("exfil host must be denied: %s", text)
	}
	t.Logf("live stack ok: vault credential injected, echoed by vendor, redacted before reaching the agent")
}
