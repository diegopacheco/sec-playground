package session

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"sync"
	"time"
)

type Ceremony struct {
	App       string
	Kind      string
	Username  string
	UserID    []byte
	Challenge []byte
	NewUser   bool
	Expires   time.Time
}

type Session struct {
	App      string    `json:"app"`
	Username string    `json:"username"`
	UserID   string    `json:"userId"`
	Expires  time.Time `json:"expiresAt"`
}

type Manager struct {
	mu          sync.Mutex
	ceremonies  map[string]Ceremony
	sessions    map[string]Session
	ceremonyTTL time.Duration
	sessionTTL  time.Duration
	now         func() time.Time
}

func NewManager(ceremonyTTL, sessionTTL time.Duration) *Manager {
	return &Manager{
		ceremonies:  map[string]Ceremony{},
		sessions:    map[string]Session{},
		ceremonyTTL: ceremonyTTL,
		sessionTTL:  sessionTTL,
		now:         time.Now,
	}
}

func RandomID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func hash(token string) string {
	h := sha256.Sum256([]byte(token))
	return base64.RawURLEncoding.EncodeToString(h[:])
}

func (m *Manager) StartCeremony(c Ceremony) string {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sweep()
	id := RandomID()
	c.Expires = m.now().Add(m.ceremonyTTL)
	m.ceremonies[id] = c
	return id
}

func (m *Manager) TakeCeremony(id, app, kind string) (Ceremony, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	c, ok := m.ceremonies[id]
	delete(m.ceremonies, id)
	if !ok || c.App != app || c.Kind != kind || m.now().After(c.Expires) {
		return Ceremony{}, false
	}
	return c, true
}

func (m *Manager) Create(app, username, userID string) (string, Session) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sweep()
	token := RandomID()
	s := Session{App: app, Username: username, UserID: userID, Expires: m.now().Add(m.sessionTTL)}
	m.sessions[hash(token)] = s
	return token, s
}

func (m *Manager) Lookup(token string) (Session, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	s, ok := m.sessions[hash(token)]
	if !ok || m.now().After(s.Expires) {
		return Session{}, false
	}
	return s, true
}

func (m *Manager) Revoke(token string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.sessions, hash(token))
}

func (m *Manager) sweep() {
	now := m.now()
	for k, c := range m.ceremonies {
		if now.After(c.Expires) {
			delete(m.ceremonies, k)
		}
	}
	for k, s := range m.sessions {
		if now.After(s.Expires) {
			delete(m.sessions, k)
		}
	}
}
