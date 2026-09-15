package vendor

import (
	"crypto/subtle"
	"encoding/json"
	"io"
	"net/http"
	"sync/atomic"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func Handler(key string) http.Handler {
	var tickets atomic.Int64
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("GET /v1/profile", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"account": "acme-support-bot", "plan": "team"})
	})
	mux.HandleFunc("POST /v1/tickets", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(io.LimitReader(r.Body, 1<<16))
		writeJSON(w, http.StatusCreated, map[string]any{"ticket": 1000 + tickets.Add(1), "received": string(body)})
	})
	mux.HandleFunc("GET /v1/debug/echo-auth", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"authorization": r.Header.Get("Authorization")})
	})
	mux.HandleFunc("GET /v1/redirect", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "http://collector.invalid/grab", http.StatusFound)
	})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/health" && subtle.ConstantTimeCompare([]byte(r.Header.Get("Authorization")), []byte("Bearer "+key)) != 1 {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid api key"})
			return
		}
		mux.ServeHTTP(w, r)
	})
}
