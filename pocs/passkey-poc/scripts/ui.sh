#!/usr/bin/env bash
set -euo pipefail
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/common.sh"

for name in bookstore-web vault-web; do
  port="$(service_port "$name")"
  [ -n "$port" ] || fail "$name is not declared in scripts/ports.env"
  port_up "$port" || fail "$name is not running on $port, run ./scripts/start-all.sh first"
  url="$(service_url "$name")"
  log "opening $url"
  open_url "$url"
done
