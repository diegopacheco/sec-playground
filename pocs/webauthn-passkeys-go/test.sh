#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"

PORT="${PORT:-8090}"
DB_PATH="${DB_PATH:-/tmp/webauthn-passkeys-test.db}"
export PORT DB_PATH
rm -f "$DB_PATH"
./stop.sh >/dev/null 2>&1 || true
trap './stop.sh >/dev/null 2>&1 || true; rm -f "$DB_PATH"' EXIT

go test -v ./...
./start.sh

curl -sf "http://localhost:${PORT}/health"
echo
BEGIN=$(curl -sf -X POST "http://localhost:${PORT}/api/register/begin" -H "Content-Type: application/json" -d '{"username":"security-test","displayName":"Security Test"}')
echo "$BEGIN"
echo "$BEGIN" | grep -q '"publicKey"'
echo "$BEGIN" | grep -q '"userVerification":"required"'
echo "browser ceremony options require user verification"
