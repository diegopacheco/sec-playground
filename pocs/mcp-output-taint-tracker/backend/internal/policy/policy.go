package policy

import (
	"encoding/json"
	"regexp"
	"slices"
	"strings"

	"mcp-output-taint-tracker/internal/taint"
)

type SourceRule struct {
	Tool    string        `json:"tool"`
	Arg     string        `json:"arg,omitempty"`
	Pattern string        `json:"pattern,omitempty"`
	Labels  []taint.Label `json:"labels"`
	re      *regexp.Regexp
}

type SinkRule struct {
	Tool   string        `json:"tool"`
	Arg    string        `json:"arg"`
	Forbid []taint.Label `json:"forbid"`
	Why    string        `json:"why"`
}

type Policy struct {
	Sources []SourceRule `json:"sources"`
	Sinks   []SinkRule   `json:"sinks"`
}

func New(sources []SourceRule, sinks []SinkRule) *Policy {
	for i := range sources {
		if sources[i].Pattern != "" {
			sources[i].re = regexp.MustCompile(sources[i].Pattern)
		}
	}
	return &Policy{sources, sinks}
}

func (p *Policy) SourceLabels(tool string, args map[string]any) []taint.Label {
	var out []taint.Label
	for _, r := range p.Sources {
		if !matchTool(r.Tool, tool) {
			continue
		}
		if r.re != nil {
			v, ok := ArgString(args, r.Arg)
			if !ok || !r.re.MatchString(v) {
				continue
			}
		}
		for _, l := range r.Labels {
			if !slices.Contains(out, l) {
				out = append(out, l)
			}
		}
	}
	return out
}

func (p *Policy) SinkArgs(tool string, args map[string]any) []taint.SinkArg {
	var out []taint.SinkArg
	for _, r := range p.Sinks {
		if !matchTool(r.Tool, tool) {
			continue
		}
		if v, ok := ArgString(args, r.Arg); ok {
			out = append(out, taint.SinkArg{Arg: r.Arg, Value: v, Forbid: r.Forbid})
		}
	}
	return out
}

func ArgString(args map[string]any, key string) (string, bool) {
	v, ok := args[key]
	if !ok || v == nil {
		return "", false
	}
	if s, ok := v.(string); ok {
		return s, true
	}
	b, err := json.Marshal(v)
	if err != nil {
		return "", false
	}
	return string(b), true
}

func matchTool(pattern, tool string) bool {
	if prefix, ok := strings.CutSuffix(pattern, "*"); ok {
		return strings.HasPrefix(tool, prefix)
	}
	return pattern == tool
}
