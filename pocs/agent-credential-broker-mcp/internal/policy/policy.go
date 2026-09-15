package policy

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

type Rule struct {
	Scheme     string   `json:"scheme"`
	Host       string   `json:"host"`
	Methods    []string `json:"methods"`
	PathPrefix string   `json:"path_prefix"`
}

type Inject struct {
	Header   string `json:"header"`
	Template string `json:"template"`
}

type Credential struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Provider    string `json:"provider"`
	Path        string `json:"path,omitempty"`
	Field       string `json:"field,omitempty"`
	Service     string `json:"service,omitempty"`
	Account     string `json:"account,omitempty"`
	Inject      Inject `json:"inject"`
	Allow       []Rule `json:"allow"`
	Approval    string `json:"approval"`
}

type Policy struct {
	Credentials []Credential `json:"credentials"`
}

func Load(path string) (*Policy, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Parse([]byte(os.ExpandEnv(string(raw))))
}

func Parse(raw []byte) (*Policy, error) {
	var p Policy
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, err
	}
	return &p, p.Validate()
}

func (p *Policy) Validate() error {
	seen := map[string]bool{}
	for _, c := range p.Credentials {
		if c.Name == "" {
			return errors.New("credential without name")
		}
		if seen[c.Name] {
			return fmt.Errorf("duplicate credential %q", c.Name)
		}
		seen[c.Name] = true
		if c.Provider != "vault" && c.Provider != "keychain" {
			return fmt.Errorf("credential %q has unknown provider %q", c.Name, c.Provider)
		}
		if c.Approval != "human" && c.Approval != "auto" {
			return fmt.Errorf("credential %q approval must be human or auto", c.Name)
		}
		if c.Inject.Header == "" || !strings.Contains(c.Inject.Template, "{{credential}}") {
			return fmt.Errorf("credential %q needs inject header and a template with {{credential}}", c.Name)
		}
		if len(c.Allow) == 0 {
			return fmt.Errorf("credential %q has no allow rules", c.Name)
		}
	}
	return nil
}

func (p *Policy) Find(name string) (Credential, bool) {
	for _, c := range p.Credentials {
		if c.Name == name {
			return c, true
		}
	}
	return Credential{}, false
}
