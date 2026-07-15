#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"
PORT="${PORT:-18081}"
INTENT_TOKEN_SECRET="test-intent-signing-key-32-bytes-minimum"
export PORT INTENT_TOKEN_SECRET
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
curl -sf "http://localhost:${PORT}/" | grep -q '"path":"/api/tokens"'
ISSUED=$(curl -sf -X POST "http://localhost:${PORT}/api/tokens" -H "Content-Type: application/json" -d '{"subject":"agent:buyer","audience":"payments-api","action":"transfer","resource":"account:operations","max_amount_cents":50000,"valid_for_seconds":60}')
TOKEN=$(printf '%s' "$ISSUED" | sed -n 's/.*"token":"\([^"]*\)".*/\1/p')
test -n "$TOKEN"
curl -s -X POST "http://localhost:${PORT}/api/verify" -H "Content-Type: application/json" -d "{\"token\":\"$TOKEN\",\"audience\":\"payments-api\",\"action\":\"delete\",\"resource\":\"account:operations\",\"amount_cents\":100}" | grep -q '"code":"action_mismatch"'
curl -s -X POST "http://localhost:${PORT}/api/verify" -H "Content-Type: application/json" -d "{\"token\":\"$TOKEN\",\"audience\":\"payments-api\",\"action\":\"transfer\",\"resource\":\"account:operations\",\"amount_cents\":50001}" | grep -q '"code":"amount_exceeded"'
curl -sf -X POST "http://localhost:${PORT}/api/verify" -H "Content-Type: application/json" -d "{\"token\":\"$TOKEN\",\"audience\":\"payments-api\",\"action\":\"transfer\",\"resource\":\"account:operations\",\"amount_cents\":49999}" | grep -q '"allowed":true'
curl -s -X POST "http://localhost:${PORT}/api/verify" -H "Content-Type: application/json" -d "{\"token\":\"$TOKEN\",\"audience\":\"payments-api\",\"action\":\"transfer\",\"resource\":\"account:operations\",\"amount_cents\":49999}" | grep -q '"code":"already_used"'
curl -sf "http://localhost:${PORT}/api/audit" | grep -q '"allowed":true'
echo "intent constraints, amount limit, one-time use, and audit records passed"
