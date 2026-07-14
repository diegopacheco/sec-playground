#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"

PORT="${PORT:-8091}"
mkdir -p bin
rm -f audit.jsonl
go build -o bin/agent-firewall .
PORT="$PORT" ./bin/agent-firewall > server.log 2>&1 &
echo $! > server.pid

for i in $(seq 1 30); do
  if curl -sf "http://localhost:${PORT}/health" >/dev/null; then
    echo "agent firewall is running at http://localhost:${PORT}"
    exit 0
  fi
  sleep 1
done

cat server.log
exit 1
