package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"mcp-rug-pull-proxy/internal/proxy"
)

func main() {
	server := flag.String("server", "", "name of the upstream server, used as the pin namespace")
	guardURL := flag.String("guard", "http://localhost:8741", "guard base url")
	tokenFile := flag.String("token-file", "", "file holding the proxy token, defaults to GUARD_PROXY_TOKEN")
	flag.Parse()
	upstream := flag.Args()
	if *server == "" || len(upstream) == 0 {
		fmt.Fprintln(os.Stderr, "usage: proxy --server <name> [--guard url] [--token-file path] -- <upstream command> [args]")
		os.Exit(2)
	}
	token, err := readToken(*tokenFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "proxy token:", err)
		os.Exit(2)
	}

	cmd := exec.Command(upstream[0], upstream[1:]...)
	cmd.Stderr = os.Stderr
	upIn, err := cmd.StdinPipe()
	if err != nil {
		fail(err)
	}
	upOut, err := cmd.StdoutPipe()
	if err != nil {
		fail(err)
	}
	if err := cmd.Start(); err != nil {
		fail(err)
	}

	p := proxy.New(*server, proxy.NewHTTPGuard(strings.TrimRight(*guardURL, "/"), token), os.Stderr)
	runErr := p.Run(os.Stdin, os.Stdout, upIn, upOut)
	_ = upIn.Close()
	done := make(chan struct{})
	go func() {
		_ = cmd.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		_ = cmd.Process.Kill()
	}
	if runErr != nil {
		fail(runErr)
	}
}

func readToken(path string) (string, error) {
	if path == "" {
		if v := os.Getenv("GUARD_PROXY_TOKEN"); v != "" {
			return v, nil
		}
		return "", fmt.Errorf("pass --token-file or set GUARD_PROXY_TOKEN")
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(raw)), nil
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "proxy:", err)
	os.Exit(1)
}
