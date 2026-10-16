#!/usr/bin/env bash
set -euo pipefail
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/common.sh"

log "starting"

[ -x "$ROOT/backend/bin/server" ] || fail "backend not built, run ./scripts/setup.sh first"
[ -d "$ROOT/frontend/node_modules" ] || fail "frontend not installed, run ./scripts/setup.sh first"

start_bg backend "$ROOT/backend" env TAINT_PORT="$(service_port backend)" ./bin/server
wait_port_up "$(service_port backend)" 60 || fail "backend did not open port $(service_port backend), see $LOGS/backend.log"
log "backend up on $(service_url backend)/api"

start_bg frontend "$ROOT/frontend" env WEB_PORT="$(service_port frontend)" TRACKER_PORT="$(service_port backend)" bun run dev
wait_port_up "$(service_port frontend)" 60 || fail "frontend did not open port $(service_port frontend), see $LOGS/frontend.log"
log "frontend up on $(service_url frontend)"

"$SCRIPTS/status.sh"

log "links"
printf "%-14s %s\n" backend "$(service_url backend)/api"
printf "%-14s %s\n" frontend "$(service_url frontend)"
