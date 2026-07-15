package token

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

	"agent-intent-bound-tokens/internal/intent"
)

var (
	ErrFormat    = errors.New("token format is invalid")
	ErrSignature = errors.New("token signature is invalid")
	ErrNotActive = errors.New("token is not active")
	ErrExpired   = errors.New("token is expired")
)

type Service struct {
	secret []byte
	now    func() time.Time
}

func New(secret string) (*Service, error) {
	if len(secret) < 32 {
		return nil, errors.New("signing secret must contain at least 32 characters")
	}
	return &Service{secret: []byte(secret), now: time.Now}, nil
}

func (s *Service) Issue(claims intent.Claims) (string, intent.Claims, error) {
	if err := validateClaims(claims); err != nil {
		return "", intent.Claims{}, err
	}
	if claims.ID == "" {
		id, err := randomID()
		if err != nil {
			return "", intent.Claims{}, err
		}
		claims.ID = id
	}
	header, err := encode(map[string]string{"alg": "HS256", "typ": "JWT"})
	if err != nil {
		return "", intent.Claims{}, err
	}
	payload, err := encode(claims)
	if err != nil {
		return "", intent.Claims{}, err
	}
	unsigned := header + "." + payload
	return unsigned + "." + s.sign(unsigned), claims, nil
}

func (s *Service) Parse(value string) (intent.Claims, error) {
	parts := strings.Split(value, ".")
	if len(parts) != 3 {
		return intent.Claims{}, ErrFormat
	}
	unsigned := parts[0] + "." + parts[1]
	expected, err := base64.RawURLEncoding.DecodeString(s.sign(unsigned))
	if err != nil {
		return intent.Claims{}, ErrFormat
	}
	provided, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil || !hmac.Equal(expected, provided) {
		return intent.Claims{}, ErrSignature
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return intent.Claims{}, ErrFormat
	}
	var claims intent.Claims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return intent.Claims{}, ErrFormat
	}
	if err := validateClaims(claims); err != nil {
		return intent.Claims{}, ErrFormat
	}
	now := s.now().Unix()
	if now < claims.NotBefore {
		return intent.Claims{}, ErrNotActive
	}
	if now >= claims.ExpiresAt {
		return intent.Claims{}, ErrExpired
	}
	return claims, nil
}

func (s *Service) sign(value string) string {
	mac := hmac.New(sha256.New, s.secret)
	_, _ = mac.Write([]byte(value))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func encode(value any) (string, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(data), nil
}

func randomID() (string, error) {
	value := make([]byte, 16)
	if _, err := rand.Read(value); err != nil {
		return "", err
	}
	return hex.EncodeToString(value), nil
}

func validateClaims(claims intent.Claims) error {
	if claims.Issuer == "" || claims.Subject == "" || claims.Audience == "" || claims.Action == "" || claims.Resource == "" {
		return errors.New("issuer, subject, audience, action, and resource are required")
	}
	if claims.MaxAmountCents < 0 {
		return errors.New("maximum amount cannot be negative")
	}
	if claims.NotBefore <= 0 || claims.ExpiresAt <= claims.NotBefore {
		return fmt.Errorf("deadline must be after activation time")
	}
	return nil
}
