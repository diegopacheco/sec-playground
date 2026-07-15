package token

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"token-blast-radius/internal/domain"
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

func (s *Service) Issue(claims domain.Claims) (string, error) {
	header, err := encode(map[string]string{"alg": "HS256", "typ": "JWT"})
	if err != nil {
		return "", err
	}
	payload, err := encode(claims)
	if err != nil {
		return "", err
	}
	unsigned := header + "." + payload
	return unsigned + "." + s.sign(unsigned), nil
}

func (s *Service) Parse(value string) (domain.Claims, error) {
	parts := strings.Split(value, ".")
	if len(parts) != 3 {
		return domain.Claims{}, errors.New("token format is invalid")
	}
	unsigned := parts[0] + "." + parts[1]
	expected, _ := base64.RawURLEncoding.DecodeString(s.sign(unsigned))
	provided, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil || !hmac.Equal(expected, provided) {
		return domain.Claims{}, errors.New("token signature is invalid")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return domain.Claims{}, errors.New("token payload is invalid")
	}
	var claims domain.Claims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return domain.Claims{}, errors.New("token payload is invalid")
	}
	if claims.Subject == "" || claims.Actor == "" || claims.Audience == "" || claims.MaximumDepth < 1 {
		return domain.Claims{}, errors.New("token claims are incomplete")
	}
	if s.now().Unix() >= claims.ExpiresAt {
		return domain.Claims{}, errors.New("token is expired")
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
