package store

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"

	"passkey-server/internal/webauthn"
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
	ErrCredNotFound = errors.New("credential not found")
)

type User struct {
	ID          []byte                `json:"id"`
	Name        string                `json:"name"`
	Credentials []webauthn.Credential `json:"credentials"`
}

type Store struct {
	mu    sync.Mutex
	path  string
	Users map[string]map[string]*User `json:"users"`
}

func Open(path string) (*Store, error) {
	s := &Store{path: path, Users: map[string]map[string]*User{}}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return s, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, s); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Store) save() error {
	if s.path == "" {
		return nil
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0o700); err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}

func (s *Store) Get(app, name string) (*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	u, ok := s.Users[app][name]
	if !ok {
		return nil, ErrUserNotFound
	}
	cp := *u
	cp.Credentials = append([]webauthn.Credential{}, u.Credentials...)
	return &cp, nil
}

func (s *Store) AddCredential(app string, user *User, cred webauthn.Credential, mustBeNew bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, users := range s.Users {
		for _, u := range users {
			for _, c := range u.Credentials {
				if bytes.Equal(c.ID, cred.ID) {
					return errors.New("credential id already registered")
				}
			}
		}
	}
	if s.Users[app] == nil {
		s.Users[app] = map[string]*User{}
	}
	existing, ok := s.Users[app][user.Name]
	switch {
	case ok && mustBeNew:
		return ErrUserExists
	case ok:
		if !bytes.Equal(existing.ID, user.ID) {
			return errors.New("user handle mismatch")
		}
		existing.Credentials = append(existing.Credentials, cred)
	default:
		s.Users[app][user.Name] = &User{ID: user.ID, Name: user.Name, Credentials: []webauthn.Credential{cred}}
	}
	return s.save()
}

func (s *Store) UpdateSignCount(app, name string, credID []byte, count uint32, backedUp bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	u, ok := s.Users[app][name]
	if !ok {
		return ErrUserNotFound
	}
	for i := range u.Credentials {
		if bytes.Equal(u.Credentials[i].ID, credID) {
			u.Credentials[i].SignCount = count
			u.Credentials[i].BackedUp = backedUp
			return s.save()
		}
	}
	return ErrCredNotFound
}
