package mcp

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type rpcRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type rpcResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Result  any             `json:"result,omitempty"`
	Error   *rpcError       `json:"error,omitempty"`
}

type Tool struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	InputSchema map[string]any `json:"inputSchema"`
}

func object(required []string, props map[string]any) map[string]any {
	return map[string]any{"type": "object", "properties": props, "required": required, "additionalProperties": false}
}

func str(desc string) map[string]any {
	return map[string]any{"type": "string", "description": desc}
}

var Tools = []Tool{
	{"list_credentials", "List credential names the broker can use on your behalf, with allowed hosts, methods and approval mode. Never returns credential values.", object([]string{}, map[string]any{})},
	{"request_access", "Declare one exact HTTP request you want to make with a credential. Returns an intent id. Human approval may be required before a token is issued.", object([]string{"credential", "method", "url", "reason"}, map[string]any{
		"credential": str("credential name from list_credentials"),
		"method":     str("HTTP method"),
		"url":        str("absolute URL of the request"),
		"body":       str("exact request body, empty for none"),
		"reason":     str("why this request is needed, shown to the approver"),
	})},
	{"check_access", "Check an intent. When approved, returns a single use token bound to that exact request.", object([]string{"intent_id"}, map[string]any{
		"intent_id": str("intent id from request_access"),
	})},
	{"call_api", "Execute the approved request. The broker injects the credential, blocks redirects and redacts the credential from the response. The token is burned on use.", object([]string{"token", "method", "url"}, map[string]any{
		"token":  str("token from check_access"),
		"method": str("same method that was approved"),
		"url":    str("same URL that was approved"),
		"body":   str("same body that was approved"),
	})},
}

type Backend struct {
	BaseURL string
	Client  *http.Client
}

func (b Backend) do(ctx context.Context, method, path string, body any) (json.RawMessage, bool, error) {
	var reader io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			return nil, false, err
		}
		reader = bytes.NewReader(raw)
	}
	req, err := http.NewRequestWithContext(ctx, method, strings.TrimRight(b.BaseURL, "/")+path, reader)
	if err != nil {
		return nil, false, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := b.Client
	if client == nil {
		client = &http.Client{Timeout: 20 * time.Second}
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, false, fmt.Errorf("broker unreachable: %w", err)
	}
	defer res.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(res.Body, 2<<20))
	return raw, res.StatusCode < 300, err
}

type Server struct {
	Backend Backend
	Name    string
	Version string
}

func (s *Server) Serve(ctx context.Context, in io.Reader, out io.Writer) error {
	sc := bufio.NewScanner(in)
	sc.Buffer(make([]byte, 4<<20), 4<<20)
	enc := json.NewEncoder(out)
	for sc.Scan() {
		line := bytes.TrimSpace(sc.Bytes())
		if len(line) == 0 {
			continue
		}
		var req rpcRequest
		if err := json.Unmarshal(line, &req); err != nil {
			enc.Encode(rpcResponse{JSONRPC: "2.0", ID: json.RawMessage("null"), Error: &rpcError{-32700, "parse error"}})
			continue
		}
		if len(req.ID) == 0 {
			continue
		}
		result, rerr := s.handle(ctx, req)
		resp := rpcResponse{JSONRPC: "2.0", ID: req.ID, Result: result, Error: rerr}
		if err := enc.Encode(resp); err != nil {
			return err
		}
	}
	return sc.Err()
}

func (s *Server) handle(ctx context.Context, req rpcRequest) (any, *rpcError) {
	switch req.Method {
	case "initialize":
		var p struct {
			ProtocolVersion string `json:"protocolVersion"`
		}
		json.Unmarshal(req.Params, &p)
		if p.ProtocolVersion == "" {
			p.ProtocolVersion = "2025-06-18"
		}
		return map[string]any{
			"protocolVersion": p.ProtocolVersion,
			"capabilities":    map[string]any{"tools": map[string]any{}},
			"serverInfo":      map[string]string{"name": s.Name, "version": s.Version},
		}, nil
	case "ping":
		return map[string]any{}, nil
	case "tools/list":
		return map[string]any{"tools": Tools}, nil
	case "tools/call":
		var p struct {
			Name      string          `json:"name"`
			Arguments json.RawMessage `json:"arguments"`
		}
		if err := json.Unmarshal(req.Params, &p); err != nil {
			return nil, &rpcError{-32602, "invalid params"}
		}
		return s.callTool(ctx, p.Name, p.Arguments)
	}
	return nil, &rpcError{-32601, "method not found: " + req.Method}
}

func toolResult(raw json.RawMessage, ok bool, err error) (any, *rpcError) {
	text := string(raw)
	if err != nil {
		text, ok = err.Error(), false
	}
	return map[string]any{
		"content": []map[string]string{{"type": "text", "text": text}},
		"isError": !ok,
	}, nil
}

func (s *Server) callTool(ctx context.Context, name string, args json.RawMessage) (any, *rpcError) {
	var fields map[string]any
	if len(args) > 0 {
		if err := json.Unmarshal(args, &fields); err != nil {
			return nil, &rpcError{-32602, "arguments must be an object"}
		}
	}
	switch name {
	case "list_credentials":
		return toolResult(s.Backend.do(ctx, http.MethodGet, "/agent/credentials", nil))
	case "request_access":
		return toolResult(s.Backend.do(ctx, http.MethodPost, "/agent/intents", fields))
	case "check_access":
		id, _ := fields["intent_id"].(string)
		if id == "" || strings.ContainsAny(id, "/?#") {
			return toolResult(nil, false, fmt.Errorf("intent_id is required"))
		}
		return toolResult(s.Backend.do(ctx, http.MethodGet, "/agent/intents/"+id, nil))
	case "call_api":
		return toolResult(s.Backend.do(ctx, http.MethodPost, "/agent/call", fields))
	}
	return nil, &rpcError{-32602, "unknown tool: " + name}
}
