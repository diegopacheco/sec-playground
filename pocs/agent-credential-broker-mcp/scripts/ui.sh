#!/usr/bin/env bash
set -euo pipefail
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/common.sh"

port="$(service_port broker)"
[ -n "$port" ] || fail "broker is not declared in scripts/ports.env"
port_up "$port" || fail "broker is not running on $port, run ./scripts/start-all.sh first"

log "opening $(service_url broker)"
open_url "$(ui_link)"
