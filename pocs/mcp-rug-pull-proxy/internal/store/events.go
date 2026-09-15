package store

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"time"
)

const genesis = "sha256:0000000000000000000000000000000000000000000000000000000000000000"

type Event struct {
	Seq    int       `json:"seq"`
	Time   time.Time `json:"time"`
	Type   string    `json:"type"`
	Server string    `json:"server,omitempty"`
	Tool   string    `json:"tool,omitempty"`
	Hash   string    `json:"toolHash,omitempty"`
	Detail string    `json:"detail"`
	Prev   string    `json:"prev"`
	Sum    string    `json:"sum"`
}

type ChainStatus struct {
	Verified bool   `json:"verified"`
	Count    int    `json:"count"`
	BrokenAt int    `json:"brokenAt,omitempty"`
	Reason   string `json:"reason,omitempty"`
}

func seal(e Event) string {
	e.Sum = ""
	raw, _ := json.Marshal(e)
	sum := sha256.Sum256(append([]byte(e.Prev), raw...))
	return "sha256:" + hex.EncodeToString(sum[:])
}

func Verify(events []Event) ChainStatus {
	prev := genesis
	for i, e := range events {
		if e.Seq != i+1 {
			return ChainStatus{Count: len(events), BrokenAt: i + 1, Reason: "sequence gap"}
		}
		if e.Prev != prev {
			return ChainStatus{Count: len(events), BrokenAt: e.Seq, Reason: "previous hash does not match"}
		}
		if seal(e) != e.Sum {
			return ChainStatus{Count: len(events), BrokenAt: e.Seq, Reason: "event content was modified"}
		}
		prev = e.Sum
	}
	return ChainStatus{Verified: true, Count: len(events)}
}

func loadEvents(path string) ([]Event, error) {
	f, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var events []Event
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1<<16), 1<<22)
	for scanner.Scan() {
		if len(scanner.Bytes()) == 0 {
			continue
		}
		var e Event
		if err := json.Unmarshal(scanner.Bytes(), &e); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, scanner.Err()
}

func appendEvent(path string, e Event) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()
	raw, err := json.Marshal(e)
	if err != nil {
		return err
	}
	_, err = f.Write(append(raw, '\n'))
	return err
}
