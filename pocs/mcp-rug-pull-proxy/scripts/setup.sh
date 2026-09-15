#!/usr/bin/env bash
set -euo pipefail
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/common.sh"

log "setup started"

require go
require lsof
go mod download
build_binaries
log "built guard, proxy, fake-upstream and client-sim into bin/"

if ! grep -q '^\.run/$' "$ROOT/.gitignore" 2>/dev/null; then
  printf ".run/\n" >>"$ROOT/.gitignore"
fi

log "setup done"
