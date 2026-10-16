package webauthn

import (
	"errors"
	"strings"
	"testing"

	"passkey-server/internal/cose"
	"passkey-server/internal/softauth"
)

var rp = RelyingParty{ID: "localhost", Name: "Test", Origins: []string{"http://localhost:5173"}}

const origin = "http://localhost:5173"

func register(t *testing.T, a *softauth.Authenticator) *Credential {
	t.Helper()
	ch := NewChallenge()
	cd, att := a.Create(rp.ID, origin, ch)
	cred, err := VerifyRegistration(rp, ch, RegistrationInput{ClientDataJSON: cd, AttestationObject: att})
	if err != nil {
		t.Fatalf("registration should succeed: %v", err)
	}
	return cred
}

func login(a *softauth.Authenticator, cred *Credential, rpID, org string) error {
	ch := NewChallenge()
	cd, ad, sig := a.Get(rpID, org, ch)
	res, err := VerifyAssertion(rp, ch, cred, AssertionInput{ClientDataJSON: cd, AuthenticatorData: ad, Signature: sig})
	if err == nil {
		cred.SignCount = res.SignCount
	}
	return err
}

func TestAllSupportedAlgorithmsRegisterAndLogin(t *testing.T) {
	for _, alg := range []int64{cose.AlgES256, cose.AlgEdDSA, cose.AlgRS256} {
		a := softauth.New(alg)
		cred := register(t, a)
		if cred.Alg != alg {
			t.Fatalf("stored alg %d, want %d", cred.Alg, alg)
		}
		if err := login(a, cred, rp.ID, origin); err != nil {
			t.Fatalf("alg %d login failed: %v", alg, err)
		}
	}
}

func TestPackedSelfAttestationIsVerified(t *testing.T) {
	a := softauth.New(cose.AlgES256)
	a.Format = "packed"
	cred := register(t, a)
	if cred.AttestationFmt != "packed" {
		t.Fatalf("format not recorded")
	}
}

func TestPhishingOriginIsRejected(t *testing.T) {
	a := softauth.New(cose.AlgES256)
	ch := NewChallenge()
	cd, att := a.Create(rp.ID, "https://evil.test", ch)
	_, err := VerifyRegistration(rp, ch, RegistrationInput{ClientDataJSON: cd, AttestationObject: att})
	if err == nil || !strings.Contains(err.Error(), "origin") {
		t.Fatalf("a passkey created on another origin must never be accepted, got %v", err)
	}
	cred := register(t, a)
	if err := login(a, cred, rp.ID, "https://evil.test"); err == nil {
		t.Fatalf("assertion from a phishing origin must be rejected")
	}
}

func TestOtherRelyingPartyIDIsRejected(t *testing.T) {
	a := softauth.New(cose.AlgES256)
	cred := register(t, a)
	if err := login(a, cred, "evil.test", origin); err == nil || !strings.Contains(err.Error(), "rpIdHash") {
		t.Fatalf("assertion scoped to another rpId must be rejected, got %v", err)
	}
}

func TestReplayedOrWrongChallengeIsRejected(t *testing.T) {
	a := softauth.New(cose.AlgES256)
	cred := register(t, a)
	cd, ad, sig := a.Get(rp.ID, origin, NewChallenge())
	_, err := VerifyAssertion(rp, NewChallenge(), cred, AssertionInput{ClientDataJSON: cd, AuthenticatorData: ad, Signature: sig})
	if err == nil || !strings.Contains(err.Error(), "challenge") {
		t.Fatalf("an assertion signed over another challenge is a replay, got %v", err)
	}
}

func TestTamperedSignatureIsRejected(t *testing.T) {
	a := softauth.New(cose.AlgES256)
	cred := register(t, a)
	ch := NewChallenge()
	cd, ad, sig := a.Get(rp.ID, origin, ch)
	sig[len(sig)-1] ^= 0xff
	if _, err := VerifyAssertion(rp, ch, cred, AssertionInput{ClientDataJSON: cd, AuthenticatorData: ad, Signature: sig}); err == nil {
		t.Fatalf("tampered signature must be rejected")
	}
}

func TestAnotherAuthenticatorCannotSignForACredential(t *testing.T) {
	owner := softauth.New(cose.AlgES256)
	cred := register(t, owner)
	attacker := softauth.New(cose.AlgES256)
	attacker.CredID = owner.CredID
	if err := login(attacker, cred, rp.ID, origin); err == nil {
		t.Fatalf("knowing the credential id must not be enough to log in")
	}
}

func TestUserVerificationIsRequired(t *testing.T) {
	a := softauth.New(cose.AlgES256)
	a.Flags = FlagUP
	ch := NewChallenge()
	cd, att := a.Create(rp.ID, origin, ch)
	if _, err := VerifyRegistration(rp, ch, RegistrationInput{ClientDataJSON: cd, AttestationObject: att}); err == nil {
		t.Fatalf("passkeys replace passwords, so UV (biometric or PIN) must be required")
	}
}

func TestClonedAuthenticatorIsDetectedBySignCount(t *testing.T) {
	a := softauth.New(cose.AlgES256)
	cred := register(t, a)
	if err := login(a, cred, rp.ID, origin); err != nil {
		t.Fatal(err)
	}
	if err := login(a, cred, rp.ID, origin); err != nil {
		t.Fatal(err)
	}
	clone := *a
	clone.SignCount = 0
	if err := login(&clone, cred, rp.ID, origin); !errors.Is(err, ErrCloned) {
		t.Fatalf("a counter going backwards means a cloned key, got %v", err)
	}
}

func TestSyncedPasskeysWithZeroSignCountAreAllowed(t *testing.T) {
	a := softauth.New(cose.AlgES256)
	a.Synced = true
	cred := register(t, a)
	for range 3 {
		if err := login(a, cred, rp.ID, origin); err != nil {
			t.Fatalf("synced passkeys report 0 forever and must keep working: %v", err)
		}
	}
}

func TestCreateClientDataCannotBeUsedToLogin(t *testing.T) {
	a := softauth.New(cose.AlgES256)
	cred := register(t, a)
	ch := NewChallenge()
	cd, ad, sig := a.Get(rp.ID, origin, ch)
	wrong := []byte(strings.Replace(string(cd), "webauthn.get", "webauthn.create", 1))
	if _, err := VerifyAssertion(rp, ch, cred, AssertionInput{ClientDataJSON: wrong, AuthenticatorData: ad, Signature: sig}); err == nil {
		t.Fatalf("ceremony type confusion must be rejected")
	}
}

func TestTruncatedAuthenticatorDataIsRejected(t *testing.T) {
	if _, err := ParseAuthenticatorData(make([]byte, 36)); err == nil {
		t.Fatalf("short authData must be rejected")
	}
	raw := make([]byte, 37)
	raw[32] = FlagAT
	if _, err := ParseAuthenticatorData(raw); err == nil {
		t.Fatalf("AT flag without credential data must be rejected")
	}
}
