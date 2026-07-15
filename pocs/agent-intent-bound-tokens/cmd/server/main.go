package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"agent-intent-bound-tokens/internal/httpapi"
	"agent-intent-bound-tokens/internal/token"
	"agent-intent-bound-tokens/internal/verification"
)

func main() {
	secret := os.Getenv("INTENT_TOKEN_SECRET")
	tokens, err := token.New(secret)
	if err != nil {
		log.Fatal(err)
	}
	port := env("PORT", "8081")
	issuer := env("INTENT_TOKEN_ISSUER", "intent.local")
	verifier := verification.New(tokens)
	api := httpapi.New(issuer, tokens, verifier)
	server := &http.Server{
		Addr:              ":" + port,
		Handler:           api.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
	log.Printf("intent token service listening on http://localhost:%s", port)
	log.Fatal(server.ListenAndServe())
}

func env(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
