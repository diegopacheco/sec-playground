#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"
if [ -f .run/server.pid ]; then
  kill "$(cat .run/server.pid)" 2>/dev/null || true
  rm -f .run/server.pid
  echo "token blast radius stopped"
else
  echo "token blast radius is not running"
fi
