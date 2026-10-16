package webauthn

import (
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"slices"
)

const (
	TypeCreate = "webauthn.create"
	TypeGet    = "webauthn.get"
)

type ClientData struct {
	Type        string `json:"type"`
	Challenge   string `json:"challenge"`
	Origin      string `json:"origin"`
	CrossOrigin bool   `json:"crossOrigin"`
}

func VerifyClientData(raw []byte, wantType string, challenge []byte, origins []string) (*ClientData, error) {
	var cd ClientData
	if err := json.Unmarshal(raw, &cd); err != nil {
		return nil, fmt.Errorf("clientDataJSON: %w", err)
	}
	if cd.Type != wantType {
		return nil, fmt.Errorf("clientDataJSON: type %q, want %q", cd.Type, wantType)
	}
	got, err := base64.RawURLEncoding.DecodeString(cd.Challenge)
	if err != nil {
		return nil, fmt.Errorf("clientDataJSON: challenge is not base64url")
	}
	if subtle.ConstantTimeCompare(got, challenge) != 1 {
		return nil, fmt.Errorf("clientDataJSON: challenge mismatch")
	}
	if !slices.Contains(origins, cd.Origin) {
		return nil, fmt.Errorf("clientDataJSON: origin %q not allowed", cd.Origin)
	}
	if cd.CrossOrigin {
		return nil, fmt.Errorf("clientDataJSON: cross origin ceremonies are not allowed")
	}
	return &cd, nil
}
