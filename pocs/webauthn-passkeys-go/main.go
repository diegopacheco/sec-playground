package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

type App struct {
	store    *Store
	webauthn *webauthn.WebAuthn
}

type beginRequest struct {
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
}

type beginResponse struct {
	SessionID string `json:"sessionId"`
	Options   any    `json:"options"`
}

func main() {
	port := env("PORT", "8090")
	origin := env("WEBAUTHN_ORIGIN", "http://localhost:"+port)
	rpID := env("WEBAUTHN_RP_ID", "localhost")
	store, err := OpenStore(env("DB_PATH", "passkeys.db"))
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()
	wa, err := webauthn.New(&webauthn.Config{
		RPID:          rpID,
		RPDisplayName: "Passkey Foundry",
		RPOrigins:     []string{origin},
		AuthenticatorSelection: protocol.AuthenticatorSelection{
			ResidentKey:      protocol.ResidentKeyRequirementRequired,
			UserVerification: protocol.VerificationRequired,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	app := &App{store: store, webauthn: wa}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", app.health)
	mux.HandleFunc("GET /api/status", app.status)
	mux.HandleFunc("POST /api/register/begin", app.beginRegistration)
	mux.HandleFunc("POST /api/register/finish", app.finishRegistration)
	mux.HandleFunc("POST /api/login/begin", app.beginLogin)
	mux.HandleFunc("POST /api/login/finish", app.finishLogin)
	mux.Handle("/", http.FileServer(http.Dir("web")))
	server := &http.Server{
		Addr:              ":" + port,
		Handler:           securityHeaders(mux),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
	log.Printf("passkey service listening on %s", origin)
	log.Fatal(server.ListenAndServe())
}

func (a *App) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (a *App) status(w http.ResponseWriter, _ *http.Request) {
	users, credentials, err := a.store.Counts()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]int{"users": users, "credentials": credentials})
}

func (a *App) beginRegistration(w http.ResponseWriter, r *http.Request) {
	var input beginRequest
	if err := readJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	input.Username = strings.TrimSpace(strings.ToLower(input.Username))
	input.DisplayName = strings.TrimSpace(input.DisplayName)
	if input.Username == "" || input.DisplayName == "" {
		writeError(w, http.StatusBadRequest, errors.New("username and display name are required"))
		return
	}
	user, err := a.store.GetOrCreateUser(input.Username, input.DisplayName)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if len(user.Credentials) != 0 {
		writeError(w, http.StatusConflict, errors.New("account already has a passkey"))
		return
	}
	options, session, err := a.webauthn.BeginRegistration(user,
		webauthn.WithResidentKeyRequirement(protocol.ResidentKeyRequirementRequired),
		webauthn.WithAuthenticatorSelection(protocol.AuthenticatorSelection{
			ResidentKey:      protocol.ResidentKeyRequirementRequired,
			UserVerification: protocol.VerificationRequired,
		}),
	)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	sessionID, err := a.store.SaveSession("registration", user.ID, session)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, beginResponse{SessionID: sessionID, Options: options})
}

func (a *App) finishRegistration(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session")
	session, userID, err := a.store.ConsumeSession(sessionID, "registration")
	if err != nil {
		writeError(w, http.StatusUnauthorized, err)
		return
	}
	user, err := a.store.GetUserByID(userID)
	if err != nil {
		writeError(w, http.StatusUnauthorized, err)
		return
	}
	credential, err := a.webauthn.FinishRegistration(user, *session, r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, err)
		return
	}
	if err := a.store.SaveCredential(user.ID, credential); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"status": "passkey registered"})
}

func (a *App) beginLogin(w http.ResponseWriter, r *http.Request) {
	var input beginRequest
	if err := readJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	user, err := a.store.GetUser(strings.TrimSpace(strings.ToLower(input.Username)))
	if err != nil {
		writeError(w, http.StatusUnauthorized, errors.New("account or passkey not found"))
		return
	}
	options, session, err := a.webauthn.BeginLogin(user, webauthn.WithUserVerification(protocol.VerificationRequired))
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	sessionID, err := a.store.SaveSession("login", user.ID, session)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, beginResponse{SessionID: sessionID, Options: options})
}

func (a *App) finishLogin(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session")
	session, userID, err := a.store.ConsumeSession(sessionID, "login")
	if err != nil {
		writeError(w, http.StatusUnauthorized, err)
		return
	}
	user, err := a.store.GetUserByID(userID)
	if err != nil {
		writeError(w, http.StatusUnauthorized, err)
		return
	}
	credential, err := a.webauthn.FinishLogin(user, *session, r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, err)
		return
	}
	if err := a.store.SaveCredential(user.ID, credential); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "authenticated", "username": user.Name})
}

func randomID(size int) ([]byte, error) {
	value := make([]byte, size)
	_, err := rand.Read(value)
	return value, err
}

func randomToken() (string, error) {
	value, err := randomID(32)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(value), nil
}

func env(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}

func readJSON(r *http.Request, value any) error {
	decoder := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	return decoder.Decode(value)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self'; connect-src 'self'; frame-ancestors 'none'")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		next.ServeHTTP(w, r)
	})
}
