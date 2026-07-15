# Agent Intent-Bound Tokens

A modular Go service that issues signed, one-time tokens for one declared agent action. A valid signature is insufficient by itself: the audience, action, resource, amount, deadline, and token-use state must all match.

## Security properties

- HMAC-SHA256 signatures protect claims from modification.
- Audience binding prevents use against a different API.
- Exact action and resource binding stop scope expansion.
- Integer monetary limits avoid floating-point authorization errors.
- Activation and expiration timestamps bound the authorization window.
- Successful use consumes the token and blocks replay.
- Every decision produces an in-memory audit record.
- The signing secret must contain at least 32 characters.

## Structure

```text
cmd/server                 process entry point
internal/httpapi           HTTP transport
internal/intent            domain types
internal/token             signing and parsing
internal/verification      constraint evaluation and audit state
```

## Run

```bash
./start.sh
```

The service listens on `http://localhost:8081`.

The root URL returns the available HTTP endpoints as JSON.

```bash
./stop.sh
```

## Interactive UI

Start the API, then run the React interface with Bun:

```bash
cd ui
bun install
bun run dev
```

Open `http://localhost:5173`. The interface issues intent-bound tokens, checks scope escalation and amount overflow, confirms one valid action, rejects replay, and renders the audit ledger.

Build the interface with TypeScript 7:

```bash
cd ui
bun run build
```

## Verify

```bash
./test.sh
```

The script runs all Go tests and then validates action mismatch, amount overflow, successful execution, replay rejection, and audit recording through the HTTP API.

## Issue a token

```bash
curl -s -X POST http://localhost:8081/api/tokens \
  -H "Content-Type: application/json" \
  -d '{"subject":"agent:buyer","audience":"payments-api","action":"transfer","resource":"account:operations","max_amount_cents":50000,"valid_for_seconds":60}'
```

## Verify an action

```bash
curl -s -X POST http://localhost:8081/api/verify \
  -H "Content-Type: application/json" \
  -d '{"token":"SIGNED_TOKEN","audience":"payments-api","action":"transfer","resource":"account:operations","amount_cents":49999}'
```

## Configuration

| Variable | Default | Purpose |
|---|---|---|
| `PORT` | `8081` | HTTP port |
| `INTENT_TOKEN_SECRET` | Local development value in `start.sh` | HMAC signing secret |
| `INTENT_TOKEN_ISSUER` | `intent.local` | Signed issuer claim |
