#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"

export KIND_EXPERIMENTAL_PROVIDER=podman
CLUSTER=book-service
CTX=kind-book-service
NS=book-service
IMAGE=localhost/book-service:latest
KCADM=/opt/keycloak/bin/kcadm.sh

kx() { kubectl --context "${CTX}" "$@"; }

provision_keycloak() {
  local pod
  pod=$(kx -n "${NS}" get pod -l app=keycloak -o jsonpath='{.items[0].metadata.name}')
  kx -n "${NS}" exec "${pod}" -- ${KCADM} config credentials --server http://localhost:8080 --realm master --user admin --password admin
  kx -n "${NS}" exec "${pod}" -- ${KCADM} create realms -s realm=books -s enabled=true -s accessTokenLifespan=3600 || true
  kx -n "${NS}" exec "${pod}" -- ${KCADM} create clients -r books -s clientId=book-client -s secret=book-secret -s enabled=true -s protocol=openid-connect -s publicClient=false -s standardFlowEnabled=false -s directAccessGrantsEnabled=false -s serviceAccountsEnabled=true || true
  local cid
  cid=$(kx -n "${NS}" exec "${pod}" -- ${KCADM} get clients -r books -q clientId=book-client --fields id --format csv --noquotes | tr -d '\r')
  kx -n "${NS}" exec "${pod}" -- ${KCADM} create clients/"${cid}"/protocol-mappers/models -r books -s name=book-audience -s protocol=openid-connect -s protocolMapper=oidc-audience-mapper -s 'config."included.custom.audience"=book-service.api' -s 'config."access.token.claim"=true' -s 'config."id.token.claim"=false' || true
}

if ! command -v istioctl >/dev/null 2>&1; then
  export ISTIO_VERSION="${ISTIO_VERSION:-1.30.2}"
  curl -L https://istio.io/downloadIstio | sh -
  export PATH="$PWD/istio-${ISTIO_VERSION}/bin:$PATH"
fi

if ! kind get clusters | grep -q "^${CLUSTER}$"; then
  kind create cluster --name "${CLUSTER}" --wait 120s
fi

istioctl install --set profile=default -y --context "${CTX}"

kx apply -f deploy/namespace.yaml

podman build -t "${IMAGE}" -f Containerfile .
podman save "${IMAGE}" -o book-service.tar
kind load image-archive book-service.tar --name "${CLUSTER}"
rm -f book-service.tar

kx apply -f deploy/keycloak.yaml
kx -n "${NS}" rollout status deployment/keycloak --timeout=240s
provision_keycloak

kx apply -f deploy/deployment.yaml
kx apply -f deploy/requestauthentication.yaml
kx apply -f deploy/authorizationpolicy.yaml
kx apply -f deploy/gateway.yaml

kx -n "${NS}" rollout status deployment/book-service --timeout=120s
kx -n istio-system rollout status deployment/istio-ingressgateway --timeout=120s

echo "cluster ready: run ./test.sh"
