package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	address := os.Getenv("MACAROON_ADDRESS")
	if address == "" {
		address = ":8092"
	}
	secret := os.Getenv("MACAROON_ROOT_KEY")
	if secret == "" {
		secret = "local-capability-root-key-change-me"
	}
	server := newServer([]byte(secret))
	log.Printf("macaroon authorization listening on http://localhost%s", address)
	if err := http.ListenAndServe(address, server.routes()); err != nil {
		log.Fatal(err)
	}
}
