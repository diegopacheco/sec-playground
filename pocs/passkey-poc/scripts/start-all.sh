#!/usr/bin/env bash
set -euo pipefail
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/common.sh"

export_ports
[ -x "$ROOT/passkey-server/bin/passkey-server" ] || fail "passkey-server is not built, run ./scripts/setup.sh first"
mkdir -p "$ROOT/data"

log "starting"

apps="bookstore=Passkey Books=http://localhost:$BOOKSTORE_WEB_PORT;vault=PasskeyVault=http://localhost:$VAULT_WEB_PORT"
start_service passkey "$ROOT" "$ROOT/passkey-server/bin/passkey-server" -addr "127.0.0.1:$PASSKEY_PORT" -rp-id localhost -apps "$apps" -data "$ROOT/data/passkeys.json"
start_service bookstore-api "$ROOT/apps/bookstore/api" bun src/server.ts
start_service vault-api "$ROOT/apps/vault/api" env VAULT_DATA="$ROOT/data/vault.json" bun src/server.ts
start_service bookstore-web "$ROOT/apps/bookstore/web" bun run dev
start_service vault-web "$ROOT/apps/vault/web" bun run dev

"$SCRIPTS/status.sh"

log "links"
for name in $(service_names); do
  printf "%-14s %s\n" "$name" "$(service_url "$name")"
done
