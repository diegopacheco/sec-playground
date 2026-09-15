package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"mcp-rug-pull-proxy/internal/sim"
)

func main() {
	mode := flag.String("mode", "rugpull", "clean, poisoned or rugpull")
	guardURL := flag.String("guard", "http://localhost:8741", "guard base url")
	runDir := flag.String("run-dir", ".run", "directory holding admin-token and proxy-token")
	binDir := flag.String("bin", "bin", "directory holding the proxy and fake-upstream binaries")
	server := flag.String("server", "", "server name, defaults to weather-<mode>")
	flag.Parse()
	if *server == "" {
		*server = "weather-" + *mode
	}

	adminToken, err := read(filepath.Join(*runDir, "admin-token"))
	if err != nil {
		fail(err)
	}
	proxyTokenFile := filepath.Join(*runDir, "proxy-token")

	cmd := exec.Command(filepath.Join(*binDir, "proxy"),
		"--server", *server, "--guard", *guardURL, "--token-file", proxyTokenFile,
		"--", filepath.Join(*binDir, "fake-upstream"), "--mode", *mode)
	cmd.Stderr = os.Stderr
	toProxy, err := cmd.StdinPipe()
	if err != nil {
		fail(err)
	}
	fromProxy, err := cmd.StdoutPipe()
	if err != nil {
		fail(err)
	}
	if err := cmd.Start(); err != nil {
		fail(err)
	}

	sc := &sim.Scenario{
		Mode:   *mode,
		Server: *server,
		Client: sim.NewClient(fromProxy, toProxy),
		Admin:  sim.NewAdmin(strings.TrimRight(*guardURL, "/"), adminToken),
		Out:    os.Stdout,
	}
	runErr := sc.Run()
	_ = toProxy.Close()
	_ = cmd.Wait()
	if runErr != nil {
		fmt.Fprintf(os.Stderr, "FAIL %s: %v\n", *mode, runErr)
		os.Exit(1)
	}
}

func read(path string) (string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(raw)), nil
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "client-sim:", err)
	os.Exit(1)
}
