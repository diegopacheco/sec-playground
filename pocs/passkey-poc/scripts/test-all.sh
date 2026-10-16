#!/usr/bin/env bash
set -euo pipefail
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/common.sh"

export_ports
log "tests started"

require go
require bun

log "go tests"
( cd "$ROOT/passkey-server" && go vet ./... && go test ./... ) || fail "go tests failed"

log "typecheck"
( cd "$ROOT" && bunx tsc -p tsconfig.json --noEmit ) || fail "typecheck failed"

log "bun tests"
( cd "$ROOT" && bun test packages apps ) || fail "bun tests failed"

log "e2e tests"
all_up=1
for name in $(service_names); do
  port_up "$(service_port "$name")" || all_up=0
done
if [ "$all_up" -eq 1 ]; then
  ( cd "$ROOT/e2e" && bun passkey.e2e.ts ) || fail "e2e tests failed"
else
  fail "e2e tests need the running app, run ./scripts/start-all.sh first"
fi

log "tests passed"
