package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"mcp-output-taint-tracker/internal/session"
)

type hookInput struct {
	SessionID    string         `json:"session_id"`
	Event        string         `json:"hook_event_name"`
	Prompt       string         `json:"prompt"`
	ToolName     string         `json:"tool_name"`
	ToolInput    map[string]any `json:"tool_input"`
	ToolResponse any            `json:"tool_response"`
}

var client = &http.Client{Timeout: 5 * time.Second}

func main() {
	var in hookInput
	if err := json.NewDecoder(os.Stdin).Decode(&in); err != nil {
		deny("taint hook could not read hook input: " + err.Error())
	}
	base := os.Getenv("TAINT_URL")
	if base == "" {
		base = "http://127.0.0.1:8751"
	}
	if in.SessionID == "" {
		in.SessionID = "default"
	}
	path := base + "/api/sessions/" + in.SessionID
	switch in.Event {
	case "UserPromptSubmit":
		send(path+"/prompt", map[string]any{"text": in.Prompt})
	case "PostToolUse":
		send(path+"/result", map[string]any{"tool": in.ToolName, "args": in.ToolInput, "output": stringify(in.ToolResponse)})
	case "PreToolUse":
		ev := send(path+"/call", map[string]any{"tool": in.ToolName, "args": in.ToolInput})
		if ev.Decision != nil && !ev.Decision.Allow {
			deny(explain(ev.Decision))
		}
	}
}

func send(url string, body any) session.Event {
	b, _ := json.Marshal(body)
	resp, err := client.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {
		deny("taint tracker unreachable, failing closed: " + err.Error())
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		deny(fmt.Sprintf("taint tracker answered %d: %s", resp.StatusCode, raw))
	}
	var ev session.Event
	if err := json.Unmarshal(raw, &ev); err != nil {
		deny("taint tracker sent invalid json: " + err.Error())
	}
	return ev
}

func explain(d *session.Decision) string {
	var b bytes.Buffer
	b.WriteString("blocked by taint tracker:")
	for _, f := range d.Flows {
		fmt.Fprintf(&b, "\n- %s data from %s (event %d) reached argument %q via %s match %q", f.Label, f.SourceTool, f.SourceID, f.Arg, f.Encoding, f.Evidence)
	}
	return b.String()
}

func stringify(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	b, _ := json.Marshal(v)
	return string(b)
}

func deny(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(2)
}
