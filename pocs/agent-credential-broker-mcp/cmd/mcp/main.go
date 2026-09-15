package main

import (
	"context"
	"log"
	"os"

	"agent-credential-broker-mcp/internal/mcp"
)

func main() {
	base := os.Getenv("BROKER_URL")
	if base == "" {
		base = "http://127.0.0.1:8731"
	}
	log.SetOutput(os.Stderr)
	srv := &mcp.Server{Backend: mcp.Backend{BaseURL: base}, Name: "agent-credential-broker", Version: "1.0.0"}
	if err := srv.Serve(context.Background(), os.Stdin, os.Stdout); err != nil {
		log.Fatal(err)
	}
}
