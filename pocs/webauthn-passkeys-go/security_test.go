package main

import (
	"crypto/sha256"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

func TestReplaySessionRejected(t *testing.T) {
	store, err := OpenStore(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	session := &webauthn.SessionData{Challenge: "challenge", Expires: time.Now().Add(time.Minute)}
	id, err := store.SaveSession("login", []byte("user"), session)
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err := store.ConsumeSession(id, "login"); err != nil {
		t.Fatal(err)
	}
	if _, _, err := store.ConsumeSession(id, "login"); err == nil {
		t.Fatal("replayed ceremony session was accepted")
	}
}

func TestSessionWithoutExpirationCanBeConsumed(t *testing.T) {
	store, err := OpenStore(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	started := time.Now()
	id, err := store.SaveSession("registration", []byte("user"), &webauthn.SessionData{Challenge: "challenge"})
	if err != nil {
		t.Fatal(err)
	}
	session, _, err := store.ConsumeSession(id, "registration")
	if err != nil {
		t.Fatal(err)
	}
	if !session.Expires.After(started) {
		t.Fatal("ceremony session did not receive an expiration")
	}
	if session.Expires.After(started.Add(5*time.Minute + time.Second)) {
		t.Fatal("ceremony session expiration is too long")
	}
}

func TestIncorrectOriginRejected(t *testing.T) {
	clientData := protocol.CollectedClientData{
		Type:      protocol.AssertCeremony,
		Challenge: "challenge",
		Origin:    "https://attacker.invalid",
	}
	err := clientData.Verify(
		"challenge",
		protocol.AssertCeremony,
		[]string{"http://localhost:8090"},
		nil,
		protocol.TopOriginExplicitVerificationMode,
		false,
	)
	if err == nil {
		t.Fatal("incorrect origin was accepted")
	}
}

func TestWrongRelyingPartyIDRejected(t *testing.T) {
	actual := sha256.Sum256([]byte("localhost"))
	wrong := sha256.Sum256([]byte("attacker.invalid"))
	data := protocol.AuthenticatorData{
		RPIDHash: actual[:],
		Flags:    protocol.FlagUserPresent | protocol.FlagUserVerified,
	}
	if err := data.Verify(wrong[:], nil, true, true); err == nil {
		t.Fatal("wrong relying-party ID was accepted")
	}
}

func TestMissingUserVerificationRejected(t *testing.T) {
	rpID := sha256.Sum256([]byte("localhost"))
	data := protocol.AuthenticatorData{
		RPIDHash: rpID[:],
		Flags:    protocol.FlagUserPresent,
	}
	if err := data.Verify(rpID[:], nil, true, true); err == nil {
		t.Fatal("missing user verification was accepted")
	}
}
