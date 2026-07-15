package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"token-blast-radius/internal/graph"
	"token-blast-radius/internal/httpapi"
	"token-blast-radius/internal/scenario"
	"token-blast-radius/internal/token"
)

func main() {
	tokens, err := token.New(os.Getenv("BLAST_RADIUS_SECRET"))
	if err != nil {
		log.Fatal(err)
	}
	scenarios, err := scenario.Build(tokens)
	if err != nil {
		log.Fatal(err)
	}
	port := env("PORT", "8082")
	frontendPath := env("FRONTEND_PATH", "frontend/dist")
	api := httpapi.New(graph.New(tokens, scenarios), frontendPath)
	server := &http.Server{
		Addr: ":" + port, Handler: api.Handler(), ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout: 10 * time.Second, WriteTimeout: 15 * time.Second, IdleTimeout: 30 * time.Second,
	}
	log.Printf("blast radius service listening on http://localhost:%s", port)
	log.Fatal(server.ListenAndServe())
}

func env(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
