package webauthn

import (
	"crypto/rand"
	"encoding/base64"

	"passkey-server/internal/cose"
)

const (
	ChallengeSize = 32
	TimeoutMillis = 120000
)

type RelyingParty struct {
	ID      string
	Name    string
	Origins []string
}

type Descriptor struct {
	Type       string   `json:"type"`
	ID         string   `json:"id"`
	Transports []string `json:"transports,omitempty"`
}

type Entity struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName,omitempty"`
}

type PubKeyParam struct {
	Type string `json:"type"`
	Alg  int64  `json:"alg"`
}

type AuthenticatorSelection struct {
	ResidentKey      string `json:"residentKey"`
	UserVerification string `json:"userVerification"`
}

type CreationOptions struct {
	RP                     Entity                 `json:"rp"`
	User                   Entity                 `json:"user"`
	Challenge              string                 `json:"challenge"`
	PubKeyCredParams       []PubKeyParam          `json:"pubKeyCredParams"`
	Timeout                int                    `json:"timeout"`
	ExcludeCredentials     []Descriptor           `json:"excludeCredentials"`
	AuthenticatorSelection AuthenticatorSelection `json:"authenticatorSelection"`
	Attestation            string                 `json:"attestation"`
}

type RequestOptions struct {
	Challenge        string       `json:"challenge"`
	RPID             string       `json:"rpId"`
	Timeout          int          `json:"timeout"`
	AllowCredentials []Descriptor `json:"allowCredentials"`
	UserVerification string       `json:"userVerification"`
}

func NewChallenge() []byte {
	b := make([]byte, ChallengeSize)
	rand.Read(b)
	return b
}

func B64(b []byte) string {
	return base64.RawURLEncoding.EncodeToString(b)
}

func NewCreationOptions(rp RelyingParty, userID []byte, username string, challenge []byte, exclude []Descriptor) CreationOptions {
	return CreationOptions{
		RP:        Entity{ID: rp.ID, Name: rp.Name},
		User:      Entity{ID: B64(userID), Name: username, DisplayName: username},
		Challenge: B64(challenge),
		PubKeyCredParams: []PubKeyParam{
			{Type: "public-key", Alg: cose.AlgES256},
			{Type: "public-key", Alg: cose.AlgEdDSA},
			{Type: "public-key", Alg: cose.AlgRS256},
		},
		Timeout:            TimeoutMillis,
		ExcludeCredentials: nonNil(exclude),
		AuthenticatorSelection: AuthenticatorSelection{
			ResidentKey:      "preferred",
			UserVerification: "required",
		},
		Attestation: "none",
	}
}

func NewRequestOptions(rp RelyingParty, challenge []byte, allow []Descriptor) RequestOptions {
	return RequestOptions{
		Challenge:        B64(challenge),
		RPID:             rp.ID,
		Timeout:          TimeoutMillis,
		AllowCredentials: nonNil(allow),
		UserVerification: "required",
	}
}

func nonNil(d []Descriptor) []Descriptor {
	if d == nil {
		return []Descriptor{}
	}
	return d
}
