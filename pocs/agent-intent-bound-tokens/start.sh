#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"
PORT="${PORT:-8081}"
UI_PORT="${UI_PORT:-5174}"
INTENT_TOKEN_SECRET="${INTENT_TOKEN_SECRET:-local-intent-signing-key-32-bytes-minimum}"
export PORT UI_PORT INTENT_TOKEN_SECRET
mkdir -p .run

running() {
  [ -f "$1" ] && kill -0 "$(cat "$1")" 2>/dev/null
}

if running .run/server.pid; then
  BACKEND_PORT="$(cat .run/server.port 2>/dev/null || printf '%s' "$PORT")"
else
  rm -f .run/server.pid .run/server.port
  go build -o .run/server ./cmd/server
  ./.run/server > .run/server.log 2>&1 &
  echo $! > .run/server.pid
  printf '%s' "$PORT" > .run/server.port
  BACKEND_PORT="$PORT"
fi

if [ ! -x ui/node_modules/.bin/vite ]; then
  (cd ui && bun install --frozen-lockfile)
fi

if running .run/ui.pid; then
  FRONTEND_PORT="$(cat .run/ui.port 2>/dev/null || printf '%s' "$UI_PORT")"
else
  rm -f .run/ui.pid .run/ui.port
  (cd ui && exec bun run dev -- --strictPort) > .run/ui.log 2>&1 &
  echo $! > .run/ui.pid
  printf '%s' "$UI_PORT" > .run/ui.port
  FRONTEND_PORT="$UI_PORT"
fi

for i in $(seq 1 30); do
  if curl -sf "http://localhost:${BACKEND_PORT}/health" >/dev/null && curl -sf "http://localhost:${FRONTEND_PORT}/" >/dev/null; then
    echo "backend is running at http://localhost:${BACKEND_PORT}"
    echo "frontend is running at http://localhost:${FRONTEND_PORT}"
    exit 0
  fi
  if ! running .run/server.pid || ! running .run/ui.pid; then
    break
  fi
  sleep 1
done

cat .run/server.log 2>/dev/null || true
cat .run/ui.log 2>/dev/null || true
exit 1
