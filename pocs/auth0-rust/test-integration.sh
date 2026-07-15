#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"
source ./auth0-env.sh
cargo test --test auth0_integration -- --ignored --test-threads=1
