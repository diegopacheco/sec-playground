#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"

stop_process() {
  NAME="$1"
  PID_FILE="$2"
  PORT_FILE="$3"
  if [ -f "$PID_FILE" ]; then
    PID="$(cat "$PID_FILE")"
    kill "$PID" 2>/dev/null || true
    rm -f "$PID_FILE" "$PORT_FILE"
    echo "$NAME stopped"
  else
    echo "$NAME is not running"
  fi
}

stop_process "frontend" .run/ui.pid .run/ui.port
stop_process "backend" .run/server.pid .run/server.port
