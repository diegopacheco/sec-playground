# auth0-rust

Rust SDK for Auth0 Authentication API and Management API.

## Scope

This crate provides:

- Authentication API request builders and typed response structs.
- Management API endpoint coverage generated from the Java SDK raw clients.
- One method per Management API operation on `ManagementApi` and `RawManagementApi`.
- Flexible JSON model wrappers for generated Management API request and response types.
- Blocking and async execution.
- Static token and client-credentials Management API setup.
- Request options for headers, query parameters, and per-request timeout.
- Proxy, logging, retry, and custom HTTP client configuration.
- Multipart request bodies for import and file endpoints.
- Offset paging iterator helpers.
- Auth0 HTTP status error structs.
- RSA client assertion signing with RS256 and RS384.

## Authentication API

Create a client with a domain and client ID:

```rust
use auth0_rust::AuthApi;

let auth = AuthApi::builder("tenant.auth0.com", "client_id")
    .client_secret("client_secret")
    .build()?;
let url = auth
    .authorize_url("https://app.local/callback")
    .scope("openid profile email")
    .build();
```

Build token, passwordless, signup, reset-password, PAR, JAR, back-channel, token exchange, and MFA requests with `AuthApi` methods. Execute with `auth.execute(&request)` or use the request directly.

## Management API

Create a Management API client with a static token:

```rust
use auth0_rust::ManagementApi;

let mgmt = ManagementApi::builder()
    .domain("tenant.auth0.com")
    .token("api_token")
    .build()?;
let request = mgmt.users_get()?.path_param("id", "auth0|user");
```

Every generated route is available as a flattened method. Dotted Java-style groups are mapped to snake case, so `organizations.members.roles.assign` becomes `organizations_members_roles_assign`.

## Raw API

Use `mgmt.raw()` for raw route access while keeping the same generated methods:

```rust
let request = mgmt.raw().clients_list()?;
```

## Flexible Models

Generated Management API Java DTO names are available as Rust wrappers around `serde_json::Value`. This keeps every generated type name available without forcing thousands of brittle hand-written fields.

## Request Options

Use `RequestOptions` to add custom headers, extra query parameters, and timeouts:

```rust
use auth0_rust::RequestOptions;
use std::time::Duration;

let request = mgmt
    .clients_get()?
    .path_param("id", "client_id")
    .options(
        RequestOptions::new()
            .query("fields", "name")
            .timeout(Duration::from_secs(10)),
    );
```

## Client Options

`ClientOptions` configures proxy, logging, retry, timeout, and custom reqwest clients. Pass it to `AuthApiBuilder::client_options` or `ManagementApiBuilder::client_options`.

Sensitive headers and bodies are redacted by default. Body logging requires an explicit `include_sensitive_bodies(true)` opt-in. Retries use bounded exponential backoff, honor numeric `Retry-After` values, and retry only idempotent methods unless `retry_non_idempotent(true)` is selected.

## Production Behavior

- Auth0 domains require HTTPS. Plain HTTP is accepted only for loopback test servers.
- RSA client assertions are signed for each execution and each retry.
- Automatic Management API tokens use `expires_in`, refresh early, and refresh once after a `401` response.
- Automatic token acquisition uses async transport in async execution paths.
- SDK clients can be created and dropped safely inside async runtimes.
- MFA OTP, OOB, and recovery-code token exchanges are available.
- `AuthApi::execute_json` and `AuthApi::execute_json_async` decode directly into typed response structs.

## Remaining Limitations

- This is not an official Auth0 SDK and owns its protocol, compatibility, and release risk.
- Management API models are wrappers around `serde_json::Value`, not field-typed Rust structs.
- Management API parameters use `path_param`, `query`, `body`, and `multipart` builders instead of route-specific typed signatures.
- Generated operation names are flattened snake case names instead of nested resource clients.
- The endpoint catalog comes from Java raw client source parsing rather than the Auth0 Fern source.
- The catalog contains 433 routes, but tests do not execute every route or validate every response schema.
- Generated route methods do not document required path parameters, query parameters, or request bodies.
- Endpoint-specific error response bodies are not field-typed.
- Authentication response structs cover common fields but leave many nested fields as flexible JSON.
- Client assertion signing supports PEM RSA keys with RS256 and RS384 only.
- Custom HTTP clients are Reqwest clients rather than transport traits.
- `Retry-After` HTTP dates are not parsed; numeric seconds are supported.
- Token refresh is not single-flight, so concurrent cache misses can request multiple replacement tokens.
- Pagination is blocking, offset-based, not generated with route-specific item types, and infers item arrays from response JSON.
- Multipart parts are in memory and do not include file-path helpers.
- The crate requires Rust 1.94 and edition 2024.
- The crate does not create application sessions, process callbacks, verify state or nonce values, manage cookies, or validate JWT signatures and claims.
- Tests use loopback HTTP servers and local request fixtures rather than a live Auth0 tenant.

## Tests

Run:

```bash
cargo test
cargo clippy --all-targets -- -D warnings
cargo doc --no-deps
```
