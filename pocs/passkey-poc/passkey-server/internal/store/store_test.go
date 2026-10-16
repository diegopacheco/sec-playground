package store

import (
	"os"
	"path/filepath"
	"testing"

	"passkey-server/internal/webauthn"
)

func TestCredentialsSurviveRestartAndFileIsPrivate(t *testing.T) {
	path := filepath.Join(t.TempDir(), "passkeys.json")
	s, _ := Open(path)
	cred := webauthn.Credential{ID: []byte{1, 2, 3}, PublicKey: []byte{9}, Alg: -7}
	if err := s.AddCredential("bookstore", &User{ID: []byte{7}, Name: "alice"}, cred, true); err != nil {
		t.Fatal(err)
	}
	if err := s.UpdateSignCount("bookstore", "alice", cred.ID, 5, false); err != nil {
		t.Fatal(err)
	}
	reopened, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	u, err := reopened.Get("bookstore", "alice")
	if err != nil || u.Credentials[0].SignCount != 5 {
		t.Fatalf("sign count must persist or clone detection resets on restart: %v %v", u, err)
	}
	info, _ := os.Stat(path)
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("credential store must be owner only, got %v", info.Mode().Perm())
	}
}

func TestSameCredentialCannotBeRegisteredTwice(t *testing.T) {
	s, _ := Open("")
	cred := webauthn.Credential{ID: []byte{1}}
	s.AddCredential("bookstore", &User{ID: []byte{1}, Name: "alice"}, cred, true)
	if err := s.AddCredential("vault", &User{ID: []byte{2}, Name: "bob"}, cred, true); err == nil {
		t.Fatalf("one credential id must map to one account")
	}
}

func TestGetReturnsACopy(t *testing.T) {
	s, _ := Open("")
	s.AddCredential("bookstore", &User{ID: []byte{1}, Name: "alice"}, webauthn.Credential{ID: []byte{1}, SignCount: 1}, true)
	u, _ := s.Get("bookstore", "alice")
	u.Credentials[0].SignCount = 99
	again, _ := s.Get("bookstore", "alice")
	if again.Credentials[0].SignCount != 1 {
		t.Fatalf("callers must not mutate stored credentials without going through the store")
	}
}
