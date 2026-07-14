#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"

PORT="${PORT:-8080}"
BASE="http://localhost:${PORT}"
IDENTITY="X-Auth-User: auth0|user-42"

echo "=== 1. request without a mesh-validated identity (expect 401) ==="
curl -s -o /dev/null -w "http status: %{http_code}\n" "${BASE}/api/v1/books"
echo

echo "=== 2. create 12 books (identity forwarded by Istio) ==="
for i in $(seq 1 12); do
  curl -s -H "${IDENTITY}" -H "Content-Type: application/json" \
    -d "{\"title\":\"Book ${i}\",\"author\":\"Author ${i}\",\"isbn\":\"978-1-0000-${i}\",\"year\":$((2000 + i))}" \
    "${BASE}/api/v1/books" > /dev/null
done
echo "created 12 books"
echo

echo "=== 3. page 1, page_size 5 ==="
curl -s -H "${IDENTITY}" "${BASE}/api/v1/books?page=1&page_size=5"
echo
echo

echo "=== 4. page 3, page_size 5 (last page, 2 items) ==="
curl -s -H "${IDENTITY}" "${BASE}/api/v1/books?page=3&page_size=5"
echo
echo

echo "=== 5. get book by id ==="
curl -s -H "${IDENTITY}" "${BASE}/api/v1/books/1"
echo
echo

echo "=== 6. delete book by id (expect 204) ==="
curl -s -o /dev/null -w "http status: %{http_code}\n" -X DELETE -H "${IDENTITY}" "${BASE}/api/v1/books/1"
