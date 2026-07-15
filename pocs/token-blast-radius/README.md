# Token Blast Radius

A modular Go and React application that combines signed token claims with authorization relationships to calculate effective credential reach. The interface explains active paths, blocked paths, trust boundaries, sensitive targets, and why each relationship is allowed.

## Capabilities

- Verifies compact HMAC-SHA256 signed tokens before analysis.
- Preserves the human-to-agent chain through `act` and `sub` claims.
- Applies audience, scopes, and maximum delegation depth.
- Traverses relationship graphs and finds shortest active paths.
- Separates reachable edges from blocked authorization paths.
- Counts reachable and critical resources.
- Produces a deterministic risk score and remediation verdict.
- Provides three distinct identity-security scenarios.
- Lets the user alter a token signature and observe verification failure.
- Uses a custom SVG graph without a visualization library.

## Interface

The React interface contains:

- Interactive credential scenario cards
- Animated delegation graph
- Reachability and risk metrics
- Signed claim inspector
- Token integrity controls
- Clickable node explanations
- Shortest active access paths
- Educational identity and authorization guidance
- Responsive layouts and reduced-motion support

## Structure

```text
cmd/server                  process entry point
internal/domain             shared domain types
internal/graph              traversal and risk engine
internal/httpapi            API and frontend delivery
internal/scenario           identity relationship catalog
internal/token              token signing and validation
frontend/src/components     modular React interface
frontend/src                API, types, application, and visual system
```

## Technology

| Layer | Technology |
|---|---|
| Backend | Go 1.26 |
| Frontend | React 19.2 |
| Language | TypeScript 7 |
| Build | Vite 8.1 |
| Package manager | Bun |
| Graph rendering | Native SVG |

## Run

```bash
./start.sh
```

Open `http://localhost:8082`.

```bash
./stop.sh
```

## Verify

```bash
./test.sh
```

The script runs Go tests, TypeScript checks, the production frontend build, contained and critical graph analysis, signature rejection, and HTTP delivery verification.

## API

```text
GET  /health
GET  /api/scenarios
POST /api/analyze
```

Analysis request:

```json
{
  "scenario_id": "01-contained",
  "token": "SIGNED_TOKEN"
}
```

## Configuration

| Variable | Default | Purpose |
|---|---|---|
| `PORT` | `8082` | HTTP port |
| `BLAST_RADIUS_SECRET` | Local development value in `start.sh` | HMAC signing secret |
| `FRONTEND_PATH` | `frontend/dist` | Compiled frontend directory |
