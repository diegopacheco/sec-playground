package egress

import (
	"testing"

	"agent-credential-broker-mcp/internal/policy"
)

var cred = policy.Credential{
	Name:  "vendor-api",
	Allow: []policy.Rule{{Scheme: "https", Host: "api.vendor.test", Methods: []string{"GET"}, PathPrefix: "/v1/"}},
}

func TestAllowedRequestPasses(t *testing.T) {
	u, err := Check(cred, "get", "https://API.vendor.test/v1/profile?x=1#frag")
	if err != nil {
		t.Fatal(err)
	}
	if u.String() != "https://api.vendor.test/v1/profile?x=1" {
		t.Fatalf("normalized url must drop fragment and lowercase host, got %s", u.String())
	}
}

func TestHostTricksCannotSendTheCredentialToAnAttackerHost(t *testing.T) {
	attacks := []string{
		"https://api.vendor.test.attacker.test/v1/profile",
		"https://api.vendor.test@attacker.test/v1/profile",
		"https://user:pass@api.vendor.test/v1/profile",
		"https://attacker.test/v1/profile?h=api.vendor.test",
		"https://api.vendor.test:8443/v1/profile",
		"http://api.vendor.test/v1/profile",
		"//api.vendor.test/v1/profile",
	}
	for _, a := range attacks {
		if _, err := Check(cred, "GET", a); err == nil {
			t.Errorf("%s must be rejected", a)
		}
	}
}

func TestPathEscapesCannotLeaveTheAllowedPrefix(t *testing.T) {
	attacks := []string{
		"https://api.vendor.test/v1/../admin/keys",
		"https://api.vendor.test/v1/%2e%2e/admin/keys",
		"https://api.vendor.test/v1%2f..%2fadmin",
		"https://api.vendor.test/admin/keys",
		"https://api.vendor.test/v1",
	}
	for _, a := range attacks {
		if _, err := Check(cred, "GET", a); err == nil {
			t.Errorf("%s must be rejected", a)
		}
	}
}

func TestMethodOutsidePolicyIsRejectedSoReadOnlyCredentialsStayReadOnly(t *testing.T) {
	if _, err := Check(cred, "DELETE", "https://api.vendor.test/v1/profile"); err == nil {
		t.Fatal("DELETE must be rejected")
	}
}
