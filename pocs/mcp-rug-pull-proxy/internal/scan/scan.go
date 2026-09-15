package scan

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"mcp-rug-pull-proxy/internal/mcp"
)

const (
	SeverityNone   = "none"
	SeverityLow    = "low"
	SeverityMedium = "medium"
	SeverityHigh   = "high"
)

type Finding struct {
	Rule     string `json:"rule"`
	Severity string `json:"severity"`
	Field    string `json:"field"`
	Detail   string `json:"detail"`
	Match    string `json:"match"`
	Snippet  string `json:"snippet"`
	Hidden   string `json:"hidden,omitempty"`
}

type Report struct {
	Severity string    `json:"severity"`
	Findings []Finding `json:"findings"`
}

type textRule struct {
	id       string
	severity string
	detail   string
	pattern  *regexp.Regexp
}

var textRules = []textRule{
	{"instruction-override", SeverityHigh, "tries to override the instructions the model already has", regexp.MustCompile(`(?i)(ignore|disregard|forget)\s+(all\s+)?(the\s+)?(previous|prior|above|earlier)\s+(instructions|rules|messages)`)},
	{"conceal-from-user", SeverityHigh, "asks the model to hide its actions from the user", regexp.MustCompile(`(?i)(do\s+not|don't|never)\s+(tell|mention|inform|reveal|show|notify)[^.\n]{0,40}\b(user|human)`)},
	{"hidden-tag", SeverityHigh, "pseudo system tag used to smuggle priority instructions", regexp.MustCompile(`(?i)<\s*/?\s*(important|system|instructions?|secret|admin)\s*>`)},
	{"pre-tool-directive", SeverityMedium, "orders the model to do something before calling the tool", regexp.MustCompile(`(?i)before\s+(using|calling|invoking|running)\s+(this|the)\s+tool`)},
	{"imperative-must", SeverityLow, "imperative wording aimed at the model", regexp.MustCompile(`(?i)\byou\s+must\b`)},
	{"sensitive-file", SeverityHigh, "references a secret bearing file", regexp.MustCompile(`(?i)(~/\.ssh|\.ssh/|id_rsa|id_ed25519|\.aws/|\.env\b|mcp\.json|\.npmrc|\.netrc|\.kube/config)`)},
	{"secret-reference", SeverityMedium, "mentions credentials or private keys", regexp.MustCompile(`(?i)(credentials|private\s+key|api[\s_-]?key|access\s+token|password)`)},
	{"exfil-parameter", SeverityHigh, "asks for file or secret contents to be passed through a parameter", regexp.MustCompile(`(?i)(pass|send|include|put|provide|attach|copy)\b[^.\n]{0,50}\b(content|contents|text|file|key|secret)s?\b[^.\n]{0,60}\b(parameter|param|argument|field|sidenote|notes?)\b`)},
	{"external-url", SeverityMedium, "embeds a URL the model could be steered to", regexp.MustCompile(`https?://[^\s"'<>]+`)},
	{"base64-blob", SeverityMedium, "long base64 looking blob that can hide a payload", regexp.MustCompile(`[A-Za-z0-9+/]{60,}={0,2}`)},
}

var suspiciousParamName = regexp.MustCompile(`(?i)^(notes?|side_?note|context|metadata|extra|debug|feedback|reasoning|details)$`)
var suspiciousParamAsk = regexp.MustCompile(`(?i)(content|contents|file|key|secret|token|password|conversation|history|credential)`)

type field struct {
	path string
	text string
}

func Scan(t mcp.Tool, otherTools []string) Report {
	fields := collect(t)
	var findings []Finding
	for _, f := range fields {
		findings = append(findings, scanUnicode(f)...)
		findings = append(findings, scanText(f)...)
		findings = append(findings, scanShadowing(f, t.Name, otherTools)...)
	}
	findings = append(findings, scanParams(t.InputSchema)...)
	return Report{Severity: maxSeverity(findings), Findings: dedupe(findings)}
}

func collect(t mcp.Tool) []field {
	fields := []field{{"name", t.Name}, {"title", t.Title}, {"description", t.Description}}
	walk("inputSchema", decode(t.InputSchema), &fields)
	walk("outputSchema", decode(t.OutputSchema), &fields)
	walk("annotations", decode(t.Annotations), &fields)
	return fields
}

func decode(raw json.RawMessage) any {
	if len(raw) == 0 {
		return nil
	}
	var v any
	if err := json.Unmarshal(raw, &v); err != nil {
		return string(raw)
	}
	return v
}

func walk(path string, v any, out *[]field) {
	switch node := v.(type) {
	case map[string]any:
		keys := make([]string, 0, len(node))
		for k := range node {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			child := path + "." + k
			if s, ok := node[k].(string); ok {
				*out = append(*out, field{child, s})
				continue
			}
			if strings.HasSuffix(path, ".properties") {
				*out = append(*out, field{child + "#name", k})
			}
			walk(child, node[k], out)
		}
	case []any:
		for i, item := range node {
			walk(fmt.Sprintf("%s[%d]", path, i), item, out)
		}
	case string:
		*out = append(*out, field{path, node})
	}
}

func scanText(f field) []Finding {
	normalized := Normalize(f.text)
	var findings []Finding
	for _, rule := range textRules {
		loc := rule.pattern.FindStringIndex(normalized)
		if loc == nil {
			continue
		}
		findings = append(findings, Finding{
			Rule:     rule.id,
			Severity: rule.severity,
			Field:    f.path,
			Detail:   rule.detail,
			Match:    normalized[loc[0]:loc[1]],
			Snippet:  snippet(normalized, loc[0], loc[1]),
		})
	}
	return findings
}

