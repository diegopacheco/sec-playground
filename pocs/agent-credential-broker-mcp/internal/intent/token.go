package intent

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"time"
)

type Claims struct {
	ID         string `json:"jti"`
	IntentID   string `json:"intent"`
	Credential string `json:"credential"`
	Method     string `json:"method"`
	URL        string `json:"url"`
	BodyHash   string `json:"body_sha256"`
	ExpiresAt  int64  `json:"exp"`
}

type Signer struct {
	key  []byte
	now  func() time.Time
	mu   sync.Mutex
	used map[string]time.Time
}

var (
	ErrMalformed = errors.New("token is malformed")
	ErrSignature = errors.New("token signature is invalid")
	ErrReplay    = errors.New("token was already used")
	ErrExpired   = errors.New("token is expired")
	ErrMismatch  = errors.New("request does not match the approved intent")
)

func NewSigner(now func() time.Time) *Signer {
	key := make([]byte, 32)
	rand.Read(key)
	return &Signer{key: key, now: now, used: map[string]time.Time{}}
}

func BodyHash(body string) string {
	sum := sha256.Sum256([]byte(body))
	return hex.EncodeToString(sum[:])
}

func RandomID(prefix string) string {
	b := make([]byte, 12)
	rand.Read(b)
	return prefix + hex.EncodeToString(b)
}

func (s *Signer) Mint(c Claims) string {
	payload, _ := json.Marshal(c)
	enc := base64.RawURLEncoding.EncodeToString(payload)
	return enc + "." + base64.RawURLEncoding.EncodeToString(s.sign(enc))
}

func (s *Signer) sign(enc string) []byte {
	m := hmac.New(sha256.New, s.key)
	m.Write([]byte(enc))
	return m.Sum(nil)
}

func (s *Signer) Redeem(token, method, url, body string) (Claims, error) {
	var c Claims
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return c, ErrMalformed
	}
	sig, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return c, ErrMalformed
	}
	if !hmac.Equal(sig, s.sign(parts[0])) {
		return c, ErrSignature
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil || json.Unmarshal(payload, &c) != nil {
		return c, ErrMalformed
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, burned := s.used[c.ID]; burned {
		return c, ErrReplay
	}
	s.used[c.ID] = time.Unix(c.ExpiresAt, 0)
	s.sweep()
	if s.now().Unix() >= c.ExpiresAt {
		return c, ErrExpired
	}
	if c.Method != method || c.URL != url || c.BodyHash != BodyHash(body) {
		return c, ErrMismatch
	}
	return c, nil
}

func (s *Signer) sweep() {
	cutoff := s.now().Add(-time.Minute)
	for id, exp := range s.used {
		if exp.Before(cutoff) {
			delete(s.used, id)
		}
	}
}
