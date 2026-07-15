#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"
source ./auth0-env.sh
sbt "IntegrationTest / test"
