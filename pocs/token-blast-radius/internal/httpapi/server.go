package httpapi

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"token-blast-radius/internal/domain"
	"token-blast-radius/internal/graph"
)

type Server struct {
	engine       *graph.Engine
	frontendPath string
}

func New(engine *graph.Engine, frontendPath string) *Server {
	return &Server{engine: engine, frontendPath: frontendPath}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.health)
	mux.HandleFunc("GET /api/scenarios", s.scenarios)
	mux.HandleFunc("POST /api/analyze", s.analyze)
	mux.Handle("/", s.frontend())
	return mux
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}

func (s *Server) scenarios(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, s.engine.Summaries())
}

func (s *Server) analyze(w http.ResponseWriter, r *http.Request) {
	var input domain.AnalysisRequest
	decoder := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "request body is invalid"})
		return
	}
	analysis, err := s.engine.Analyze(input.ScenarioID, input.Token)
	if err != nil {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, analysis)
}

func (s *Server) frontend() http.Handler {
	files := http.FileServer(http.Dir(s.frontendPath))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cleaned := filepath.Clean(strings.TrimPrefix(r.URL.Path, "/"))
		if cleaned == "." {
			cleaned = "index.html"
		}
		path := filepath.Join(s.frontendPath, cleaned)
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			files.ServeHTTP(w, r)
			return
		}
		http.ServeFile(w, r, filepath.Join(s.frontendPath, "index.html"))
	})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
