package httpapi

import (
	"crypto/subtle"
	"embed"
	"encoding/json"
	"errors"
	"io/fs"
	"net/http"

	"agent-credential-broker-mcp/internal/audit"
	"agent-credential-broker-mcp/internal/broker"
	"agent-credential-broker-mcp/internal/lab"
)

//go:embed web
var webFS embed.FS

type Server struct {
	Broker      *broker.Broker
	AdminToken  string
	UpstreamURL string
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	static, _ := fs.Sub(webFS, "web")
	mux.Handle("GET /", http.FileServerFS(static))
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("GET /agent/credentials", s.credentials)
	mux.HandleFunc("POST /agent/intents", s.requestAccess)
	mux.HandleFunc("GET /agent/intents/{id}", s.checkAccess)
	mux.HandleFunc("POST /agent/call", s.call)
	mux.HandleFunc("GET /api/credentials", s.admin(s.credentials))
	mux.HandleFunc("GET /api/intents", s.admin(s.intents))
	mux.HandleFunc("POST /api/intents/{id}/approve", s.admin(s.decide(true)))
	mux.HandleFunc("POST /api/intents/{id}/deny", s.admin(s.decide(false)))
	mux.HandleFunc("GET /api/audit", s.admin(s.audit))
	mux.HandleFunc("GET /api/lab", s.admin(s.scenarios))
	mux.HandleFunc("POST /api/lab/{id}", s.admin(s.runScenario))
	return secureHeaders(mux)
}

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("Content-Security-Policy", "default-src 'self'; frame-ancestors 'none'")
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("Referrer-Policy", "no-referrer")
		h.Set("Cache-Control", "no-store")
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
		next.ServeHTTP(w, r)
	})
}

func (s *Server) admin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		given := r.Header.Get("X-Admin-Token")
		if s.AdminToken == "" || subtle.ConstantTimeCompare([]byte(given), []byte(s.AdminToken)) != 1 {
			writeError(w, http.StatusUnauthorized, errors.New("admin token required"))
			return
		}
		next(w, r)
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

func (s *Server) credentials(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.Broker.Credentials())
}

func (s *Server) requestAccess(w http.ResponseWriter, r *http.Request) {
	var req broker.AccessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	in, err := s.Broker.RequestAccess(req)
	if err != nil {
		writeError(w, http.StatusForbidden, err)
		return
	}
	writeJSON(w, http.StatusCreated, in)
}

func (s *Server) checkAccess(w http.ResponseWriter, r *http.Request) {
	in, token, err := s.Broker.Check(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"intent": in, "token": token})
}

func (s *Server) call(w http.ResponseWriter, r *http.Request) {
	var req broker.CallRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	res, err := s.Broker.Call(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusForbidden, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (s *Server) intents(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.Broker.Intents())
}

func (s *Server) decide(approve bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var in broker.Intent
		var err error
		if approve {
			in, err = s.Broker.Approve(r.PathValue("id"))
		} else {
			in, err = s.Broker.Deny(r.PathValue("id"))
		}
		switch {
		case errors.Is(err, broker.ErrNotFound):
			writeError(w, http.StatusNotFound, err)
		case err != nil:
			writeError(w, http.StatusConflict, err)
		default:
			writeJSON(w, http.StatusOK, in)
		}
	}
}

func (s *Server) audit(w http.ResponseWriter, r *http.Request) {
	entries := s.Broker.Audit().Entries()
	ok, broken := audit.Verify(entries)
	writeJSON(w, http.StatusOK, map[string]any{"verified": ok, "broken_at": broken, "entries": entries})
}

func (s *Server) scenarios(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, lab.Scenarios)
}

func (s *Server) runScenario(w http.ResponseWriter, r *http.Request) {
	steps, err := lab.Run(r.Context(), s.Broker, s.UpstreamURL, r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	writeJSON(w, http.StatusOK, steps)
}
