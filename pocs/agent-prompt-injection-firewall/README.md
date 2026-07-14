# Agent Prompt Injection Firewall

A deterministic Go service that treats retrieved instructions as untrusted data and places a reference monitor between an agent and its tools. It provides read and HTTP tools while enforcing capabilities before any operation runs.

The protections align with the attack surfaces described by [OWASP Agentic AI Threats and Mitigations](https://genai.owasp.org/resource/agentic-ai-threats-and-mitigations/).

## Security properties

- Retrieved files are always labeled untrusted.
- Read access is confined to `fixtures/data`.
- HTTP requests are confined to an exact host allowlist.
- Redirects are blocked to prevent allowlist bypass.
- HTTP methods are restricted to a small capability set.
- State-changing HTTP calls require an explicit approval token.
- Canary values are blocked before tool execution.
- Canary values are redacted from decisions and audit records.
- Every candidate action receives an auditable allow or deny decision.

## Request flow

```text
retrieved document
       |
       v
untrusted action candidates
       |
       v
capability and data firewall
       |
       +---- deny and audit
       |
       v
read or HTTP tool
       |
       v
redacted result and audit record
```

The action lines in `fixtures/instructions` represent tool candidates produced after an agent interprets retrieved content. They are never trusted simply because they came from an agent.

## Run

```bash
./start.sh
```

The service listens on `http://localhost:8091`.

Run a trusted retrieval:

```bash
curl -s -X POST http://localhost:8091/run -H "Content-Type: application/json" -d '{"document":"safe.txt"}'
```

Run a poisoned retrieval:

```bash
curl -s -X POST http://localhost:8091/run -H "Content-Type: application/json" -d '{"document":"poisoned.txt"}'
```

Approve a state-changing request:

```bash
curl -s -X POST http://localhost:8091/run -H "Content-Type: application/json" -d '{"document":"approval.txt","approval":"approve-once"}'
```

Inspect audit records:

```bash
curl -s http://localhost:8091/audit
```

## Stop

```bash
./stop.sh
```

## Test

```bash
./test.sh
```

The test proves:

- a valid read and allowlisted HTTP request run
- path traversal is denied
- a foreign HTTP destination is denied
- sensitive canary exfiltration is denied and redacted
- an HTTP write is denied without approval
- the same write runs with explicit approval
- all security decisions are present in the audit trail
