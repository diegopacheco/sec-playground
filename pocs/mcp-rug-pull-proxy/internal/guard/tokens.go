package guard

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func LoadOrCreateTokens(runDir string) (Tokens, error) {
	if err := os.MkdirAll(runDir, 0o700); err != nil {
		return Tokens{}, err
	}
	admin, err := loadOrCreate(filepath.Join(runDir, "admin-token"))
	if err != nil {
		return Tokens{}, err
	}
	proxy, err := loadOrCreate(filepath.Join(runDir, "proxy-token"))
	if err != nil {
		return Tokens{}, err
	}
	return Tokens{Admin: admin, Proxy: proxy}, nil
}

func loadOrCreate(path string) (string, error) {
	raw, err := os.ReadFile(path)
	if err == nil && len(strings.TrimSpace(string(raw))) >= 32 {
		return strings.TrimSpace(string(raw)), nil
	}
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return "", err
	}
	b := make([]byte, 24)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	token := hex.EncodeToString(b)
	return token, os.WriteFile(path, []byte(token+"\n"), 0o600)
}
