#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"
source ../auth0-env.sh

cd webapp
npm install --silent
npm run build
cd ..

./stop.sh >/dev/null 2>&1 || true

nohup cargo run -p auth0-rust-sample-app > server.log 2>&1 &
echo $! > server.pid

for i in $(seq 1 120); do
  if curl -fsS -o /dev/null http://localhost:3000/api/me 2>/dev/null; then
    echo "rust sample app on http://localhost:3000"
    exit 0
  fi
  sleep 1
done

echo "sample app did not start, see sample-app/server.log"
exit 1
