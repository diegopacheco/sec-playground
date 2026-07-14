#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"

if [ -f server.pid ]; then
  kill "$(cat server.pid)" 2>/dev/null || true
  rm -f server.pid
  echo "passkey service stopped"
else
  echo "passkey service is not running"
fi
