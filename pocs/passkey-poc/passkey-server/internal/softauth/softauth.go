package softauth

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"math/big"

	"passkey-server/internal/cbor"
	"passkey-server/internal/cose"
)

type Authenticator struct {
	Alg       int64
	Signer    crypto.Signer
	CredID    []byte
	SignCount uint32
	Flags     byte
	Format    string
	Synced    bool
}

const DefaultFlags byte = 0x01 | 0x04

func New(alg int64) *Authenticator {
	a := &Authenticator{Alg: alg, CredID: random(32), Flags: DefaultFlags, Format: "none"}
	switch alg {
	case cose.AlgES256:
		a.Signer, _ = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case cose.AlgEdDSA:
		_, a.Signer, _ = ed25519.GenerateKey(rand.Reader)
	case cose.AlgRS256:
		a.Signer, _ = rsa.GenerateKey(rand.Reader, 2048)
	}
	return a
}

func random(n int) []byte {
	b := make([]byte, n)
	rand.Read(b)
	return b
}

func ClientData(typ string, challenge []byte, origin string) []byte {
	b, _ := json.Marshal(map[string]any{
		"type":        typ,
		"challenge":   base64.RawURLEncoding.EncodeToString(challenge),
		"origin":      origin,
		"crossOrigin": false,
	})
	return b
}

func (a *Authenticator) COSEKey() []byte {
	switch pub := a.Signer.Public().(type) {
	case *ecdsa.PublicKey:
		raw, _ := pub.Bytes()
		return cbor.MustEncode(cbor.Map{{int64(1), int64(2)}, {int64(3), cose.AlgES256}, {int64(-1), int64(1)}, {int64(-2), raw[1:33]}, {int64(-3), raw[33:]}})
	case ed25519.PublicKey:
		return cbor.MustEncode(cbor.Map{{int64(1), int64(1)}, {int64(3), cose.AlgEdDSA}, {int64(-1), int64(6)}, {int64(-2), []byte(pub)}})
	case *rsa.PublicKey:
		return cbor.MustEncode(cbor.Map{{int64(1), int64(3)}, {int64(3), cose.AlgRS256}, {int64(-1), pub.N.Bytes()}, {int64(-2), big.NewInt(int64(pub.E)).Bytes()}})
	}
	return nil
}

func (a *Authenticator) authData(rpID string, attested bool) []byte {
	h := sha256.Sum256([]byte(rpID))
	out := append([]byte{}, h[:]...)
	flags := a.Flags
	if attested {
		flags |= 0x40
	}
	out = append(out, flags)
	out = binary.BigEndian.AppendUint32(out, a.SignCount)
	if attested {
		out = append(out, make([]byte, 16)...)
		out = binary.BigEndian.AppendUint16(out, uint16(len(a.CredID)))
		out = append(out, a.CredID...)
		out = append(out, a.COSEKey()...)
	}
	return out
}

func (a *Authenticator) Sign(data []byte) []byte {
	if _, ok := a.Signer.(ed25519.PrivateKey); ok {
		sig, _ := a.Signer.Sign(rand.Reader, data, crypto.Hash(0))
		return sig
	}
	digest := sha256.Sum256(data)
	sig, _ := a.Signer.Sign(rand.Reader, digest[:], crypto.SHA256)
	return sig
}

func (a *Authenticator) Create(rpID, origin string, challenge []byte) (clientData, attestationObject []byte) {
	clientData = ClientData("webauthn.create", challenge, origin)
	auth := a.authData(rpID, true)
	stmt := cbor.Map{}
	if a.Format == "packed" {
		hash := sha256.Sum256(clientData)
		stmt = cbor.Map{{"alg", a.Alg}, {"sig", a.Sign(append(append([]byte{}, auth...), hash[:]...))}}
	}
	attestationObject = cbor.MustEncode(cbor.Map{{"fmt", a.Format}, {"attStmt", stmt}, {"authData", auth}})
	return clientData, attestationObject
}

func (a *Authenticator) Get(rpID, origin string, challenge []byte) (clientData, authData, signature []byte) {
	if !a.Synced {
		a.SignCount++
	}
	clientData = ClientData("webauthn.get", challenge, origin)
	authData = a.authData(rpID, false)
	hash := sha256.Sum256(clientData)
	signature = a.Sign(append(append([]byte{}, authData...), hash[:]...))
	return clientData, authData, signature
}
