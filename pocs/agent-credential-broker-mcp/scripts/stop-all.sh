#!/usr/bin/env bash
set -euo pipefail
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/common.sh"

log "stopping"

stop_bg broker
stop_bg vendor

if podman container exists "$VAULT_CONTAINER" 2>/dev/null; then
  podman rm -f "$VAULT_CONTAINER" >/dev/null
fi
wait_port_down "$(service_port vault)" 30 || true
log "vault stopped"

for name in $(service_names); do
  port="$(service_port "$name")"
  if port_up "$port"; then
    fail "$name still listening on $port"
  fi
done

log "stopped"
