package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
)

type server struct {
	macaroons *macaroonService
}

type mintRequest struct {
	Resource         string `json:"resource"`
	Operation        string `json:"operation"`
	Location         string `json:"location"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

type tokenResponse struct {
	Token     string   `json:"token"`
	ID        string   `json:"id"`
	Caveats   []Caveat `json:"caveats"`
	Signature string   `json:"signature"`
}

type attenuationRequest struct {
	Token  string `json:"token"`
	Caveat Caveat `json:"caveat"`
}

type verificationRequest struct {
	Token string `json:"token"`
	AccessRequest
}

func newServer(rootKey []byte) *server {
	return &server{macaroons: newMacaroonService(rootKey)}
}

func (s *server) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", s.health)
	mux.HandleFunc("POST /api/macaroon", s.mint)
	mux.HandleFunc("POST /api/attenuate", s.attenuate)
	mux.HandleFunc("POST /api/verify", s.verify)
	mux.Handle("/", http.FileServer(http.Dir("web")))
	return secure(mux)
}

func (s *server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *server) mint(w http.ResponseWriter, r *http.Request) {
	var request mintRequest
	if err := readJSON(w, r, &request); err != nil {
		writeError(w, err)
		return
	}
	if request.ExpiresInSeconds < 30 || request.ExpiresInSeconds > 3600 {
		writeError(w, errors.New("expiration must be between 30 and 3600 seconds"))
		return
	}
	for _, value := range []string{request.Resource, request.Operation, request.Location} {
		if strings.TrimSpace(value) == "" {
			writeError(w, errors.New("resource, operation, and location are required"))
			return
		}
	}
	token, macaroon, err := s.macaroons.mint(request.Resource, request.Operation, request.Location, time.Duration(request.ExpiresInSeconds)*time.Second)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, tokenResponse{Token: token, ID: macaroon.ID, Caveats: macaroon.Caveats, Signature: macaroon.Signature})
}

func (s *server) attenuate(w http.ResponseWriter, r *http.Request) {
	var request attenuationRequest
	if err := readJSON(w, r, &request); err != nil {
		writeError(w, err)
		return
	}
	if request.Token == "" {
		writeError(w, errors.New("token is required"))
		return
	}
	token, macaroon, err := s.macaroons.attenuate(request.Token, request.Caveat)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, tokenResponse{Token: token, ID: macaroon.ID, Caveats: macaroon.Caveats, Signature: macaroon.Signature})
}

func (s *server) verify(w http.ResponseWriter, r *http.Request) {
	var request verificationRequest
	if err := readJSON(w, r, &request); err != nil {
		writeError(w, err)
		return
	}
	if request.Token == "" {
		writeError(w, errors.New("token is required"))
		return
	}
	writeJSON(w, http.StatusOK, s.macaroons.verify(request.Token, request.AccessRequest))
}

func secure(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self'; script-src 'self'; img-src 'self' data:; connect-src 'self'; frame-ancestors 'none'")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "no-referrer")
		next.ServeHTTP(w, r)
	})
}

func readJSON(w http.ResponseWriter, r *http.Request, destination any) error {
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		return errors.New("request body is not valid")
	}
	return nil
}

func writeError(w http.ResponseWriter, err error) {
	writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(value)
}
