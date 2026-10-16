package taint

import (
	"encoding/base64"
	"encoding/hex"
	"net/url"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

type variant struct {
	encoding string
	raw      string
	text     string
}

var (
	b64Re = regexp.MustCompile(`[A-Za-z0-9+/_\-]{16,}={0,2}`)
	hexRe = regexp.MustCompile(`[0-9a-fA-F]{16,}`)
)

func variants(s string) []variant {
	base := []variant{{"plain", s, s}}
	if strings.Contains(s, "%") || strings.Contains(s, "+") {
		if d, err := url.QueryUnescape(s); err == nil && d != s {
			base = append(base, variant{"url", s, d})
		}
	}
	out := append([]variant{}, base...)
	for _, v := range base {
		out = append(out, decodeRuns(v)...)
	}
	return out
}

func decodeRuns(v variant) []variant {
	var out []variant
	for _, run := range hexRe.FindAllString(v.text, -1) {
		if len(run)%2 != 0 {
			continue
		}
		if b, err := hex.DecodeString(run); err == nil && printable(b) {
			out = append(out, variant{chain(v.encoding, "hex"), run, string(b)})
		}
	}
	for _, run := range b64Re.FindAllString(v.text, -1) {
		if b, ok := decodeBase64(run); ok {
			out = append(out, variant{chain(v.encoding, "base64"), run, string(b)})
		}
	}
	return out
}

func chain(prev, next string) string {
	if prev == "plain" {
		return next
	}
	return prev + "+" + next
}

func decodeBase64(s string) ([]byte, bool) {
	bare := strings.TrimRight(s, "=")
	for _, enc := range []*base64.Encoding{base64.RawStdEncoding, base64.RawURLEncoding} {
		if b, err := enc.DecodeString(bare); err == nil && len(b) >= minEntityLen && printable(b) {
			return b, true
		}
	}
	return nil, false
}

func printable(b []byte) bool {
	if !utf8.Valid(b) {
		return false
	}
	total, ok := 0, 0
	for _, r := range string(b) {
		total++
		if unicode.IsPrint(r) || unicode.IsSpace(r) {
			ok++
		}
	}
	return total > 0 && ok*10 >= total*9
}
