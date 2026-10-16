package main

import (
	"log"
	"net/http"
	"os"

	"mcp-output-taint-tracker/internal/api"
	"mcp-output-taint-tracker/internal/policy"
	"mcp-output-taint-tracker/internal/session"
)

func main() {
	port := os.Getenv("TAINT_PORT")
	if port == "" {
		port = "8751"
	}
	addr := "127.0.0.1:" + port
	handler := api.New(session.NewStore(policy.Default()))
	log.Printf("taint tracker listening on http://%s", addr)
	log.Fatal(http.ListenAndServe(addr, handler))
}
