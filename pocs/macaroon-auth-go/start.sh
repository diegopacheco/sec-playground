#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"
mkdir -p .run
if [[ -f .run/server.pid ]] && kill -0 "$(cat .run/server.pid)" 2>/dev/null; then
  echo "macaroon-auth-go is already running"
  exit 0
fi
go build -o .run/macaroon-auth-go .
./.run/macaroon-auth-go > .run/server.log 2>&1 &
echo $! > .run/server.pid
for _ in {1..60}; do
  if curl -fsS http://127.0.0.1:8092/api/health >/dev/null 2>&1; then
    echo "macaroon-auth-go started at http://127.0.0.1:8092"
    exit 0
  fi
  sleep 1
done
echo "macaroon-auth-go failed to start"
cat .run/server.log
exit 1
