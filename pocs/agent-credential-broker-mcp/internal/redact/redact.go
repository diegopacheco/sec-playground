package redact

import (
	"encoding/base64"
	"encoding/hex"
	"net/url"
	"strings"
)

func Variants(value string) []string {
	if value == "" {
		return nil
	}
	raw := []byte(value)
	candidates := []string{
		value,
		base64.StdEncoding.EncodeToString(raw),
		base64.RawStdEncoding.EncodeToString(raw),
		base64.URLEncoding.EncodeToString(raw),
		base64.RawURLEncoding.EncodeToString(raw),
		url.QueryEscape(value),
		url.PathEscape(value),
		hex.EncodeToString(raw),
		strings.ToUpper(hex.EncodeToString(raw)),
	}
	seen := map[string]bool{}
	var out []string
	for _, c := range candidates {
		if !seen[c] {
			seen[c] = true
			out = append(out, c)
		}
	}
	return out
}

func Apply(text, value, label string) (string, int) {
	count := 0
	for _, v := range Variants(value) {
		n := strings.Count(text, v)
		if n > 0 {
			count += n
			text = strings.ReplaceAll(text, v, "[REDACTED:"+label+"]")
		}
	}
	return text, count
}
