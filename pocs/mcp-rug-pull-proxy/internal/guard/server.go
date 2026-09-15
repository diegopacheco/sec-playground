package guard

import (
	"crypto/subtle"
	"embed"
	"encoding/json"
	"errors"
	"io/fs"
	"net"
	"net/http"
	"strings"

	"mcp-rug-pull-proxy/internal/mcp"
	"mcp-rug-pull-proxy/internal/scan"
	"mcp-rug-pull-proxy/internal/store"
)

//go:embed web
var web embed.FS

type Tokens struct {
	Admin string
	Proxy string
}

type Server struct {
	store  *store.Store
	tokens Tokens
}

type checkRequest struct {
	Server string     `json:"server"`
	Tools  []mcp.Tool `json:"tools"`
}

type toolRequest struct {
	Server string `json:"server"`
	Tool   string `json:"tool"`
	Hash   string `json:"hash"`
}

type deniedRequest struct {
	Server string `json:"server"`
	Tool   string `json:"tool"`
	Reason string `json:"reason"`
}

func New(s *store.Store, tokens Tokens) *Server {
	return &Server{store: s, tokens: tokens}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	static, _ := fs.Sub(web, "web")
	mux.Handle("GET /", http.FileServerFS(static))
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("POST /api/check", s.requireToken(s.tokens.Proxy, "X-Proxy-Token", s.check))
	mux.HandleFunc("POST /api/denied", s.requireToken(s.tokens.Proxy, "X-Proxy-Token", s.denied))
	mux.HandleFunc("GET /api/tools", s.tools)
	mux.HandleFunc("GET /api/alerts", s.alerts)
	mux.HandleFunc("GET /api/events", s.events)
	mux.HandleFunc("POST /api/scan", s.scan)
	mux.HandleFunc("POST /api/tools/approve", s.requireToken(s.tokens.Admin, "X-Admin-Token", s.approve))
	mux.HandleFunc("POST /api/tools/reject", s.requireToken(s.tokens.Admin, "X-Admin-Token", s.reject))
	mux.HandleFunc("POST /api/servers/forget", s.requireToken(s.tokens.Admin, "X-Admin-Token", s.forget))
	return localOnly(securityHeaders(mux))
}

func (s *Server) check(w http.ResponseWriter, r *http.Request) {
	var req checkRequest
	if !decode(w, r, &req) {
		return
	}
	if req.Server == "" {
		writeError(w, http.StatusBadRequest, "server is required")
		return
	}
	decisions, err := s.store.Check(req.Server, req.Tools)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"decisions": decisions})
}

func (s *Server) denied(w http.ResponseWriter, r *http.Request) {
	var req deniedRequest
	if !decode(w, r, &req) {
		return
	}
	_ = s.store.RecordDenied(req.Server, req.Tool, req.Reason)
	writeJSON(w, http.StatusOK, map[string]string{"status": "recorded"})
}

func (s *Server) tools(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"pins": s.store.Pins()})
}

func (s *Server) alerts(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"alerts": s.store.Alerts()})
}

func (s *Server) events(w http.ResponseWriter, r *http.Request) {
	events, chain := s.store.Events()
	if events == nil {
		events = []store.Event{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"events": events, "chain": chain})
}

func (s *Server) scan(w http.ResponseWriter, r *http.Request) {
	var tool mcp.Tool
	if !decode(w, r, &tool) {
		return
	}
	writeJSON(w, http.StatusOK, scan.Scan(tool, s.store.ToolNamesExcept("", tool.Name)))
}

func (s *Server) approve(w http.ResponseWriter, r *http.Request) {
	var req toolRequest
	if !decode(w, r, &req) {
		return
	}
	s.mutation(w, s.store.Approve(req.Server, req.Tool, req.Hash))
}

func (s *Server) reject(w http.ResponseWriter, r *http.Request) {
	var req toolRequest
	if !decode(w, r, &req) {
		return
	}
	s.mutation(w, s.store.Reject(req.Server, req.Tool))
}

func (s *Server) forget(w http.ResponseWriter, r *http.Request) {
	var req toolRequest
	if !decode(w, r, &req) {
		return
	}
	s.mutation(w, s.store.Forget(req.Server))
}

func (s *Server) mutation(w http.ResponseWriter, err error) {
	switch {
	case err == nil:
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	case errors.Is(err, store.ErrNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, store.ErrHashMismatch):
		writeError(w, http.StatusConflict, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, err.Error())
	}
}

func (s *Server) requireToken(expected, header string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		got := r.Header.Get(header)
		if expected == "" || subtle.ConstantTimeCompare([]byte(got), []byte(expected)) != 1 {
			writeError(w, http.StatusUnauthorized, "missing or invalid "+header)
			return
		}
		next(w, r)
	}
}

func localOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host := r.Host
		if h, _, err := net.SplitHostPort(host); err == nil {
			host = h
		}
		host = strings.Trim(host, "[]")
		if host != "localhost" && host != "127.0.0.1" && host != "::1" {
			writeError(w, http.StatusForbidden, "guard only answers on localhost")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self'; script-src 'self'; img-src 'self' data:")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "no-referrer")
		next.ServeHTTP(w, r)
	})
}

func decode(w http.ResponseWriter, r *http.Request, v any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json: "+err.Error())
		return false
	}
	return true
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
