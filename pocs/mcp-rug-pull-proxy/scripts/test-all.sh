#!/usr/bin/env bash
set -euo pipefail
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/common.sh"

log "tests started"

require go
go vet ./... || fail "go vet failed"
go test -race -count=1 ./... || fail "go tests failed"

started_here=0
if ! port_up "$(service_port guard)"; then
  "$SCRIPTS/start-all.sh" >/dev/null
  started_here=1
else
  build_binaries
fi

cleanup() {
  if [ "$started_here" = 1 ]; then
    "$SCRIPTS/stop-all.sh" >/dev/null
  fi
}
trap cleanup EXIT

for mode in clean poisoned rugpull; do
  "$BIN/client-sim" --mode "$mode" --guard "$(service_url guard)" --run-dir "$RUN" --bin "$BIN" 2>>"$LOGS/client-sim.log" || fail "end to end scenario $mode failed, see .run/logs/client-sim.log"
done

log "tests passed"
