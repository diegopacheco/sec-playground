#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"
if [[ ! -f .run/server.pid ]]; then
  echo "macaroon-auth-go is not running"
  exit 0
fi
pid="$(cat .run/server.pid)"
if kill -0 "$pid" 2>/dev/null; then
  kill "$pid"
  for _ in {1..60}; do
    if ! kill -0 "$pid" 2>/dev/null; then
      break
    fi
    sleep 1
  done
fi
rm -f .run/server.pid
echo "macaroon-auth-go stopped"
