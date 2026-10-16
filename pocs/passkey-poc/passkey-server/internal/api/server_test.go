package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"passkey-server/internal/cose"
	"passkey-server/internal/session"
	"passkey-server/internal/softauth"
	"passkey-server/internal/store"
	"passkey-server/internal/webauthn"
)

const (
	bookOrigin  = "http://localhost:5173"
	vaultOrigin = "http://localhost:5174"
)

func newServer(t *testing.T) *httptest.Server {
	t.Helper()
	st, _ := store.Open("")
	apps := map[string]webauthn.RelyingParty{
		"bookstore": {ID: "localhost", Name: "Bookstore", Origins: []string{bookOrigin}},
		"vault":     {ID: "localhost", Name: "Vault", Origins: []string{vaultOrigin}},
	}
	srv := httptest.NewServer(New(apps, st, session.NewManager(time.Minute, time.Hour)).Handler())
	t.Cleanup(srv.Close)
	return srv
}

func b64(b []byte) string {
	return base64.RawURLEncoding.EncodeToString(b)
}

func call(t *testing.T, srv *httptest.Server, method, path, token string, body any) (int, map[string]any) {
	t.Helper()
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, srv.URL+path, bytes.NewReader(data))
	req.Host = "localhost"
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	out := map[string]any{}
	json.NewDecoder(resp.Body).Decode(&out)
	return resp.StatusCode, out
}

func challengeOf(t *testing.T, begin map[string]any) []byte {
	t.Helper()
	raw, err := base64.RawURLEncoding.DecodeString(begin["publicKey"].(map[string]any)["challenge"].(string))
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

func signup(t *testing.T, srv *httptest.Server, a *softauth.Authenticator, app, origin, user string) (int, map[string]any) {
	t.Helper()
	code, begin := call(t, srv, "POST", "/api/"+app+"/register/begin", "", map[string]string{"username": user})
	if code != 200 {
		return code, begin
	}
	cd, att := a.Create("localhost", origin, challengeOf(t, begin))
	return call(t, srv, "POST", "/api/"+app+"/register/finish", "", map[string]any{
		"ceremonyId": begin["ceremonyId"],
		"credential": map[string]any{"rawId": b64(a.CredID), "type": "public-key", "response": map[string]any{"clientDataJSON": b64(cd), "attestationObject": b64(att)}},
	})
}

func signin(t *testing.T, srv *httptest.Server, a *softauth.Authenticator, app, origin, user string) (int, map[string]any) {
	t.Helper()
	code, begin := call(t, srv, "POST", "/api/"+app+"/login/begin", "", map[string]string{"username": user})
	if code != 200 {
		return code, begin
	}
	cd, ad, sig := a.Get("localhost", origin, challengeOf(t, begin))
	return call(t, srv, "POST", "/api/"+app+"/login/finish", "", map[string]any{
		"ceremonyId": begin["ceremonyId"],
		"credential": map[string]any{"rawId": b64(a.CredID), "type": "public-key", "response": map[string]any{"clientDataJSON": b64(cd), "authenticatorData": b64(ad), "signature": b64(sig)}},
	})
}

func TestSignupLoginSessionAndLogout(t *testing.T) {
	srv := newServer(t)
	a := softauth.New(cose.AlgES256)
	if code, out := signup(t, srv, a, "bookstore", bookOrigin, "alice"); code != 200 {
		t.Fatalf("signup failed %d %v", code, out)
	}
	code, out := signin(t, srv, a, "bookstore", bookOrigin, "alice")
	if code != 200 {
		t.Fatalf("login failed %d %v", code, out)
	}
	token := out["token"].(string)
	code, sess := call(t, srv, "GET", "/api/session", token, nil)
	if code != 200 || sess["username"] != "alice" || sess["app"] != "bookstore" {
		t.Fatalf("session must identify the user and the app: %d %v", code, sess)
	}
	call(t, srv, "POST", "/api/logout", token, nil)
	if code, _ := call(t, srv, "GET", "/api/session", token, nil); code != 401 {
		t.Fatalf("logged out token must stop working, got %d", code)
	}
}

func TestUsernameCannotBeTakenOver(t *testing.T) {
	srv := newServer(t)
	signup(t, srv, softauth.New(cose.AlgES256), "bookstore", bookOrigin, "alice")
	if code, _ := signup(t, srv, softauth.New(cose.AlgES256), "bookstore", bookOrigin, "alice"); code != 409 {
		t.Fatalf("a second signup must not attach an attacker passkey to an existing account, got %d", code)
	}
}

func TestPasskeyFromOneAppDoesNotWorkOnTheOther(t *testing.T) {
	srv := newServer(t)
	a := softauth.New(cose.AlgES256)
	signup(t, srv, a, "bookstore", bookOrigin, "alice")
	signup(t, srv, softauth.New(cose.AlgES256), "vault", vaultOrigin, "alice")
	if code, _ := signin(t, srv, a, "vault", vaultOrigin, "alice"); code != 401 {
		t.Fatalf("a bookstore passkey must not unlock the vault account, got %d", code)
	}
	if code, _ := signin(t, srv, a, "bookstore", vaultOrigin, "alice"); code != 401 {
		t.Fatalf("a ceremony run on the vault origin must not be accepted by the bookstore, got %d", code)
	}
}

func TestCeremonyIsSingleUse(t *testing.T) {
	srv := newServer(t)
	a := softauth.New(cose.AlgES256)
	signup(t, srv, a, "bookstore", bookOrigin, "alice")
	_, begin := call(t, srv, "POST", "/api/bookstore/login/begin", "", map[string]string{"username": "alice"})
	cd, ad, sig := a.Get("localhost", bookOrigin, challengeOf(t, begin))
	body := map[string]any{
		"ceremonyId": begin["ceremonyId"],
		"credential": map[string]any{"rawId": b64(a.CredID), "type": "public-key", "response": map[string]any{"clientDataJSON": b64(cd), "authenticatorData": b64(ad), "signature": b64(sig)}},
	}
	if code, _ := call(t, srv, "POST", "/api/bookstore/login/finish", "", body); code != 200 {
		t.Fatalf("first use must succeed")
	}
	if code, _ := call(t, srv, "POST", "/api/bookstore/login/finish", "", body); code != 400 {
		t.Fatalf("a captured assertion must not be replayable, got %d", code)
	}
}

func TestForeignHostIsRefused(t *testing.T) {
	srv := newServer(t)
	req, _ := http.NewRequest("GET", srv.URL+"/api/health", nil)
	req.Host = "attacker.test"
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 403 {
		t.Fatalf("DNS rebinding host must be refused, got %d", resp.StatusCode)
	}
}

func TestPasskeyListIsScopedToTheSessionApp(t *testing.T) {
	srv := newServer(t)
	_, out := signup(t, srv, softauth.New(cose.AlgEdDSA), "bookstore", bookOrigin, "alice")
	token := out["token"].(string)
	req, _ := http.NewRequest("GET", srv.URL+"/api/bookstore/passkeys", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, _ := http.DefaultClient.Do(req)
	var list []map[string]any
	json.NewDecoder(resp.Body).Decode(&list)
	resp.Body.Close()
	if len(list) != 1 || list[0]["alg"] != float64(cose.AlgEdDSA) {
		t.Fatalf("expected the one Ed25519 passkey, got %v", list)
	}
	if code, _ := call(t, srv, "GET", "/api/vault/passkeys", token, nil); code != 401 {
		t.Fatalf("a bookstore session must not read vault data, got %d", code)
	}
}
