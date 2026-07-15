# Macaroon Capability Authorization in Go

A dependency-free Go service for issuing, attenuating, and verifying capability tokens restricted by resource, operation, location, and expiration.

## Run

```bash
./start.sh
```

Open `http://127.0.0.1:8092`.

```bash
./stop.sh
```

## Check

```bash
./test.sh
```

The check mints authority for `/records/*`, attenuates it to `/records/payroll`, proves that payroll access succeeds, and proves that budget access fails.

## What is a macaroon?

A macaroon is a bearer capability with an integrity signature and a list of restrictions called caveats. Possession grants only the authority described by every caveat in the token.

Traditional bearer tokens are commonly minted with a fixed set of claims. Changing their scope usually requires returning to the issuer. A macaroon holder can add restrictions without knowing the issuer's root key. This makes delegation safer: a service can receive broad authority, narrow it for one job, and pass the reduced capability onward.

The important property is monotonic attenuation. New caveats are joined with logical AND. They can reduce authority but cannot expand it.

## Signature chain

The issuer starts with a secret root key and a random identifier:

```text
signature_0 = HMAC-SHA256(root_key, identifier)
```

Each caveat replaces the signature with another HMAC:

```text
signature_1 = HMAC-SHA256(signature_0, caveat_1)
signature_2 = HMAC-SHA256(signature_1, caveat_2)
```

The current signature becomes the key for the next link. A holder can therefore append a caveat using the signature already present in the token. The holder cannot remove or alter an earlier caveat because doing so would require reconstructing a valid chain from the secret root key.

The verifier knows the root key, rebuilds the entire chain, compares signatures in constant time, and evaluates every caveat. Unknown caveat types fail closed.

## Supported caveats

| Type | Meaning | Accepted value |
|---|---|---|
| `resource` | Resource path the request may reach | Exact text, `*`, or a trailing wildcard such as `/records/*` |
| `operation` | Action the request may perform | Exact text or `*` |
| `location` | Request location | Exact text or `*` |
| `expires` | Final valid instant | RFC3339 timestamp |

Multiple caveats of the same type must all pass. Adding `resource=*` after `resource=/records/payroll` does not widen access because both conditions remain in the chain.

## Capability flow

```text
Issuer root key
      |
      v
resource=/records/*
operation=*
location=*
expires=2026-07-14T20:00:00Z
      |
      v
Holder adds resource=/records/payroll
      |
      v
Payroll read: allowed
Budget read: denied
```

## HTTP API

### Mint

```bash
curl -sS -X POST http://127.0.0.1:8092/api/macaroon \
  -H 'Content-Type: application/json' \
  -d '{"resource":"/records/*","operation":"*","location":"*","expires_in_seconds":900}'
```

### Attenuate

```bash
curl -sS -X POST http://127.0.0.1:8092/api/attenuate \
  -H 'Content-Type: application/json' \
  -d '{"token":"TOKEN","caveat":{"type":"operation","value":"read"}}'
```

### Verify

```bash
curl -sS -X POST http://127.0.0.1:8092/api/verify \
  -H 'Content-Type: application/json' \
  -d '{"token":"TOKEN","resource":"/records/payroll","operation":"read","location":"us-west"}'
```

## Security properties

- HMAC-SHA256 protects the identifier and the ordered caveat chain.
- Constant-time signature comparison reduces timing leakage.
- Random 128-bit identifiers prevent predictable token identity.
- Every caveat is evaluated using conjunction.
- Unsupported caveats deny access.
- Tokens have a required expiration caveat.
- HTTP requests are size-limited and reject unknown JSON fields.
- Browser responses include restrictive security headers.

## Scope

This project implements first-party caveats. It does not implement third-party caveats or discharge macaroons. The root key is configurable through `MACAROON_ROOT_KEY`; the built-in value is suitable only for local use.

The serialized token is base64url-encoded JSON so the chain remains easy to inspect. Encoding is not encryption. Anyone holding the token can read its identifier and caveats.
