package httpapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"agent-intent-bound-tokens/internal/intent"
	"agent-intent-bound-tokens/internal/token"
	"agent-intent-bound-tokens/internal/verification"
)

type Server struct {
	issuer       string
	tokens       *token.Service
	verification *verification.Service
}

type issueRequest struct {
	Subject         string `json:"subject"`
	Audience        string `json:"audience"`
	Action          string `json:"action"`
	Resource        string `json:"resource"`
	MaxAmountCents  int64  `json:"max_amount_cents"`
	ValidForSeconds int64  `json:"valid_for_seconds"`
}

type issueResponse struct {
	Token  string        `json:"token"`
	Claims intent.Claims `json:"claims"`
}

type endpoint struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}

func New(issuer string, tokens *token.Service, verification *verification.Service) *Server {
	return &Server{issuer: issuer, tokens: tokens, verification: verification}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", s.root)
	mux.HandleFunc("GET /health", s.health)
	mux.HandleFunc("POST /api/tokens", s.issue)
	mux.HandleFunc("POST /api/verify", s.verify)
	mux.HandleFunc("GET /api/audit", s.audit)
	return mux
}

func (s *Server) root(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string][]endpoint{
		"endpoints": {
			{Method: http.MethodGet, Path: "/"},
			{Method: http.MethodGet, Path: "/health"},
			{Method: http.MethodPost, Path: "/api/tokens"},
			{Method: http.MethodPost, Path: "/api/verify"},
			{Method: http.MethodGet, Path: "/api/audit"},
		},
	})
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}

func (s *Server) issue(w http.ResponseWriter, r *http.Request) {
	var input issueRequest
	if err := decode(r, &input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if input.ValidForSeconds < 1 || input.ValidForSeconds > 3600 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "valid_for_seconds must be between 1 and 3600"})
		return
	}
	now := time.Now().UTC()
	value, claims, err := s.tokens.Issue(intent.Claims{
		Issuer: s.issuer, Subject: input.Subject, Audience: input.Audience,
		Action: input.Action, Resource: input.Resource, MaxAmountCents: input.MaxAmountCents,
		NotBefore: now.Unix(), ExpiresAt: now.Add(time.Duration(input.ValidForSeconds) * time.Second).Unix(),
	})
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, issueResponse{Token: value, Claims: claims})
}

func (s *Server) verify(w http.ResponseWriter, r *http.Request) {
	var input intent.Request
	if err := decode(r, &input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	decision := s.verification.Verify(input)
	status := http.StatusForbidden
	if decision.Allowed {
		status = http.StatusOK
	}
	writeJSON(w, status, decision)
}

func (s *Server) audit(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, s.verification.Audit())
}

func decode(r *http.Request, target any) error {
	decoder := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return errors.New("request body is invalid")
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
