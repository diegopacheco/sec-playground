#!/usr/bin/env bash
set -euo pipefail
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/common.sh"

[ -x "$BIN/broker" ] && [ -s "$DATA/vault-token" ] || "$SCRIPTS/setup.sh"

VAULT_PORT="$(service_port vault)"
VENDOR_PORT="$(service_port vendor)"
BROKER_PORT="$(service_port broker)"
VAULT_TOKEN_VALUE="$(cat "$DATA/vault-token")"
VENDOR_KEY="$(cat "$DATA/vendor-api-key")"

log "starting"

if [ "$(podman inspect -f '{{.State.Running}}' "$VAULT_CONTAINER" 2>/dev/null || true)" = "true" ]; then
  log "vault already running"
else
  podman rm -f "$VAULT_CONTAINER" >/dev/null 2>&1 || true
  podman run -d --name "$VAULT_CONTAINER" --cap-add=IPC_LOCK \
    -p "127.0.0.1:$VAULT_PORT:8200" \
    -e "VAULT_DEV_ROOT_TOKEN_ID=$VAULT_TOKEN_VALUE" \
    -e "VAULT_DEV_LISTEN_ADDRESS=0.0.0.0:8200" \
    "$VAULT_IMAGE" >/dev/null || fail "vault container did not start"
  log "vault started"
fi
wait_http "$(service_url vault)/v1/sys/health" 60 || fail "vault is not healthy on port $VAULT_PORT"

curl -sf -o /dev/null -X POST \
  -H "X-Vault-Token: $VAULT_TOKEN_VALUE" \
  -H "Content-Type: application/json" \
  -d "{\"data\":{\"api_key\":\"$VENDOR_KEY\"}}" \
  "$(service_url vault)/v1/secret/data/agent/vendor-api" || fail "could not seed vault"
log "vault up on $(service_url vault), vendor key stored at agent/vendor-api"

start_bg vendor env UPSTREAM_API_KEY="$VENDOR_KEY" UPSTREAM_ADDR="127.0.0.1:$VENDOR_PORT" "$BIN/upstream"
wait_port_up "$VENDOR_PORT" 60 || fail "vendor did not open port $VENDOR_PORT, see $LOGS/vendor.log"
log "vendor up on $(service_url vendor)"

start_bg broker env -i PATH="$PATH" HOME="$HOME" \
  BROKER_ADDR="127.0.0.1:$BROKER_PORT" \
  POLICY_FILE="$ROOT/policy.json" \
  DATA_DIR="$DATA" \
  VAULT_ADDR="$(service_url vault)" \
  VAULT_TOKEN="$VAULT_TOKEN_VALUE" \
  UPSTREAM_HOST="localhost:$VENDOR_PORT" \
  UPSTREAM_URL="$(service_url vendor)" \
  "$BIN/broker"
wait_port_up "$BROKER_PORT" 60 || fail "broker did not open port $BROKER_PORT, see $LOGS/broker.log"
wait_http "$(service_url broker)/health" 30 || fail "broker is not healthy"
log "broker up on $(service_url broker)"

"$SCRIPTS/status.sh"

log "links"
for name in $(service_names); do
  printf "%-10s %s\n" "$name" "$(service_url "$name")"
done
printf "%-10s %s\n" "ui" "$(ui_link)"
printf "%-10s %s\n" "mcp" "BROKER_URL=$(service_url broker) $BIN/mcp"
