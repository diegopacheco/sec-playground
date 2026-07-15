#!/usr/bin/env bash
set -uo pipefail
cd "$(dirname "$0")"

holders() {
  lsof -nP -ti tcp:3000 2>/dev/null
}

report() {
  echo "kotlin sample app stopped"
  exit 0
}

if [ -f server.pid ]; then
  kill "$(cat server.pid)" 2>/dev/null
  rm -f server.pid
fi

pids=$(holders)
[ -z "$pids" ] && report
echo "$pids" | xargs kill 2>/dev/null

for i in $(seq 1 10); do
  [ -z "$(holders)" ] && report
  sleep 1
done

echo "$(holders)" | xargs kill -9 2>/dev/null

for i in $(seq 1 5); do
  [ -z "$(holders)" ] && report
  sleep 1
done

echo "port 3000 still held by: $(holders)"
exit 1
