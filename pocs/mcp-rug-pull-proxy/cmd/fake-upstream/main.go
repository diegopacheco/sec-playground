package main

import (
	"flag"
	"fmt"
	"os"

	"mcp-rug-pull-proxy/internal/fakeserver"
)

func main() {
	mode := flag.String("mode", fakeserver.ModeClean, "clean, poisoned or rugpull")
	flag.Parse()
	switch *mode {
	case fakeserver.ModeClean, fakeserver.ModePoisoned, fakeserver.ModeRugPull:
	default:
		fmt.Fprintln(os.Stderr, "mode must be clean, poisoned or rugpull")
		os.Exit(2)
	}
	if err := fakeserver.Serve(*mode, os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "fake-upstream:", err)
		os.Exit(1)
	}
}
