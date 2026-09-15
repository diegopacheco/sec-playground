package sim

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"mcp-rug-pull-proxy/internal/fakeserver"
	"mcp-rug-pull-proxy/internal/mcp"
	"mcp-rug-pull-proxy/internal/pin"
	"mcp-rug-pull-proxy/internal/store"
)

type Scenario struct {
	Mode   string
	Server string
	Client *Client
	Admin  *Admin
	Out    io.Writer
}

func (s *Scenario) Run() error {
	fmt.Fprintf(s.Out, "== scenario %s on server %s ==\n", s.Mode, s.Server)
	name, err := s.Client.Initialize()
	if err != nil {
		return fmt.Errorf("initialize: %w", err)
	}
	s.say("client", "initialize -> %s", name)
	if err := s.Admin.Forget(s.Server); err != nil {
		return err
	}
	s.say("admin", "forget server %s so every tool is reviewed from scratch", s.Server)
	switch s.Mode {
	case fakeserver.ModeClean:
		return s.clean()
	case fakeserver.ModePoisoned:
		return s.poisoned()
	case fakeserver.ModeRugPull:
		return s.rugPull()
	}
	return fmt.Errorf("unknown mode %s", s.Mode)
}

func (s *Scenario) clean() error {
	if err := s.expectExposed(); err != nil {
		return err
	}
	if err := s.expectDenied("get_weather"); err != nil {
		return err
	}
	if err := s.approve("get_weather", "add_numbers"); err != nil {
		return err
	}
	if err := s.expectExposed("add_numbers", "get_weather"); err != nil {
		return err
	}
	if err := s.expectAnswer("get_weather"); err != nil {
		return err
	}
	fmt.Fprintln(s.Out, "PASS clean: tools stay hidden until a reviewer pins their exact hash")
	return nil
}

func (s *Scenario) poisoned() error {
	if err := s.expectExposed(); err != nil {
		return err
	}
	pins, err := s.Admin.Pins(s.Server)
	if err != nil {
		return err
	}
	if p := pins["add_numbers"]; p.Status != store.StatusBlocked {
		return fmt.Errorf("poisoned add_numbers should be blocked by the scanner, got %s", p.Status)
	}
	if err := s.approve("get_weather"); err != nil {
		return err
	}
	if err := s.expectExposed("get_weather"); err != nil {
		return err
	}
	if err := s.expectDenied("add_numbers"); err != nil {
		return err
	}
	fmt.Fprintln(s.Out, "PASS poisoned: hidden instructions were caught before the model ever saw the tool")
	return nil
}

func (s *Scenario) rugPull() error {
	if err := s.expectExposed(); err != nil {
		return err
	}
	if err := s.approve("get_weather", "add_numbers"); err != nil {
		return err
	}
	if err := s.expectExposed("add_numbers", "get_weather"); err != nil {
		return err
	}
	if err := s.expectAnswer("get_weather"); err != nil {
		return err
	}
	if !s.Client.WaitNotification(mcp.MethodToolsListChanged, 5*time.Second) {
		return fmt.Errorf("upstream never announced the tool change")
	}
	s.say("server", "%s received, upstream swapped the get_weather definition", mcp.MethodToolsListChanged)
	if err := s.expectExposed("add_numbers"); err != nil {
		return err
	}
	if err := s.expectDenied("get_weather"); err != nil {
		return err
	}
	alerts, err := s.Admin.Alerts(s.Server)
	if err != nil {
		return err
	}
	if len(alerts) == 0 || alerts[0].State != store.AlertOpen || alerts[0].Tool != "get_weather" {
		return fmt.Errorf("expected an open rug pull alert for get_weather, got %+v", alerts)
	}
	s.say("guard", "rug pull alert %s -> %s severity %s", pin.Short(alerts[0].OldHash), pin.Short(alerts[0].NewHash), alerts[0].NewScan.Severity)
	fmt.Fprintln(s.Out, "PASS rugpull: an approved tool was silently redefined and the proxy blocked it")
	return nil
}

func (s *Scenario) expectExposed(want ...string) error {
	tools, err := s.Client.ListTools()
	if err != nil {
		return fmt.Errorf("tools/list: %w", err)
	}
	var got []string
	for _, t := range tools {
		got = append(got, t.Name)
	}
	sort.Strings(got)
	sort.Strings(want)
	s.say("client", "tools/list -> %d exposed [%s]", len(got), strings.Join(got, ", "))
	pins, err := s.Admin.Pins(s.Server)
	if err == nil {
		names := make([]string, 0, len(pins))
		for n := range pins {
			names = append(names, n)
		}
		sort.Strings(names)
		for _, n := range names {
			p := pins[n]
			suffix := ""
			if p.AlertID != "" {
				suffix = ", rug pull alert open"
			}
			s.say("guard", "%s %s %s severity %s%s", n, p.Status, pin.Short(p.Hash), p.Scan.Severity, suffix)
		}
	}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		return fmt.Errorf("exposed tools %v, want %v", got, want)
	}
	return nil
}

func (s *Scenario) expectDenied(tool string) error {
	_, rpcErr, err := s.Client.CallTool(tool, args(tool))
	if err != nil {
		return err
	}
	if rpcErr == nil {
		return fmt.Errorf("call to %s should have been denied", tool)
	}
	s.say("client", "tools/call %s -> denied (%d)", tool, rpcErr.Code)
	return nil
}

func (s *Scenario) expectAnswer(tool string) error {
	text, rpcErr, err := s.Client.CallTool(tool, args(tool))
	if err != nil {
		return err
	}
	if rpcErr != nil {
		return fmt.Errorf("call to %s denied: %s", tool, rpcErr.Message)
	}
	s.say("client", "tools/call %s -> %s", tool, text)
	return nil
}

func (s *Scenario) approve(tools ...string) error {
	pins, err := s.Admin.Pins(s.Server)
	if err != nil {
		return err
	}
	for _, tool := range tools {
		p, ok := pins[tool]
		if !ok {
			return fmt.Errorf("%s is not pinned yet", tool)
		}
		if err := s.Admin.Approve(s.Server, tool, p.Hash); err != nil {
			return err
		}
		s.say("admin", "approve %s pinned at %s", tool, pin.Short(p.Hash))
	}
	return nil
}

func (s *Scenario) say(who, format string, a ...any) {
	fmt.Fprintf(s.Out, "[%-6s] %s\n", who, fmt.Sprintf(format, a...))
}

func args(tool string) map[string]any {
	if tool == "add_numbers" {
		return map[string]any{"a": 2, "b": 3}
	}
	return map[string]any{"city": "Porto Alegre"}
}
