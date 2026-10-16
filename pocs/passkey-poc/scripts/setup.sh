#!/usr/bin/env bash
set -euo pipefail
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/common.sh"

log "setup started"

require go
require bun

( cd "$ROOT/passkey-server" && go build -o "$ROOT/passkey-server/bin/passkey-server" ./cmd/passkey-server ) || fail "go build failed"
log "passkey-server built"

( cd "$ROOT" && bun install ) || fail "bun install failed"
( cd "$ROOT/e2e" && bunx playwright install chromium ) || fail "playwright chromium install failed"

mkdir -p "$ROOT/data"
chmod 700 "$ROOT/data"

log "setup done"
