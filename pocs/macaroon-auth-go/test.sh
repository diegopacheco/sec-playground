#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"
go test ./...
./start.sh
trap './stop.sh >/dev/null 2>&1 || true' EXIT
root_response="$(curl -fsS -X POST http://127.0.0.1:8092/api/macaroon -H 'Content-Type: application/json' -d '{"resource":"/records/*","operation":"*","location":"*","expires_in_seconds":300}')"
root_token="$(printf '%s' "$root_response" | sed -n 's/.*"token":"\([^"]*\)".*/\1/p')"
attenuated_response="$(curl -fsS -X POST http://127.0.0.1:8092/api/attenuate -H 'Content-Type: application/json' -d "{\"token\":\"$root_token\",\"caveat\":{\"type\":\"resource\",\"value\":\"/records/payroll\"}}")"
attenuated_token="$(printf '%s' "$attenuated_response" | sed -n 's/.*"token":"\([^"]*\)".*/\1/p')"
read_response="$(curl -fsS -X POST http://127.0.0.1:8092/api/verify -H 'Content-Type: application/json' -d "{\"token\":\"$attenuated_token\",\"resource\":\"/records/payroll\",\"operation\":\"read\",\"location\":\"us-west\"}")"
other_response="$(curl -fsS -X POST http://127.0.0.1:8092/api/verify -H 'Content-Type: application/json' -d "{\"token\":\"$attenuated_token\",\"resource\":\"/records/budget\",\"operation\":\"read\",\"location\":\"us-west\"}")"
[[ "$read_response" == *'"allowed":true'* ]]
[[ "$other_response" == *'"allowed":false'* ]]
echo "Root capability: /records/*"
echo "Added restriction: /records/payroll"
echo "Payroll request: $read_response"
echo "Budget request: $other_response"
echo "Capability attenuation passed"
./stop.sh
