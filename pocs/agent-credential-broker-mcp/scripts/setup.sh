#!/usr/bin/env bash
set -euo pipefail
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/common.sh"

require go
require podman
require curl
require openssl
require lsof

log "building binaries into bin/"
go build -o "$BIN/" ./cmd/...

if podman image exists "$VAULT_IMAGE"; then
  log "vault image present"
else
  log "pulling $VAULT_IMAGE"
  podman pull "$VAULT_IMAGE" >/dev/null || fail "could not pull $VAULT_IMAGE"
fi

ensure_file "$DATA/vault-token" "$(random_hex 24)"
ensure_file "$DATA/vendor-api-key" "sk_live_$(random_hex 20)"

log "setup complete"
