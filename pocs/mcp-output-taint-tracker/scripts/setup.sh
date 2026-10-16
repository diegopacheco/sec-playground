#!/usr/bin/env bash
set -euo pipefail
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/common.sh"

log "setup started"

require go
( cd "$ROOT/backend" && go mod download && go build -o bin/ ./cmd/... ) || fail "backend build failed"
log "backend built: backend/bin/server and backend/bin/hook"

require bun
( cd "$ROOT/frontend" && bun install ) || fail "frontend install failed"

log "setup done"
