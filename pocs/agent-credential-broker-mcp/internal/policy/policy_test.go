package policy

import (
	"path/filepath"
	"testing"
)

func TestShippedPolicyLoadsWithHostFromEnvironment(t *testing.T) {
	t.Setenv("UPSTREAM_HOST", "localhost:9999")
	p, err := Load(filepath.Join("..", "..", "policy.json"))
	if err != nil {
		t.Fatal(err)
	}
	c, ok := p.Find("vendor-api")
	if !ok || c.Allow[0].Host != "localhost:9999" {
		t.Fatalf("vendor-api host not expanded: %+v", c)
	}
}

func TestInvalidPoliciesAreRefusedSoAMistakeCannotSilentlyWidenAccess(t *testing.T) {
	bad := map[string]string{
		"no allow rules":   `{"credentials":[{"name":"a","provider":"vault","approval":"human","inject":{"header":"Authorization","template":"Bearer {{credential}}"}}]}`,
		"unknown approval": `{"credentials":[{"name":"a","provider":"vault","approval":"yolo","inject":{"header":"A","template":"{{credential}}"},"allow":[{}]}]}`,
		"no placeholder":   `{"credentials":[{"name":"a","provider":"vault","approval":"auto","inject":{"header":"A","template":"static"},"allow":[{}]}]}`,
		"duplicate":        `{"credentials":[{"name":"a","provider":"vault","approval":"auto","inject":{"header":"A","template":"{{credential}}"},"allow":[{}]},{"name":"a","provider":"vault","approval":"auto","inject":{"header":"A","template":"{{credential}}"},"allow":[{}]}]}`,
		"bad provider":     `{"credentials":[{"name":"a","provider":"env","approval":"auto","inject":{"header":"A","template":"{{credential}}"},"allow":[{}]}]}`,
	}
	for name, raw := range bad {
		if _, err := Parse([]byte(raw)); err == nil {
			t.Errorf("%s must be rejected", name)
		}
	}
}
