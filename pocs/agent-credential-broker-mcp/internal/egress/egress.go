package egress

import (
	"errors"
	"fmt"
	"net/url"
	"path"
	"slices"
	"strings"

	"agent-credential-broker-mcp/internal/policy"
)

func Normalize(raw string) (*url.URL, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("invalid url: %w", err)
	}
	if u.Scheme == "" || u.Host == "" {
		return nil, errors.New("url must be absolute")
	}
	if u.User != nil {
		return nil, errors.New("userinfo in url is not allowed")
	}
	if u.RawPath != "" {
		return nil, errors.New("encoded characters in path are not allowed")
	}
	p := u.Path
	if p == "" {
		p = "/"
	}
	cleaned := path.Clean(p)
	if strings.HasSuffix(p, "/") && cleaned != "/" {
		cleaned += "/"
	}
	if cleaned != p {
		return nil, errors.New("dot segments in path are not allowed")
	}
	u.Scheme = strings.ToLower(u.Scheme)
	u.Host = strings.ToLower(u.Host)
	u.Path = p
	u.Fragment = ""
	return u, nil
}

func Check(c policy.Credential, method, raw string) (*url.URL, error) {
	u, err := Normalize(raw)
	if err != nil {
		return nil, err
	}
	method = strings.ToUpper(method)
	for _, r := range c.Allow {
		if r.Scheme != u.Scheme || !strings.EqualFold(r.Host, u.Host) {
			continue
		}
		if !slices.Contains(r.Methods, method) {
			continue
		}
		if !strings.HasPrefix(u.Path, r.PathPrefix) {
			continue
		}
		return u, nil
	}
	return nil, fmt.Errorf("%s %s://%s%s is not allowed for credential %q", method, u.Scheme, u.Host, u.Path, c.Name)
}
