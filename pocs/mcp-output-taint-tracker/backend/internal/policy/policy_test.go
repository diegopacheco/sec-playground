package policy

import (
	"slices"
	"testing"

	"mcp-output-taint-tracker/internal/taint"
)

func TestSensitivePathAddsSecretOnTopOfUntrusted(t *testing.T) {
	p := Default()
	got := p.SourceLabels("Read", map[string]any{"file_path": "app/.env"})
	if !slices.Contains(got, taint.Untrusted) || !slices.Contains(got, taint.Secret) {
		t.Fatalf("reading .env must be untrusted and secret, got %v", got)
	}
	if got := p.SourceLabels("Read", map[string]any{"file_path": "README.md"}); slices.Contains(got, taint.Secret) {
		t.Fatalf("README is not a secret, got %v", got)
	}
}

func TestEveryMCPToolOutputIsUntrustedByDefault(t *testing.T) {
	if got := Default().SourceLabels("mcp__anything__tool", nil); !slices.Contains(got, taint.Untrusted) {
		t.Fatalf("unknown MCP server output must be untrusted, got %v", got)
	}
}

func TestSinkArgsOnlyReturnsArgumentsThatArePresent(t *testing.T) {
	got := Default().SinkArgs("mcp__mail__send_email", map[string]any{"to": "a@b.co"})
	if len(got) != 1 || got[0].Arg != "to" {
		t.Fatalf("expected only the to argument, got %+v", got)
	}
}

func TestNonStringArgumentIsCheckedAsJSON(t *testing.T) {
	v, ok := ArgString(map[string]any{"body": map[string]any{"k": "v"}}, "body")
	if !ok || v != `{"k":"v"}` {
		t.Fatalf("nested body must be serialized for checking, got %q", v)
	}
}
