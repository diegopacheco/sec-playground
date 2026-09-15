#!/usr/bin/env bash
set -euo pipefail
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/common.sh"

require go

log "go vet"
go vet ./... || fail "go vet failed"
go vet -tags e2e ./e2e || fail "go vet e2e failed"

log "unit and integration tests"
go test -count=1 ./... || fail "go tests failed"

for name in $(service_names); do
  port_up "$(service_port "$name")" || fail "$name is down, run ./scripts/start-all.sh before the live e2e suite"
done

log "live e2e through bin/mcp, broker, vault and vendor"
go build -o "$BIN/" ./cmd/...
BROKER_URL="$(service_url broker)" \
VENDOR_URL="$(service_url vendor)" \
ADMIN_TOKEN_FILE="$DATA/admin-token" \
VENDOR_KEY_FILE="$DATA/vendor-api-key" \
MCP_BIN="$BIN/mcp" \
go test -count=1 -tags e2e -v ./e2e || fail "live e2e failed"

log "all tests passed"
