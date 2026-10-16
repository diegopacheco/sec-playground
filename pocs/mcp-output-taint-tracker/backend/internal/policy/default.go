package policy

import "mcp-output-taint-tracker/internal/taint"

var (
	untrusted = []taint.Label{taint.Untrusted}
	secret    = []taint.Label{taint.Secret}
	both      = []taint.Label{taint.Untrusted, taint.Secret}
)

func Default() *Policy {
	return New(
		[]SourceRule{
			{Tool: "WebFetch", Labels: untrusted},
			{Tool: "WebSearch", Labels: untrusted},
			{Tool: "Read", Labels: untrusted},
			{Tool: "Read", Arg: "file_path", Pattern: `(?i)(\.env|\.ssh/|id_rsa|id_ed25519|\.aws/|\.npmrc|\.netrc|\.pem$|\.key$|secret|credential)`, Labels: secret},
			{Tool: "Bash", Labels: untrusted},
			{Tool: "Bash", Arg: "command", Pattern: `(?i)(printenv|\benv\b|cat .*\.env|security find-generic-password|vault (kv )?read)`, Labels: secret},
			{Tool: "mcp__*", Labels: untrusted},
			{Tool: "mcp__db__*", Labels: secret},
		},
		[]SinkRule{
			{Tool: "Bash", Arg: "command", Forbid: both, Why: "a command built from tool output is attacker controlled code, and a command holding a secret can ship it anywhere"},
			{Tool: "WebFetch", Arg: "url", Forbid: secret, Why: "a secret inside a URL leaves the machine in the request line"},
			{Tool: "Write", Arg: "file_path", Forbid: untrusted, Why: "a path chosen by tool output can overwrite shell profiles, hooks or CI files"},
			{Tool: "Edit", Arg: "file_path", Forbid: untrusted, Why: "a path chosen by tool output can overwrite shell profiles, hooks or CI files"},
			{Tool: "mcp__mail__send_email", Arg: "to", Forbid: untrusted, Why: "a recipient taken from an email or a page is how injected text redirects mail"},
			{Tool: "mcp__mail__send_email", Arg: "body", Forbid: secret, Why: "a secret in a mail body is exfiltration"},
			{Tool: "mcp__http__request", Arg: "url", Forbid: both, Why: "a destination chosen by tool output or carrying a secret is exfiltration"},
			{Tool: "mcp__http__request", Arg: "body", Forbid: secret, Why: "a secret in a request body is exfiltration"},
			{Tool: "mcp__github__*", Arg: "body", Forbid: secret, Why: "issues and comments are public, a secret there is leaked"},
		},
	)
}
