package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"regexp"
	"strings"

	"passkey-server/internal/session"
	"passkey-server/internal/store"
	"passkey-server/internal/webauthn"
)

const maxBody = 64 << 10

var usernamePattern = regexp.MustCompile(`^[a-zA-Z0-9._-]{3,32}$`)

type Server struct {
	apps     map[string]webauthn.RelyingParty
	store    *store.Store
	sessions *session.Manager
}

func New(apps map[string]webauthn.RelyingParty, st *store.Store, sm *session.Manager) *Server {
	return &Server{apps: apps, store: st, sessions: sm}
}

type B64 []byte

func (b *B64) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	raw, err := base64.RawURLEncoding.DecodeString(strings.TrimRight(s, "="))
	if err != nil {
		return errors.New("invalid base64url")
	}
	*b = raw
	return nil
}

type beginRequest struct {
	Username string `json:"username"`
}

type credentialResponse struct {
	ClientDataJSON    B64      `json:"clientDataJSON"`
	AttestationObject B64      `json:"attestationObject"`
	Transports        []string `json:"transports"`
	AuthenticatorData B64      `json:"authenticatorData"`
	Signature         B64      `json:"signature"`
	UserHandle        B64      `json:"userHandle"`
}

type finishRequest struct {
	CeremonyID string `json:"ceremonyId"`
	Credential struct {
		RawID    B64                `json:"rawId"`
		Type     string             `json:"type"`
		Response credentialResponse `json:"response"`
	} `json:"credential"`
}

type passkeyView struct {
	ID             string   `json:"id"`
	Alg            int64    `json:"alg"`
	SignCount      uint32   `json:"signCount"`
	AttestationFmt string   `json:"attestationFmt"`
	BackupEligible bool     `json:"backupEligible"`
	BackedUp       bool     `json:"backedUp"`
	Transports     []string `json:"transports"`
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("POST /api/{app}/register/begin", s.withApp(s.registerBegin))
	mux.HandleFunc("POST /api/{app}/register/finish", s.withApp(s.registerFinish))
	mux.HandleFunc("POST /api/{app}/login/begin", s.withApp(s.loginBegin))
	mux.HandleFunc("POST /api/{app}/login/finish", s.withApp(s.loginFinish))
	mux.HandleFunc("GET /api/{app}/passkeys", s.withApp(s.passkeys))
	mux.HandleFunc("GET /api/session", s.session)
	mux.HandleFunc("POST /api/logout", s.logout)
	return localOnly(mux)
}

func localOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host, _, err := net.SplitHostPort(r.Host)
		if err != nil {
			host = r.Host
		}
		if host != "localhost" && host != "127.0.0.1" {
			writeError(w, http.StatusForbidden, "host not allowed")
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, maxBody)
		next.ServeHTTP(w, r)
	})
}

func (s *Server) withApp(h func(http.ResponseWriter, *http.Request, string, webauthn.RelyingParty)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		app := r.PathValue("app")
		rp, ok := s.apps[app]
		if !ok {
			writeError(w, http.StatusNotFound, "unknown app")
			return
		}
		h(w, r, app, rp)
	}
}

func (s *Server) registerBegin(w http.ResponseWriter, r *http.Request, app string, rp webauthn.RelyingParty) {
	var req beginRequest
	if !readJSON(w, r, &req) {
		return
	}
	if !usernamePattern.MatchString(req.Username) {
		writeError(w, http.StatusBadRequest, "username must be 3-32 chars of letters, digits, dot, dash or underscore")
		return
	}
	if _, err := s.store.Get(app, req.Username); err == nil {
		writeError(w, http.StatusConflict, "username already taken")
		return
	}
	userID := webauthn.NewChallenge()
	challenge := webauthn.NewChallenge()
	id := s.sessions.StartCeremony(session.Ceremony{App: app, Kind: "register", Username: req.Username, UserID: userID, Challenge: challenge, NewUser: true})
	writeJSON(w, http.StatusOK, map[string]any{
		"ceremonyId": id,
		"publicKey":  webauthn.NewCreationOptions(rp, userID, req.Username, challenge, nil),
	})
}

func (s *Server) registerFinish(w http.ResponseWriter, r *http.Request, app string, rp webauthn.RelyingParty) {
	var req finishRequest
	if !readJSON(w, r, &req) {
		return
	}
	c, ok := s.sessions.TakeCeremony(req.CeremonyID, app, "register")
	if !ok {
		writeError(w, http.StatusBadRequest, "unknown or expired ceremony")
		return
	}
	resp := req.Credential.Response
	cred, err := webauthn.VerifyRegistration(rp, c.Challenge, webauthn.RegistrationInput{
		ClientDataJSON:    resp.ClientDataJSON,
		AttestationObject: resp.AttestationObject,
		Transports:        resp.Transports,
	})
	if err != nil {
		log.Printf("register %s/%s rejected: %v", app, c.Username, err)
		writeError(w, http.StatusUnauthorized, "registration rejected: "+err.Error())
		return
	}
	if req.Credential.RawID != nil && !bytes.Equal(req.Credential.RawID, cred.ID) {
		writeError(w, http.StatusBadRequest, "rawId does not match attested credential id")
		return
	}
	user := &store.User{ID: c.UserID, Name: c.Username}
	if err := s.store.AddCredential(app, user, *cred, c.NewUser); err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	log.Printf("register %s/%s ok alg=%d fmt=%s", app, c.Username, cred.Alg, cred.AttestationFmt)
	s.issue(w, app, user)
}

