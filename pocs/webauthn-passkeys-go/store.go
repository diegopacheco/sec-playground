package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-webauthn/webauthn/webauthn"
	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

type User struct {
	ID          []byte
	Name        string
	DisplayName string
	Credentials []webauthn.Credential
}

func (u *User) WebAuthnID() []byte {
	return u.ID
}

func (u *User) WebAuthnName() string {
	return u.Name
}

func (u *User) WebAuthnDisplayName() string {
	return u.DisplayName
}

func (u *User) WebAuthnCredentials() []webauthn.Credential {
	return u.Credentials
}

func OpenStore(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	schema := `
CREATE TABLE IF NOT EXISTS users (
  id BLOB PRIMARY KEY,
  name TEXT NOT NULL UNIQUE,
  display_name TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS credentials (
  id BLOB PRIMARY KEY,
  user_id BLOB NOT NULL,
  data BLOB NOT NULL,
  FOREIGN KEY(user_id) REFERENCES users(id)
);
CREATE TABLE IF NOT EXISTS sessions (
  id TEXT PRIMARY KEY,
  kind TEXT NOT NULL,
  user_id BLOB NOT NULL,
  data BLOB NOT NULL,
  expires INTEGER NOT NULL,
  used INTEGER NOT NULL DEFAULT 0
);`
	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, err
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) GetOrCreateUser(name, displayName string) (*User, error) {
	user, err := s.GetUser(name)
	if err == nil {
		return user, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	id, err := randomID(32)
	if err != nil {
		return nil, err
	}
	if _, err := s.db.Exec("INSERT INTO users(id, name, display_name) VALUES(?, ?, ?)", id, name, displayName); err != nil {
		return nil, err
	}
	return &User{ID: id, Name: name, DisplayName: displayName}, nil
}

func (s *Store) GetUser(name string) (*User, error) {
	user := &User{}
	if err := s.db.QueryRow("SELECT id, name, display_name FROM users WHERE name = ?", name).Scan(&user.ID, &user.Name, &user.DisplayName); err != nil {
		return nil, err
	}
	credentials, err := s.credentials(user.ID)
	if err != nil {
		return nil, err
	}
	user.Credentials = credentials
	return user, nil
}

func (s *Store) GetUserByID(id []byte) (*User, error) {
	user := &User{}
	if err := s.db.QueryRow("SELECT id, name, display_name FROM users WHERE id = ?", id).Scan(&user.ID, &user.Name, &user.DisplayName); err != nil {
		return nil, err
	}
	credentials, err := s.credentials(user.ID)
	if err != nil {
		return nil, err
	}
	user.Credentials = credentials
	return user, nil
}

func (s *Store) credentials(userID []byte) ([]webauthn.Credential, error) {
	rows, err := s.db.Query("SELECT data FROM credentials WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	credentials := []webauthn.Credential{}
	for rows.Next() {
		var data []byte
		var credential webauthn.Credential
		if err := rows.Scan(&data); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(data, &credential); err != nil {
			return nil, err
		}
		credentials = append(credentials, credential)
	}
	return credentials, rows.Err()
}

func (s *Store) SaveCredential(userID []byte, credential *webauthn.Credential) error {
	data, err := json.Marshal(credential)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(`INSERT INTO credentials(id, user_id, data) VALUES(?, ?, ?)
ON CONFLICT(id) DO UPDATE SET user_id = excluded.user_id, data = excluded.data`, credential.ID, userID, data)
	return err
}

func (s *Store) SaveSession(kind string, userID []byte, session *webauthn.SessionData) (string, error) {
	id, err := randomToken()
	if err != nil {
		return "", err
	}
	storedSession := *session
	if storedSession.Expires.IsZero() {
		storedSession.Expires = time.Now().Add(5 * time.Minute)
	}
	data, err := json.Marshal(&storedSession)
	if err != nil {
		return "", err
	}
	_, err = s.db.Exec("INSERT INTO sessions(id, kind, user_id, data, expires) VALUES(?, ?, ?, ?, ?)", id, kind, userID, data, storedSession.Expires.Unix())
	return id, err
}

func (s *Store) ConsumeSession(id, kind string) (*webauthn.SessionData, []byte, error) {
	if id == "" {
		return nil, nil, errors.New("missing ceremony session")
	}
	tx, err := s.db.Begin()
	if err != nil {
		return nil, nil, err
	}
	defer tx.Rollback()
	result, err := tx.Exec("UPDATE sessions SET used = 1 WHERE id = ? AND kind = ? AND used = 0 AND expires > ?", id, kind, time.Now().Unix())
	if err != nil {
		return nil, nil, err
	}
	changed, err := result.RowsAffected()
	if err != nil {
		return nil, nil, err
	}
	if changed != 1 {
		return nil, nil, errors.New("ceremony session is invalid, expired, or already used")
	}
	var data []byte
	var userID []byte
	if err := tx.QueryRow("SELECT data, user_id FROM sessions WHERE id = ?", id).Scan(&data, &userID); err != nil {
		return nil, nil, err
	}
	var session webauthn.SessionData
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, nil, err
	}
	return &session, userID, nil
}

func (s *Store) Counts() (int, int, error) {
	var users int
	var credentials int
	if err := s.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&users); err != nil {
		return 0, 0, err
	}
	if err := s.db.QueryRow("SELECT COUNT(*) FROM credentials").Scan(&credentials); err != nil {
		return 0, 0, err
	}
	return users, credentials, nil
}
