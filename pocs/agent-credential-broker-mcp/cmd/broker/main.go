package main

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"agent-credential-broker-mcp/internal/audit"
	"agent-credential-broker-mcp/internal/broker"
	"agent-credential-broker-mcp/internal/httpapi"
	"agent-credential-broker-mcp/internal/policy"
	"agent-credential-broker-mcp/internal/provider"
)

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	addr := env("BROKER_ADDR", "127.0.0.1:8731")
	dataDir := env("DATA_DIR", ".run/data")
	if err := os.MkdirAll(dataDir, 0o700); err != nil {
		log.Fatal(err)
	}
	p, err := policy.Load(env("POLICY_FILE", "policy.json"))
	if err != nil {
		log.Fatalf("policy: %v", err)
	}
	auditLog, err := audit.Open(filepath.Join(dataDir, "audit.jsonl"), time.Now)
	if err != nil {
		log.Fatalf("audit: %v", err)
	}
	source := provider.Router{
		Vault:    provider.Vault{Addr: os.Getenv("VAULT_ADDR"), Token: os.Getenv("VAULT_TOKEN")},
		Keychain: provider.Keychain{},
	}
	raw := make([]byte, 24)
	rand.Read(raw)
	adminToken := hex.EncodeToString(raw)
	tokenFile := env("ADMIN_TOKEN_FILE", filepath.Join(dataDir, "admin-token"))
	if err := os.WriteFile(tokenFile, []byte(adminToken), 0o600); err != nil {
		log.Fatalf("admin token: %v", err)
	}
	srv := &httpapi.Server{
		Broker:      broker.New(p, source, auditLog, time.Now),
		AdminToken:  adminToken,
		UpstreamURL: env("UPSTREAM_URL", "http://localhost:8732"),
	}
	log.Printf("broker listening on http://%s with %d credential(s)", addr, len(p.Credentials))
	server := &http.Server{Addr: addr, Handler: srv.Handler(), ReadHeaderTimeout: 5 * time.Second}
	log.Fatal(server.ListenAndServe())
}
