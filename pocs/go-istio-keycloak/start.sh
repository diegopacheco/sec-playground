#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"

PORT="${PORT:-8080}"
DB_PATH="${DB_PATH:-books.db}"

go build -o bin/book-service ./cmd/server

PORT="$PORT" DB_PATH="$DB_PATH" ./bin/book-service > server.log 2>&1 &
echo $! > server.pid

for i in $(seq 1 30); do
  if curl -sf "http://localhost:${PORT}/health" > /dev/null; then
    echo "book-service is up on http://localhost:${PORT}"
    exit 0
  fi
  sleep 1
done

echo "book-service failed to start"
cat server.log
exit 1
