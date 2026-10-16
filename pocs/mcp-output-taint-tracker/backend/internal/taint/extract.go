package taint

import (
	"net/url"
	"regexp"
	"strings"
)

const minEntityLen = 6

type Entity struct {
	Kind  string `json:"kind"`
	Value string `json:"value"`
}

var (
	urlRe   = regexp.MustCompile("https?://[^\\s\"'<>()\\[\\]{}`]+")
	emailRe = regexp.MustCompile(`[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}`)
	tokenRe = regexp.MustCompile(`[A-Za-z0-9_\-+/=]{24,}`)
	keyRes  = []*regexp.Regexp{
		regexp.MustCompile(`AKIA[0-9A-Z]{16}`),
		regexp.MustCompile(`gh[pousr]_[A-Za-z0-9]{36}`),
		regexp.MustCompile(`sk-[A-Za-z0-9_\-]{20,}`),
		regexp.MustCompile(`xox[baprs]-[A-Za-z0-9\-]{10,}`),
		regexp.MustCompile(`-----BEGIN [A-Z ]*PRIVATE KEY-----`),
	}
	assignRe = regexp.MustCompile(`(?im)^\s*(?:export\s+)?[A-Z0-9_]*(?:KEY|SECRET|TOKEN|PASSWORD|PASSWD|PWD)[A-Z0-9_]*\s*[=:]\s*["']?([^\s"']{6,})`)
	credURL  = regexp.MustCompile(`://[^/\s:@]+:([^@\s/]{6,})@`)
)

type found struct {
	Entity
	secret bool
}

func extract(text string) []found {
	seen := map[string]bool{}
	var out []found
	add := func(kind, value string, secret bool) {
		if len(value) < minEntityLen {
			return
		}
		key := kind + "|" + strings.ToLower(value)
		if seen[key] {
			return
		}
		seen[key] = true
		out = append(out, found{Entity{kind, value}, secret})
	}
	for _, re := range keyRes {
		for _, m := range re.FindAllString(text, -1) {
			add("secret", m, true)
		}
	}
	for _, m := range assignRe.FindAllStringSubmatch(text, -1) {
		add("secret", m[1], true)
	}
	for _, m := range credURL.FindAllStringSubmatch(text, -1) {
		add("secret", m[1], true)
	}
	for _, m := range urlRe.FindAllString(text, -1) {
		m = strings.TrimRight(m, ".,;:!?")
		add("url", m, false)
		if u, err := url.Parse(m); err == nil {
			add("host", u.Hostname(), false)
		}
	}
	for _, m := range emailRe.FindAllString(text, -1) {
		add("email", m, false)
	}
	for _, m := range tokenRe.FindAllString(text, -1) {
		if !strings.Contains(m, "/") {
			add("token", m, false)
		}
	}
	return out
}
