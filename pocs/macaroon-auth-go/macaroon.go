package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Caveat struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Macaroon struct {
	Version   int      `json:"version"`
	ID        string   `json:"id"`
	Caveats   []Caveat `json:"caveats"`
	Signature string   `json:"signature"`
}

type AccessRequest struct {
	Resource  string `json:"resource"`
	Operation string `json:"operation"`
	Location  string `json:"location"`
}

type CaveatCheck struct {
	Type     string `json:"type"`
	Required string `json:"required"`
	Actual   string `json:"actual"`
	Passed   bool   `json:"passed"`
}

type Verification struct {
	Allowed bool          `json:"allowed"`
	Reason  string        `json:"reason"`
	Checks  []CaveatCheck `json:"checks"`
}

type macaroonService struct {
	rootKey []byte
	now     func() time.Time
}

func newMacaroonService(rootKey []byte) *macaroonService {
	return &macaroonService{rootKey: rootKey, now: time.Now}
}

func (s *macaroonService) mint(resource, operation, location string, lifetime time.Duration) (string, Macaroon, error) {
	idBytes := make([]byte, 16)
	if _, err := rand.Read(idBytes); err != nil {
		return "", Macaroon{}, err
	}
	macaroon := Macaroon{
		Version: 1,
		ID:      hex.EncodeToString(idBytes),
		Caveats: []Caveat{
			{Type: "resource", Value: resource},
			{Type: "operation", Value: operation},
			{Type: "location", Value: location},
			{Type: "expires", Value: s.now().UTC().Add(lifetime).Format(time.RFC3339)},
		},
	}
	for _, caveat := range macaroon.Caveats {
		if err := validateCaveat(caveat); err != nil {
			return "", Macaroon{}, err
		}
	}
	signature := sign(s.rootKey, macaroon.ID)
	for _, caveat := range macaroon.Caveats {
		signature = sign(signature, caveat.Type+"="+caveat.Value)
	}
	macaroon.Signature = hex.EncodeToString(signature)
	token, err := encodeMacaroon(macaroon)
	return token, macaroon, err
}

func (s *macaroonService) attenuate(token string, caveat Caveat) (string, Macaroon, error) {
	macaroon, err := decodeMacaroon(token)
	if err != nil {
		return "", Macaroon{}, err
	}
	if err := validateCaveat(caveat); err != nil {
		return "", Macaroon{}, err
	}
	current, err := hex.DecodeString(macaroon.Signature)
	if err != nil || len(current) != sha256.Size {
		return "", Macaroon{}, errors.New("invalid signature encoding")
	}
	macaroon.Caveats = append(macaroon.Caveats, caveat)
	macaroon.Signature = hex.EncodeToString(sign(current, caveat.Type+"="+caveat.Value))
	attenuated, err := encodeMacaroon(macaroon)
	return attenuated, macaroon, err
}

func (s *macaroonService) verify(token string, access AccessRequest) Verification {
	macaroon, err := decodeMacaroon(token)
	if err != nil {
		return Verification{Allowed: false, Reason: err.Error(), Checks: []CaveatCheck{}}
	}
	expected := sign(s.rootKey, macaroon.ID)
	for _, caveat := range macaroon.Caveats {
		expected = sign(expected, caveat.Type+"="+caveat.Value)
	}
	provided, err := hex.DecodeString(macaroon.Signature)
	if err != nil || !hmac.Equal(expected, provided) {
		return Verification{Allowed: false, Reason: "signature chain is invalid", Checks: []CaveatCheck{}}
	}
	checks := make([]CaveatCheck, 0, len(macaroon.Caveats))
	allowed := true
	for _, caveat := range macaroon.Caveats {
		check := s.checkCaveat(caveat, access)
		checks = append(checks, check)
		if !check.Passed {
			allowed = false
		}
	}
	reason := "every caveat passed"
	if !allowed {
		reason = "one or more caveats rejected the request"
	}
	return Verification{Allowed: allowed, Reason: reason, Checks: checks}
}

func (s *macaroonService) checkCaveat(caveat Caveat, access AccessRequest) CaveatCheck {
	check := CaveatCheck{Type: caveat.Type, Required: caveat.Value}
	switch caveat.Type {
	case "resource":
		check.Actual = access.Resource
		check.Passed = matches(caveat.Value, access.Resource)
	case "operation":
		check.Actual = access.Operation
		check.Passed = matches(caveat.Value, access.Operation)
	case "location":
		check.Actual = access.Location
		check.Passed = matches(caveat.Value, access.Location)
	case "expires":
		check.Actual = s.now().UTC().Format(time.RFC3339)
		expires, err := time.Parse(time.RFC3339, caveat.Value)
		check.Passed = err == nil && s.now().Before(expires)
	default:
		check.Actual = "unsupported"
		check.Passed = false
	}
	return check
}

func sign(key []byte, value string) []byte {
	hash := hmac.New(sha256.New, key)
	hash.Write([]byte(value))
	return hash.Sum(nil)
}

func matches(required, actual string) bool {
	if required == "*" {
		return true
	}
	if strings.HasSuffix(required, "*") {
		return strings.HasPrefix(actual, strings.TrimSuffix(required, "*"))
	}
	return required == actual
}

func validateCaveat(caveat Caveat) error {
	if strings.TrimSpace(caveat.Value) == "" {
		return errors.New("caveat value is required")
	}
	switch caveat.Type {
	case "resource":
		if strings.Count(caveat.Value, "*") > 1 || strings.Contains(strings.TrimSuffix(caveat.Value, "*"), "*") {
			return errors.New("resource wildcard must appear once at the end")
		}
		return nil
	case "operation", "location":
		if strings.Contains(caveat.Value, "*") && caveat.Value != "*" {
			return fmt.Errorf("%s wildcard must be the complete value", caveat.Type)
		}
		return nil
	case "expires":
		if _, err := time.Parse(time.RFC3339, caveat.Value); err != nil {
			return errors.New("expiration must use RFC3339")
		}
		return nil
	default:
		return fmt.Errorf("unsupported caveat type: %s", caveat.Type)
	}
}

func encodeMacaroon(macaroon Macaroon) (string, error) {
	data, err := json.Marshal(macaroon)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(data), nil
}

func decodeMacaroon(token string) (Macaroon, error) {
	data, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return Macaroon{}, errors.New("token is not valid base64url")
	}
	var macaroon Macaroon
	if err := json.Unmarshal(data, &macaroon); err != nil {
		return Macaroon{}, errors.New("token payload is not valid JSON")
	}
	if macaroon.Version != 1 || macaroon.ID == "" || macaroon.Signature == "" {
		return Macaroon{}, errors.New("token payload is incomplete")
	}
	return macaroon, nil
}
