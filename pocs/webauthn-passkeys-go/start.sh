#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"

PORT="${PORT:-8090}"
mkdir -p bin
go build -o bin/passkey-service .
PORT="$PORT" ./bin/passkey-service > server.log 2>&1 &
echo $! > server.pid

for i in $(seq 1 30); do
  if curl -sf "http://localhost:${PORT}/health" >/dev/null; then
    echo "passkey service is running at http://localhost:${PORT}"
    exit 0
  fi
  sleep 1
done

cat server.log
exit 1
