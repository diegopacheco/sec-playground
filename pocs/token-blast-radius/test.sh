#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"
PORT="${PORT:-18082}"
BLAST_RADIUS_SECRET="test-blast-radius-key-32-bytes-minimum"
export PORT BLAST_RADIUS_SECRET
mkdir -p .run
cleanup() {
  if [ -f .run/test-server.pid ]; then
    PID=$(cat .run/test-server.pid)
    kill "$PID" 2>/dev/null || true
    wait "$PID" 2>/dev/null || true
    rm -f .run/test-server.pid
  fi
}
trap cleanup EXIT
cleanup
go test -count=1 ./...
cd frontend
bun install --frozen-lockfile
bun run typecheck
bun run build
cd ..
go build -o .run/test-server ./cmd/server
./.run/test-server > .run/test-server.log 2>&1 &
echo $! > .run/test-server.pid
for i in $(seq 1 30); do
  if curl -sf "http://localhost:${PORT}/health" >/dev/null; then
    break
  fi
  sleep 1
done
curl -sf "http://localhost:${PORT}/health" | grep -q '"status":"ready"'
SCENARIOS=$(curl -sf "http://localhost:${PORT}/api/scenarios")
TOKEN_CONTAINED=$(printf '%s' "$SCENARIOS" | grep -o '"token":"[^"]*"' | sed -n '1p' | cut -d '"' -f 4)
TOKEN_WIDE=$(printf '%s' "$SCENARIOS" | grep -o '"token":"[^"]*"' | sed -n '2p' | cut -d '"' -f 4)
test -n "$TOKEN_CONTAINED"
test -n "$TOKEN_WIDE"
curl -sf -X POST "http://localhost:${PORT}/api/analyze" -H "Content-Type: application/json" -d "{\"scenario_id\":\"01-contained\",\"token\":\"$TOKEN_CONTAINED\"}" | grep -q '"risk_level":"contained"'
curl -sf -X POST "http://localhost:${PORT}/api/analyze" -H "Content-Type: application/json" -d "{\"scenario_id\":\"02-over-delegated\",\"token\":\"$TOKEN_WIDE\"}" | grep -q '"risk_level":"critical"'
curl -s -X POST "http://localhost:${PORT}/api/analyze" -H "Content-Type: application/json" -d "{\"scenario_id\":\"01-contained\",\"token\":\"${TOKEN_CONTAINED}A\"}" | grep -q '"error":"token signature is invalid"'
curl -sf "http://localhost:${PORT}/" | grep -q '<title>Token Blast Radius</title>'
echo "signed claims, graph traversal, risk scoring, signature rejection, and UI build passed"
