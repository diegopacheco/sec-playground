package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"
)

type Server struct {
	firewall       *Firewall
	instructionDir string
	receiverCalls  atomic.Int64
}

type runRequest struct {
	Document string `json:"document"`
	Approval string `json:"approval"`
}

func main() {
	port := env("PORT", "8091")
	canary := env("CANARY_SECRET", "FIREWALL_CANARY_7f3c91")
	approval := env("APPROVAL_TOKEN", "approve-once")
	localURL := "http://127.0.0.1:" + port
	firewall := NewFirewall(FirewallConfig{
		ReadRoot:      "fixtures/data",
		AllowedHosts:  []string{"127.0.0.1:" + port},
		Canary:        canary,
		ApprovalToken: approval,
		LocalURL:      localURL,
		AuditPath:     env("AUDIT_PATH", "audit.jsonl"),
	})
	server := &Server{firewall: firewall, instructionDir: "fixtures/instructions"}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", server.health)
	mux.HandleFunc("POST /run", server.run)
	mux.HandleFunc("GET /audit", server.audit)
	mux.HandleFunc("GET /receiver", server.receiver)
	mux.HandleFunc("POST /receiver", server.receiver)
	httpServer := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
	log.Printf("agent firewall listening on %s", localURL)
	log.Fatal(httpServer.ListenAndServe())
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) run(w http.ResponseWriter, r *http.Request) {
	var input runRequest
	if err := decodeJSON(r.Body, &input); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	path, err := safeJoin(s.instructionDir, input.Document)
	if err != nil {
		writeError(w, http.StatusForbidden, err)
		return
	}
	content, err := os.ReadFile(path)
	if err != nil {
		writeError(w, http.StatusNotFound, errors.New("instruction document not found"))
		return
	}
	candidates, err := ParseCandidates(string(content))
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	decisions := make([]Decision, 0, len(candidates))
	for _, candidate := range candidates {
		decisions = append(decisions, s.firewall.Execute(input.Document, candidate, input.Approval))
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"document":  input.Document,
		"untrusted": true,
		"decisions": decisions,
	})
}

func (s *Server) receiver(w http.ResponseWriter, r *http.Request) {
	count := s.receiverCalls.Add(1)
	body, _ := io.ReadAll(io.LimitReader(r.Body, 4096))
	writeJSON(w, http.StatusOK, map[string]any{
		"accepted": true,
		"method":   r.Method,
		"body":     string(body),
		"calls":    count,
	})
}

func (s *Server) audit(w http.ResponseWriter, _ *http.Request) {
	entries, err := s.firewall.AuditEntries()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, entries)
}

func safeJoin(root, requested string) (string, error) {
	if requested == "" || filepath.IsAbs(requested) {
		return "", errors.New("path outside capability root")
	}
	rootPath, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	joined := filepath.Join(rootPath, filepath.Clean(requested))
	relative, err := filepath.Rel(rootPath, joined)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "", errors.New("path outside capability root")
	}
	return joined, nil
}

func decodeJSON(reader io.Reader, value any) error {
	decoder := json.NewDecoder(io.LimitReader(reader, 1<<20))
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

func env(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
