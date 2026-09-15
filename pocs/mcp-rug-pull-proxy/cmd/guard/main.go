package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"mcp-rug-pull-proxy/internal/guard"
	"mcp-rug-pull-proxy/internal/store"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:8741", "listen address")
	runDir := flag.String("run-dir", ".run", "directory for tokens and data")
	flag.Parse()

	st, err := store.Open(filepath.Join(*runDir, "data"))
	if err != nil {
		log.Fatal(err)
	}
	tokens, err := guard.LoadOrCreateTokens(*runDir)
	if err != nil {
		log.Fatal(err)
	}
	srv := &http.Server{
		Addr:              *addr,
		Handler:           guard.New(st, tokens).Handler(),
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Printf("guard listening on http://%s", *addr)
	log.Fatal(srv.ListenAndServe())
}
