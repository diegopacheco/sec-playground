package redact

import (
	"encoding/base64"
	"encoding/hex"
	"net/url"
	"strings"
	"testing"
)

func TestEncodedFormsOfTheCredentialAreRedactedSoEncodingCannotSmuggleItOut(t *testing.T) {
	value := "sk_live+abc/123=="
	leaks := []string{
		value,
		base64.StdEncoding.EncodeToString([]byte(value)),
		base64.RawURLEncoding.EncodeToString([]byte(value)),
		url.QueryEscape(value),
		hex.EncodeToString([]byte(value)),
		strings.ToUpper(hex.EncodeToString([]byte(value))),
	}
	for _, leak := range leaks {
		out, n := Apply(`{"echo":"`+leak+`"}`, value, "vendor-api")
		if n == 0 || strings.Contains(out, leak) {
			t.Errorf("form %q leaked: %s", leak, out)
		}
		if !strings.Contains(out, "[REDACTED:vendor-api]") {
			t.Errorf("marker missing for %q", leak)
		}
	}
}

func TestTextWithoutTheCredentialIsUntouched(t *testing.T) {
	in := `{"account":"acme-support-bot"}`
	out, n := Apply(in, "sk_live_abc", "vendor-api")
	if n != 0 || out != in {
		t.Fatalf("unexpected change: %s", out)
	}
}
