package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"agent-credential-broker-mcp/internal/policy"
)

type Source interface {
	Fetch(ctx context.Context, c policy.Credential) (string, error)
}

type Vault struct {
	Addr   string
	Token  string
	Client *http.Client
}

func (v Vault) Fetch(ctx context.Context, c policy.Credential) (string, error) {
	if v.Addr == "" || v.Token == "" {
		return "", errors.New("vault address or token not configured")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, strings.TrimRight(v.Addr, "/")+"/v1/secret/data/"+c.Path, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("X-Vault-Token", v.Token)
	client := v.Client
	if client == nil {
		client = &http.Client{Timeout: 5 * time.Second}
	}
	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("vault unreachable: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("vault returned %d for %s", res.StatusCode, c.Path)
	}
	var body struct {
		Data struct {
			Data map[string]string `json:"data"`
		} `json:"data"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return "", err
	}
	value, ok := body.Data.Data[c.Field]
	if !ok || value == "" {
		return "", fmt.Errorf("vault field %q missing at %s", c.Field, c.Path)
	}
	return value, nil
}

type Runner func(ctx context.Context, name string, args ...string) ([]byte, error)

func ExecRunner(ctx context.Context, name string, args ...string) ([]byte, error) {
	return exec.CommandContext(ctx, name, args...).Output()
}

type Keychain struct {
	Run Runner
}

func (k Keychain) Fetch(ctx context.Context, c policy.Credential) (string, error) {
	run := k.Run
	if run == nil {
		run = ExecRunner
	}
	out, err := run(ctx, "security", "find-generic-password", "-s", c.Service, "-a", c.Account, "-w")
	if err != nil {
		return "", fmt.Errorf("keychain item %s/%s not found", c.Service, c.Account)
	}
	value := strings.TrimRight(string(out), "\r\n")
	if value == "" {
		return "", errors.New("keychain item is empty")
	}
	return value, nil
}

type Router struct {
	Vault    Source
	Keychain Source
}

func (r Router) Fetch(ctx context.Context, c policy.Credential) (string, error) {
	switch c.Provider {
	case "vault":
		return r.Vault.Fetch(ctx, c)
	case "keychain":
		return r.Keychain.Fetch(ctx, c)
	}
	return "", fmt.Errorf("unknown provider %q", c.Provider)
}
