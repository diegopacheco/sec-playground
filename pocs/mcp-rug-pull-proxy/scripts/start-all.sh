#!/usr/bin/env bash
set -euo pipefail
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/common.sh"

log "starting"

port="$(service_port guard)"
if port_up "$port"; then
  log "guard already listening on $port"
else
  build_binaries
  start_bg guard "$ROOT" "$BIN/guard" --addr "127.0.0.1:$port" --run-dir "$RUN"
  wait_port_up "$port" 60 || fail "guard did not open port $port, see .run/logs/guard.log"
fi
log "guard up on $(service_url guard)"

"$SCRIPTS/status.sh"

log "links"
for name in $(service_names); do
  printf "%-14s %s\n" "$name" "$(service_url "$name")"
done
printf "%-14s %s\n" "review ui" "$(review_url)"
