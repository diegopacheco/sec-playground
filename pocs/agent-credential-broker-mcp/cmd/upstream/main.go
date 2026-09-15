package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"agent-credential-broker-mcp/internal/vendor"
)

func main() {
	key := os.Getenv("UPSTREAM_API_KEY")
	if key == "" {
		log.Fatal("UPSTREAM_API_KEY is required")
	}
	addr := os.Getenv("UPSTREAM_ADDR")
	if addr == "" {
		addr = "127.0.0.1:8732"
	}
	log.Printf("vendor api listening on http://%s", addr)
	server := &http.Server{Addr: addr, Handler: vendor.Handler(key), ReadHeaderTimeout: 5 * time.Second}
	log.Fatal(server.ListenAndServe())
}