func (s *Server) loginBegin(w http.ResponseWriter, r *http.Request, app string, rp webauthn.RelyingParty) {
	var req beginRequest
	if !readJSON(w, r, &req) {
		return
	}
	user, err := s.store.Get(app, req.Username)
	if err != nil {
		writeError(w, http.StatusNotFound, "no passkey registered for this username")
		return
	}
	allow := make([]webauthn.Descriptor, 0, len(user.Credentials))
	for _, c := range user.Credentials {
		allow = append(allow, webauthn.Descriptor{Type: "public-key", ID: webauthn.B64(c.ID), Transports: c.Transports})
	}
	challenge := webauthn.NewChallenge()
	id := s.sessions.StartCeremony(session.Ceremony{App: app, Kind: "login", Username: user.Name, UserID: user.ID, Challenge: challenge})
	writeJSON(w, http.StatusOK, map[string]any{
		"ceremonyId": id,
		"publicKey":  webauthn.NewRequestOptions(rp, challenge, allow),
	})
}

func (s *Server) loginFinish(w http.ResponseWriter, r *http.Request, app string, rp webauthn.RelyingParty) {
	var req finishRequest
	if !readJSON(w, r, &req) {
		return
	}
	c, ok := s.sessions.TakeCeremony(req.CeremonyID, app, "login")
	if !ok {
		writeError(w, http.StatusBadRequest, "unknown or expired ceremony")
		return
	}
	user, err := s.store.Get(app, c.Username)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "login rejected")
		return
	}
	resp := req.Credential.Response
	if resp.UserHandle != nil && !bytes.Equal(resp.UserHandle, user.ID) {
		writeError(w, http.StatusUnauthorized, "login rejected: user handle mismatch")
		return
	}
	var cred *webauthn.Credential
	for i := range user.Credentials {
		if bytes.Equal(user.Credentials[i].ID, req.Credential.RawID) {
			cred = &user.Credentials[i]
		}
	}
	if cred == nil {
		writeError(w, http.StatusUnauthorized, "login rejected: credential not registered for this user")
		return
	}
	res, err := webauthn.VerifyAssertion(rp, c.Challenge, cred, webauthn.AssertionInput{
		ClientDataJSON:    resp.ClientDataJSON,
		AuthenticatorData: resp.AuthenticatorData,
		Signature:         resp.Signature,
	})
	if err != nil {
		log.Printf("login %s/%s rejected: %v", app, c.Username, err)
		writeError(w, http.StatusUnauthorized, "login rejected: "+err.Error())
		return
	}
	if err := s.store.UpdateSignCount(app, user.Name, cred.ID, res.SignCount, res.BackedUp); err != nil {
		writeError(w, http.StatusInternalServerError, "could not update credential")
		return
	}
	log.Printf("login %s/%s ok signCount=%d", app, c.Username, res.SignCount)
	s.issue(w, app, user)
}

func (s *Server) issue(w http.ResponseWriter, app string, user *store.User) {
	token, sess := s.sessions.Create(app, user.Name, webauthn.B64(user.ID))
	writeJSON(w, http.StatusOK, map[string]any{"token": token, "session": sess})
}

func bearer(r *http.Request) string {
	h := r.Header.Get("Authorization")
	if t, ok := strings.CutPrefix(h, "Bearer "); ok {
		return t
	}
	return ""
}

func (s *Server) session(w http.ResponseWriter, r *http.Request) {
	sess, ok := s.sessions.Lookup(bearer(r))
	if !ok {
		writeError(w, http.StatusUnauthorized, "invalid or expired session")
		return
	}
	writeJSON(w, http.StatusOK, sess)
}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	s.sessions.Revoke(bearer(r))
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) passkeys(w http.ResponseWriter, r *http.Request, app string, _ webauthn.RelyingParty) {
	sess, ok := s.sessions.Lookup(bearer(r))
	if !ok || sess.App != app {
		writeError(w, http.StatusUnauthorized, "invalid or expired session")
		return
	}
	user, err := s.store.Get(app, sess.Username)
	if err != nil {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}
	out := make([]passkeyView, 0, len(user.Credentials))
	for _, c := range user.Credentials {
		out = append(out, passkeyView{
			ID:             webauthn.B64(c.ID),
			Alg:            c.Alg,
			SignCount:      c.SignCount,
			AttestationFmt: c.AttestationFmt,
			BackupEligible: c.BackupEligible,
			BackedUp:       c.BackedUp,
			Transports:     c.Transports,
		})
	}
	writeJSON(w, http.StatusOK, out)
}

func readJSON(w http.ResponseWriter, r *http.Request, v any) bool {
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(v); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return false
	}
	return true
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