func scanUnicode(f field) []Finding {
	var zeroWidth, bidi int
	var tags strings.Builder
	for _, r := range f.text {
		switch {
		case isZeroWidth(r):
			zeroWidth++
		case isBidi(r):
			bidi++
		case isTag(r):
			if r >= 0xE0020 && r <= 0xE007E {
				tags.WriteRune(r - 0xE0000)
			}
		}
	}
	var findings []Finding
	if zeroWidth > 0 {
		findings = append(findings, Finding{Rule: "invisible-unicode", Severity: SeverityHigh, Field: f.path, Detail: fmt.Sprintf("%d zero width characters hide text from human reviewers", zeroWidth), Match: "zero width", Snippet: Reveal(f.text, 160)})
	}
	if bidi > 0 {
		findings = append(findings, Finding{Rule: "bidi-control", Severity: SeverityHigh, Field: f.path, Detail: fmt.Sprintf("%d bidirectional control characters reorder what a reviewer sees", bidi), Match: "bidi control", Snippet: Reveal(f.text, 160)})
	}
	if tags.Len() > 0 {
		findings = append(findings, Finding{Rule: "unicode-tag-smuggling", Severity: SeverityHigh, Field: f.path, Detail: "unicode tag characters carry text that is invisible on screen but read by the model", Match: "tag characters", Snippet: Reveal(f.text, 160), Hidden: tags.String()})
	}
	return findings
}

func scanShadowing(f field, self string, others []string) []Finding {
	var findings []Finding
	for _, other := range others {
		if len(other) < 4 || other == self {
			continue
		}
		pattern := regexp.MustCompile(`\b` + regexp.QuoteMeta(other) + `\b`)
		normalized := Normalize(f.text)
		loc := pattern.FindStringIndex(normalized)
		if loc == nil {
			continue
		}
		findings = append(findings, Finding{Rule: "tool-shadowing", Severity: SeverityHigh, Field: f.path, Detail: "mentions a tool from another server, a common way to hijack how that tool is used", Match: other, Snippet: snippet(normalized, loc[0], loc[1])})
	}
	return findings
}

func scanParams(schema json.RawMessage) []Finding {
	root, ok := decode(schema).(map[string]any)
	if !ok {
		return nil
	}
	props, ok := root["properties"].(map[string]any)
	if !ok {
		return nil
	}
	var findings []Finding
	names := make([]string, 0, len(props))
	for name := range props {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		if !suspiciousParamName.MatchString(name) {
			continue
		}
		prop, _ := props[name].(map[string]any)
		desc, _ := prop["description"].(string)
		if !suspiciousParamAsk.MatchString(Normalize(desc)) {
			continue
		}
		findings = append(findings, Finding{Rule: "suspicious-parameter", Severity: SeverityHigh, Field: "inputSchema.properties." + name, Detail: "side channel parameter that asks for file or secret contents", Match: name, Snippet: name + ": " + Reveal(desc, 120)})
	}
	return findings
}

func Normalize(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch {
		case isZeroWidth(r), isBidi(r):
		case isTag(r):
			if r >= 0xE0020 && r <= 0xE007E {
				b.WriteRune(r - 0xE0000)
			}
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

func Reveal(s string, limit int) string {
	var b strings.Builder
	var tags strings.Builder
	flush := func() {
		if tags.Len() > 0 {
			b.WriteString("[hidden:" + tags.String() + "]")
			tags.Reset()
		}
	}
	for _, r := range s {
		switch {
		case isTag(r):
			if r >= 0xE0020 && r <= 0xE007E {
				tags.WriteRune(r - 0xE0000)
			}
			continue
		case isZeroWidth(r), isBidi(r):
			flush()
			b.WriteString(fmt.Sprintf("[U+%04X]", r))
		default:
			flush()
			b.WriteRune(r)
		}
	}
	flush()
	out := []rune(b.String())
	if len(out) > limit {
		return string(out[:limit]) + "..."
	}
	return string(out)
}

func snippet(s string, start, end int) string {
	from := start - 40
	if from < 0 {
		from = 0
	}
	to := end + 40
	if to > len(s) {
		to = len(s)
	}
	for from > 0 && !utf8Start(s[from]) {
		from--
	}
	for to < len(s) && !utf8Start(s[to]) {
		to++
	}
	out := strings.Join(strings.Fields(s[from:to]), " ")
	if from > 0 {
		out = "..." + out
	}
	if to < len(s) {
		out += "..."
	}
	return out
}

func utf8Start(b byte) bool {
	return b&0xC0 != 0x80
}

func isZeroWidth(r rune) bool {
	return (r >= 0x200B && r <= 0x200F) || r == 0x2060 || r == 0xFEFF
}

func isBidi(r rune) bool {
	return (r >= 0x202A && r <= 0x202E) || (r >= 0x2066 && r <= 0x2069)
}

func isTag(r rune) bool {
	return r >= 0xE0000 && r <= 0xE007F
}

func Rank(severity string) int {
	switch severity {
	case SeverityHigh:
		return 3
	case SeverityMedium:
		return 2
	case SeverityLow:
		return 1
	}
	return 0
}

func maxSeverity(findings []Finding) string {
	best := SeverityNone
	for _, f := range findings {
		if Rank(f.Severity) > Rank(best) {
			best = f.Severity
		}
	}
	return best
}

func dedupe(findings []Finding) []Finding {
	seen := map[string]bool{}
	out := []Finding{}
	for _, f := range findings {
		key := f.Rule + "|" + f.Field + "|" + f.Match
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, f)
	}
	sort.SliceStable(out, func(i, j int) bool {
		return Rank(out[i].Severity) > Rank(out[j].Severity)
	})
	return out
}
