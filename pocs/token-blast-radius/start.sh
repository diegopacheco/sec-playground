#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"
PORT="${PORT:-8082}"
BLAST_RADIUS_SECRET="${BLAST_RADIUS_SECRET:-local-blast-radius-key-32-bytes-minimum}"
export PORT BLAST_RADIUS_SECRET
mkdir -p .run
if [ -f .run/server.pid ] && kill -0 "$(cat .run/server.pid)" 2>/dev/null; then
  echo "token blast radius is already running at http://localhost:${PORT}"
  exit 0
fi
if [ ! -d frontend/node_modules ]; then
  cd frontend
  bun install --frozen-lockfile
  cd ..
fi
cd frontend
bun run build
cd ..
go build -o .run/server ./cmd/server
./.run/server > .run/server.log 2>&1 &
echo $! > .run/server.pid
for i in $(seq 1 30); do
  if curl -sf "http://localhost:${PORT}/health" >/dev/null; then
    echo "token blast radius is running at http://localhost:${PORT}"
    exit 0
  fi
  sleep 1
done
cat .run/server.log
exit 1
