package sim

import (
	"bytes"
	"io"
	"net/http/httptest"
	"testing"

	"mcp-rug-pull-proxy/internal/fakeserver"
	"mcp-rug-pull-proxy/internal/guard"
	"mcp-rug-pull-proxy/internal/proxy"
	"mcp-rug-pull-proxy/internal/store"
)

func TestEndToEndScenariosThroughGuardProxyAndRealMcpServer(t *testing.T) {
	st, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	tokens := guard.Tokens{Admin: "admin-token-for-tests-0123456789abcdef", Proxy: "proxy-token-for-tests-0123456789abcdef"}
	srv := httptest.NewServer(guard.New(st, tokens).Handler())
	defer srv.Close()
	for _, mode := range []string{fakeserver.ModeClean, fakeserver.ModePoisoned, fakeserver.ModeRugPull} {
		t.Run(mode, func(t *testing.T) {
			clientToProxyR, clientToProxyW := io.Pipe()
			proxyToClientR, proxyToClientW := io.Pipe()
			proxyToUpR, proxyToUpW := io.Pipe()
			upToProxyR, upToProxyW := io.Pipe()
			go func() { _ = fakeserver.Serve(mode, proxyToUpR, upToProxyW) }()
			p := proxy.New("weather-"+mode, proxy.NewHTTPGuard(srv.URL, tokens.Proxy), io.Discard)
			go func() { _ = p.Run(clientToProxyR, proxyToClientW, proxyToUpW, upToProxyR) }()
			var out bytes.Buffer
			sc := &Scenario{Mode: mode, Server: "weather-" + mode, Client: NewClient(proxyToClientR, clientToProxyW), Admin: NewAdmin(srv.URL, tokens.Admin), Out: &out}
			if err := sc.Run(); err != nil {
				t.Fatalf("%v\n%s", err, out.String())
			}
			_ = clientToProxyW.Close()
			_ = proxyToUpW.Close()
		})
	}
	if _, chain := st.Events(); !chain.Verified {
		t.Fatalf("event chain broken after scenarios: %+v", chain)
	}
}
