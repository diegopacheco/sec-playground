#!/usr/bin/env bash
set -euo pipefail
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/common.sh"

port="$(service_port guard)"
[ -n "$port" ] || fail "guard is not declared in scripts/ports.env"
port_up "$port" || fail "guard is not running on $port, run ./scripts/start-all.sh first"

log "opening $(service_url guard)"
open_url "$(review_url)"
