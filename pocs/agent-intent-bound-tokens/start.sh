#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"
PORT="${PORT:-8081}"
INTENT_TOKEN_SECRET="${INTENT_TOKEN_SECRET:-local-intent-signing-key-32-bytes-minimum}"
export PORT INTENT_TOKEN_SECRET
mkdir -p .run
if [ -f .run/server.pid ] && kill -0 "$(cat .run/server.pid)" 2>/dev/null; then
  echo "intent token service is already running at http://localhost:${PORT}"
  exit 0
fi
go build -o .run/server ./cmd/server
./.run/server > .run/server.log 2>&1 &
echo $! > .run/server.pid
for i in $(seq 1 30); do
  if curl -sf "http://localhost:${PORT}/health" >/dev/null; then
    echo "intent token service is running at http://localhost:${PORT}"
    exit 0
  fi
  sleep 1
done
cat .run/server.log
exit 1
