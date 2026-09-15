package audit

import (
	"path/filepath"
	"testing"
	"time"
)

func fill(l *Log) {
	for _, ev := range []string{"request", "approve", "call"} {
		l.Append(Entry{Event: ev, Outcome: "allowed", Credential: "vendor-api"})
	}
}

func TestEditingAnEntryBreaksTheChainSoHistoryCannotBeRewrittenSilently(t *testing.T) {
	l, _ := Open("", time.Now)
	fill(l)
	entries := l.Entries()
	if ok, _ := Verify(entries); !ok {
		t.Fatal("fresh chain must verify")
	}
	entries[1].Outcome = "denied"
	if ok, at := Verify(entries); ok || at != 2 {
		t.Fatalf("edited entry must break chain at 2, got ok=%v at=%d", ok, at)
	}
}

func TestDeletingAnEntryBreaksTheChain(t *testing.T) {
	l, _ := Open("", time.Now)
	fill(l)
	entries := l.Entries()
	cut := append([]Entry{entries[0]}, entries[2])
	if ok, at := Verify(cut); ok || at != 2 {
		t.Fatalf("deleted entry must break chain at 2, got ok=%v at=%d", ok, at)
	}
}

func TestChainSurvivesRestart(t *testing.T) {
	file := filepath.Join(t.TempDir(), "audit.jsonl")
	l, _ := Open(file, time.Now)
	fill(l)
	again, err := Open(file, time.Now)
	if err != nil {
		t.Fatal(err)
	}
	again.Append(Entry{Event: "call", Outcome: "denied"})
	entries := again.Entries()
	if len(entries) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(entries))
	}
	if ok, at := Verify(entries); !ok {
		t.Fatalf("chain must continue after reload, broken at %d", at)
	}
}
