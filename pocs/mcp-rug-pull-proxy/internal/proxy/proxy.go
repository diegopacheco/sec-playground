package proxy

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
	"sync"

	"mcp-rug-pull-proxy/internal/mcp"
	"mcp-rug-pull-proxy/internal/store"
)

const CodeToolBlocked = -32001

type Proxy struct {
	server     string
	guard      Guard
	logw       io.Writer
	toClient   *mcp.Writer
	toUpstream *mcp.Writer
	mu         sync.Mutex
	lists      map[string]bool
	internal   map[string]bool
	allowed    map[string]bool
	seq        int
}

func New(server string, guard Guard, logw io.Writer) *Proxy {
	return &Proxy{
		server:   server,
		guard:    guard,
		logw:     logw,
		lists:    map[string]bool{},
		internal: map[string]bool{},
		allowed:  map[string]bool{},
	}
}

func (p *Proxy) Run(clientIn io.Reader, clientOut io.Writer, upstreamIn io.Writer, upstreamOut io.Reader) error {
	p.toClient = mcp.NewWriter(clientOut)
	p.toUpstream = mcp.NewWriter(upstreamIn)
	errs := make(chan error, 2)
	go func() { errs <- p.pumpClient(clientIn) }()
	go func() { errs <- p.pumpUpstream(upstreamOut) }()
	err := <-errs
	if errors.Is(err, io.EOF) {
		return nil
	}
	return err
}

func (p *Proxy) pumpClient(in io.Reader) error {
	reader := mcp.NewReader(in)
	for {
		line, err := reader.ReadLine()
		if err != nil {
			return err
		}
		var msg mcp.Message
		if json.Unmarshal(line, &msg) == nil && msg.IsRequest() {
			switch msg.Method {
			case mcp.MethodToolsList:
				p.mu.Lock()
				p.lists[string(msg.ID)] = !hasCursor(msg.Params)
				p.mu.Unlock()
			case mcp.MethodToolsCall:
				if blocked, reason := p.blockCall(msg); blocked {
					if err := p.toClient.Write(mcp.NewError(msg.ID, CodeToolBlocked, reason)); err != nil {
						return err
					}
					continue
				}
			}
		}
		if err := p.toUpstream.WriteRaw(line); err != nil {
			return err
		}
	}
}

func (p *Proxy) pumpUpstream(out io.Reader) error {
	reader := mcp.NewReader(out)
	for {
		line, err := reader.ReadLine()
		if err != nil {
			return err
		}
		var msg mcp.Message
		if json.Unmarshal(line, &msg) != nil {
			if err := p.toClient.WriteRaw(line); err != nil {
				return err
			}
			continue
		}
		if msg.IsResponse() {
			key := string(msg.ID)
			p.mu.Lock()
			first, isList := p.lists[key]
			internal := p.internal[key]
			delete(p.lists, key)
			delete(p.internal, key)
			p.mu.Unlock()
			if internal {
				p.filter(msg, true)
				continue
			}
			if isList {
				if err := p.toClient.Write(p.filter(msg, first)); err != nil {
					return err
				}
				continue
			}
		}
		if msg.IsNotification() && msg.Method == mcp.MethodToolsListChanged {
			p.mu.Lock()
			p.allowed = map[string]bool{}
			p.mu.Unlock()
			p.logf("upstream announced a tool list change, calls are held until the new definitions are checked")
			if err := p.refresh(); err != nil {
				return err
			}
		}
		if err := p.toClient.WriteRaw(line); err != nil {
			return err
		}
	}
}

func (p *Proxy) blockCall(msg mcp.Message) (bool, string) {
	var params mcp.CallParams
	if err := json.Unmarshal(msg.Params, &params); err != nil {
		return true, "mcp-rug-pull-proxy: invalid tools/call params"
	}
	p.mu.Lock()
	ok := p.allowed[params.Name]
	p.mu.Unlock()
	if ok {
		return false, ""
	}
	reason := fmt.Sprintf("mcp-rug-pull-proxy: tool %q on server %q is not approved or its definition changed after approval, review it in the guard UI", params.Name, p.server)
	p.logf("denied call to %s", params.Name)
	p.guard.Denied(p.server, params.Name, "call blocked, tool is not in the approved set")
	return true, reason
}

func (p *Proxy) refresh() error {
	p.mu.Lock()
	p.seq++
	id := fmt.Sprintf("mcp-rug-pull-proxy-%d", p.seq)
	p.internal[fmt.Sprintf("%q", id)] = true
	p.mu.Unlock()
	req, err := mcp.NewRequest(id, mcp.MethodToolsList, map[string]any{})
	if err != nil {
		return err
	}
	return p.toUpstream.Write(req)
}

func (p *Proxy) filter(msg mcp.Message, first bool) mcp.Message {
	if msg.Error != nil || len(msg.Result) == 0 {
		return msg
	}
	var result map[string]json.RawMessage
	var rawTools []json.RawMessage
	if json.Unmarshal(msg.Result, &result) != nil || json.Unmarshal(result["tools"], &rawTools) != nil {
		return p.emptyList(msg, result, first)
	}
	tools := make([]mcp.Tool, 0, len(rawTools))
	for _, raw := range rawTools {
		var t mcp.Tool
		if json.Unmarshal(raw, &t) != nil {
			return p.emptyList(msg, result, first)
		}
		tools = append(tools, t)
	}
	decisions, err := p.guard.Check(p.server, tools)
	if err != nil {
		p.logf("guard unreachable, failing closed and exposing no tools: %v", err)
		return p.emptyList(msg, result, first)
	}
	allowedNow := map[string]bool{}
	for _, d := range decisions {
		if d.Decision == store.DecisionAllow {
			allowedNow[d.Tool] = true
		} else {
			p.logf("hiding %s: %s (%s)", d.Tool, d.Decision, d.Reason)
		}
	}
	kept := []mcp.Tool{}
	for _, t := range tools {
		if allowedNow[t.Name] {
			kept = append(kept, t)
		}
	}
	p.setAllowed(allowedNow, first)
	result["tools"], _ = json.Marshal(kept)
	msg.Result, _ = json.Marshal(result)
	return msg
}

func (p *Proxy) emptyList(msg mcp.Message, result map[string]json.RawMessage, first bool) mcp.Message {
	if result == nil {
		result = map[string]json.RawMessage{}
	}
	p.setAllowed(map[string]bool{}, first)
	result["tools"] = json.RawMessage(`[]`)
	msg.Result, _ = json.Marshal(result)
	return msg
}

func (p *Proxy) setAllowed(allowed map[string]bool, replace bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if replace {
		p.allowed = allowed
		return
	}
	for name := range allowed {
		p.allowed[name] = true
	}
}

func (p *Proxy) Allowed() []string {
	p.mu.Lock()
	defer p.mu.Unlock()
	names := make([]string, 0, len(p.allowed))
	for name := range p.allowed {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func (p *Proxy) logf(format string, args ...any) {
	if p.logw == nil {
		return
	}
	fmt.Fprintf(p.logw, "[mcp-rug-pull-proxy %s] %s\n", p.server, strings.TrimSpace(fmt.Sprintf(format, args...)))
}

func hasCursor(params json.RawMessage) bool {
	var v struct {
		Cursor string `json:"cursor"`
	}
	return json.Unmarshal(params, &v) == nil && v.Cursor != ""
}
