#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"

CTX=kind-book-service
HOST="books.example.com"
KC=http://localhost:8081
GW=http://localhost:18080

kubectl --context "${CTX}" -n book-service port-forward svc/keycloak 8081:8080 >/dev/null 2>&1 &
KC_PID=$!
kubectl --context "${CTX}" -n istio-system port-forward svc/istio-ingressgateway 18080:80 >/dev/null 2>&1 &
GW_PID=$!
trap 'kill ${KC_PID} ${GW_PID} 2>/dev/null || true' EXIT

for i in $(seq 1 60); do
  if curl -sf "${KC}/realms/books" >/dev/null 2>&1; then break; fi
  sleep 1
done
for i in $(seq 1 60); do
  if curl -s -o /dev/null -H "Host: ${HOST}" "${GW}/"; then break; fi
  sleep 1
done

echo "=== 1. fetch a real token from Keycloak (client_credentials) ==="
TOKEN=$(curl -s -X POST "${KC}/realms/books/protocol/openid-connect/token" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=client_credentials" \
  -d "client_id=book-client" \
  -d "client_secret=book-secret" | sed -n 's/.*"access_token":"\([^"]*\)".*/\1/p')
if [ -z "${TOKEN}" ]; then echo "failed to obtain token"; exit 1; fi
echo "token acquired (${#TOKEN} chars)"
echo

echo "=== 2. call through Istio WITHOUT a token (expect 403, rejected by the mesh) ==="
curl -s -o /dev/null -w "http status: %{http_code}\n" -H "Host: ${HOST}" "${GW}/api/v1/books"
echo

echo "=== 3. call through Istio with an INVALID token (expect 401) ==="
curl -s -o /dev/null -w "http status: %{http_code}\n" -H "Host: ${HOST}" \
  -H "Authorization: Bearer not.a.real.token" "${GW}/api/v1/books"
echo

echo "=== 4. create 12 books with the valid token ==="
for i in $(seq 1 12); do
  curl -s -H "Host: ${HOST}" -H "Authorization: Bearer ${TOKEN}" -H "Content-Type: application/json" \
    -d "{\"title\":\"Book ${i}\",\"author\":\"Author ${i}\",\"isbn\":\"978-1-0000-${i}\",\"year\":$((2000 + i))}" \
    "${GW}/api/v1/books" >/dev/null
done
echo "created 12 books"
echo

echo "=== 5. page 1, page_size 5 ==="
curl -s -H "Host: ${HOST}" -H "Authorization: Bearer ${TOKEN}" "${GW}/api/v1/books?page=1&page_size=5"
echo
echo

echo "=== 6. page 3, page_size 5 (last page, 2 items) ==="
curl -s -H "Host: ${HOST}" -H "Authorization: Bearer ${TOKEN}" "${GW}/api/v1/books?page=3&page_size=5"
echo
echo

echo "=== 7. get book by id ==="
curl -s -H "Host: ${HOST}" -H "Authorization: Bearer ${TOKEN}" "${GW}/api/v1/books/1"
echo
echo

echo "=== 8. delete book by id (expect 204) ==="
curl -s -o /dev/null -w "http status: %{http_code}\n" -X DELETE \
  -H "Host: ${HOST}" -H "Authorization: Bearer ${TOKEN}" "${GW}/api/v1/books/1"
