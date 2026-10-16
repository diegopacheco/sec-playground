package session

import (
	"testing"
	"time"
)

func TestCeremonyBoundToAppAndKind(t *testing.T) {
	m := NewManager(time.Minute, time.Hour)
	id := m.StartCeremony(Ceremony{App: "bookstore", Kind: "register"})
	if _, ok := m.TakeCeremony(id, "vault", "register"); ok {
		t.Fatalf("a bookstore challenge must not be finished on the vault")
	}
	if _, ok := m.TakeCeremony(id, "bookstore", "register"); ok {
		t.Fatalf("a failed attempt must burn the ceremony so it cannot be retried")
	}
}

func TestExpiredCeremonyAndSessionAreRejected(t *testing.T) {
	m := NewManager(time.Minute, time.Hour)
	now := time.Now()
	m.now = func() time.Time { return now }
	id := m.StartCeremony(Ceremony{App: "a", Kind: "login"})
	token, _ := m.Create("a", "alice", "x")
	now = now.Add(2 * time.Hour)
	if _, ok := m.TakeCeremony(id, "a", "login"); ok {
		t.Fatalf("expired challenge must be rejected")
	}
	if _, ok := m.Lookup(token); ok {
		t.Fatalf("expired session must be rejected")
	}
}

func TestTokensAreStoredHashed(t *testing.T) {
	m := NewManager(time.Minute, time.Hour)
	token, _ := m.Create("a", "alice", "x")
	if _, raw := m.sessions[token]; raw {
		t.Fatalf("a memory dump must not reveal usable bearer tokens")
	}
}
