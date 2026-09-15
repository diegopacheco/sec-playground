package audit

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"sync"
	"time"
)

type Entry struct {
	Seq        int    `json:"seq"`
	Time       string `json:"time"`
	Event      string `json:"event"`
	Outcome    string `json:"outcome"`
	IntentID   string `json:"intent_id,omitempty"`
	Credential string `json:"credential,omitempty"`
	Method     string `json:"method,omitempty"`
	URL        string `json:"url,omitempty"`
	Detail     string `json:"detail,omitempty"`
	Prev       string `json:"prev"`
	Hash       string `json:"hash"`
}

type Log struct {
	mu      sync.Mutex
	entries []Entry
	file    string
	now     func() time.Time
}

func Open(file string, now func() time.Time) (*Log, error) {
	l := &Log{file: file, now: now}
	if file == "" {
		return l, nil
	}
	f, err := os.Open(file)
	if os.IsNotExist(err) {
		return l, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 1<<20), 1<<20)
	for sc.Scan() {
		var e Entry
		if json.Unmarshal(sc.Bytes(), &e) == nil {
			l.entries = append(l.entries, e)
		}
	}
	return l, sc.Err()
}

func digest(e Entry) string {
	e.Hash = ""
	raw, _ := json.Marshal(e)
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:])
}

func (l *Log) Append(e Entry) Entry {
	l.mu.Lock()
	defer l.mu.Unlock()
	e.Seq = len(l.entries) + 1
	e.Time = l.now().UTC().Format(time.RFC3339Nano)
	e.Prev = "genesis"
	if n := len(l.entries); n > 0 {
		e.Prev = l.entries[n-1].Hash
	}
	e.Hash = digest(e)
	l.entries = append(l.entries, e)
	if l.file != "" {
		if f, err := os.OpenFile(l.file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600); err == nil {
			raw, _ := json.Marshal(e)
			f.Write(append(raw, '\n'))
			f.Close()
		}
	}
	return e
}

func (l *Log) Entries() []Entry {
	l.mu.Lock()
	defer l.mu.Unlock()
	return append([]Entry(nil), l.entries...)
}

func Verify(entries []Entry) (bool, int) {
	prev := "genesis"
	for i, e := range entries {
		if e.Prev != prev || digest(e) != e.Hash {
			return false, i + 1
		}
		prev = e.Hash
	}
	return true, 0
}
