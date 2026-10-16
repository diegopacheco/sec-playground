#!/usr/bin/env bash
set -euo pipefail
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/common.sh"

log "stopping"

for name in vault-web bookstore-web vault-api bookstore-api passkey; do
  stop_bg "$name"
done

for name in $(service_names); do
  port="$(service_port "$name")"
  if port_up "$port"; then
    fail "$name still listening on $port"
  fi
done

log "stopped"
