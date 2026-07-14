#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"

PORT="${PORT:-8091}"
export PORT
./stop.sh >/dev/null 2>&1 || true
trap './stop.sh >/dev/null 2>&1 || true' EXIT

go test -v ./...
./start.sh

SAFE=$(curl -sf -X POST "http://localhost:${PORT}/run" -H "Content-Type: application/json" -d '{"document":"safe.txt"}')
echo "$SAFE"
echo "$SAFE" | grep -q '"allowed":true'

POISONED=$(curl -sf -X POST "http://localhost:${PORT}/run" -H "Content-Type: application/json" -d '{"document":"poisoned.txt"}')
echo "$POISONED"
echo "$POISONED" | grep -q 'read target outside capability root'
echo "$POISONED" | grep -q 'HTTP destination not allowed'
echo "$POISONED" | grep -q 'sensitive canary detected'

DENIED=$(curl -sf -X POST "http://localhost:${PORT}/run" -H "Content-Type: application/json" -d '{"document":"approval.txt"}')
echo "$DENIED"
echo "$DENIED" | grep -q 'explicit approval required'

APPROVED=$(curl -sf -X POST "http://localhost:${PORT}/run" -H "Content-Type: application/json" -d '{"document":"approval.txt","approval":"approve-once"}')
echo "$APPROVED"
echo "$APPROVED" | grep -q 'approved HTTP write allowed'
echo "$APPROVED" | grep -q '"allowed":true'

AUDIT=$(curl -sf "http://localhost:${PORT}/audit")
echo "$AUDIT"
echo "$AUDIT" | grep -q '\[REDACTED\]'
if echo "$AUDIT" | grep -q 'FIREWALL_CANARY_7f3c91'; then
  exit 1
fi

echo "all firewall protections passed"
