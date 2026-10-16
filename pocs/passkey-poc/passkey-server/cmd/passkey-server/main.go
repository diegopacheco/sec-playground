package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"passkey-server/internal/api"
	"passkey-server/internal/session"
	"passkey-server/internal/store"
	"passkey-server/internal/webauthn"
)

func parseApps(spec, rpID string) (map[string]webauthn.RelyingParty, error) {
	apps := map[string]webauthn.RelyingParty{}
	for entry := range strings.SplitSeq(spec, ";") {
		parts := strings.SplitN(entry, "=", 3)
		if len(parts) != 3 {
			return nil, fmt.Errorf("app spec %q must be key=Name=origin", entry)
		}
		apps[parts[0]] = webauthn.RelyingParty{ID: rpID, Name: parts[1], Origins: []string{parts[2]}}
	}
	return apps, nil
}

func main() {
	addr := flag.String("addr", "127.0.0.1:8090", "listen address")
	rpID := flag.String("rp-id", "localhost", "relying party id")
	appSpec := flag.String("apps", "bookstore=Bookstore=http://localhost:5173;vault=PasskeyVault=http://localhost:5174", "key=Name=origin entries separated by ;")
	data := flag.String("data", "data/passkeys.json", "credential store file")
	flag.Parse()

	apps, err := parseApps(*appSpec, *rpID)
	if err != nil {
		log.Fatal(err)
	}
	st, err := store.Open(*data)
	if err != nil {
		log.Fatal(err)
	}
	srv := api.New(apps, st, session.NewManager(2*time.Minute, time.Hour))
	for k, rp := range apps {
		log.Printf("app %s rpId=%s origins=%v", k, rp.ID, rp.Origins)
	}
	log.Printf("passkey-server listening on http://%s", *addr)
	server := &http.Server{Addr: *addr, Handler: srv.Handler(), ReadHeaderTimeout: 5 * time.Second}
	log.Fatal(server.ListenAndServe())
}
