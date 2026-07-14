#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"

export KIND_EXPERIMENTAL_PROVIDER=podman
CLUSTER=book-service

kind delete cluster --name "${CLUSTER}"
rm -f book-service.tar
echo "cluster destroyed"
