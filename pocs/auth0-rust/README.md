# auth0-rust

Rust SDK for Auth0 Authentication API and Management API.

## Idiomatic Rust Design

The audit found the crate idiomatic without requiring source changes:

- Builders own configurable construction and return `Result` when domains, clients, or requests can be invalid.
- Read-only inputs use `&str`, `AsRef`, or borrowed requests where ownership is unnecessary, while stored values use owned `String` data.
- Enums model HTTP methods, bodies, signing algorithms, passwordless modes, and API error categories with exhaustive matching.
- `Auth0Error` composes lower-level failures through `From`, implements the standard error traits, and keeps recoverable failure in return types.
- Iterator implementations provide lazy pagination, and collection transformations avoid manual indexing.
- Client, token, lock, and transport lifetimes rely on ownership and scope-based cleanup.
- Generated JSON models derive standard traits, and configurable clients use `Default` where a conventional baseline exists.

## Idiomatic Audit Changes

No Rust API rewrite was needed. The code already passed `cargo fmt --check` and `cargo clippy --all-targets --all-features -- -D warnings`; this README now records the ownership, error, enum, builder, iterator, and cleanup conventions that make the crate idiomatic.

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

## Rust Compared To The Java SDK

Side by side against `com.auth0:auth0:3.10.0`. Unlike the Scala 3 and Kotlin wrappers, this crate is an independent implementation that never calls the Java SDK, so every row is real Rust code rather than a forwarded call. Counts come from `javap` on the published jar, from the generated sources, and from running the suites.

| Capability | Java SDK 3.10.0 | Rust crate | Tested |
| --- | --- | --- | --- |
| Client construction | `AuthAPI.newBuilder(...)` | `AuthApi::builder(...)` returns `Result` instead of throwing | Unit |
| Authorization and logout URLs | `authorizeUrl`, `logoutUrl`, PAR and JAR variants | Same, as owned builders | Unit and live tenant |
| Remaining 27 auth operations | 27 typed methods | Full parity, one method each | 14 unit tests, `exchange_code` on live tenant |
| Management API clients | 45 nested resource clients | 866 flat methods over 433 routes, on `ManagementApi` and `RawManagementApi` | Catalog, routing, and request-render unit tests |
| Management construction | `ManagementApiBuilder` | `ManagementApiBuilder` returning `Result` | Unit |
| Automatic token renewal on 401 | Implemented | Implemented, blocking and async | 2 HTTP tests against a loopback server |
| Async execution | `executeAsync()` returns `CompletableFuture` | Native async alongside blocking | HTTP tests |
| Paging | Page iterators | `SyncPagingIterable`, offset based, blocking only | Unit |
| Multipart bodies | Implemented | Implemented, in memory only | Unit |
| Errors | Thrown `APIException` | Returned `Result` with `Auth0Error` and per-status structs | Unit |
| Client assertion signing | RS256 and RS384, multiple key sources | PEM RSA keys, RS256 and RS384 only | Unit |
| Models | Field-typed Java POJOs | `serde_json::Value` wrappers, not field typed | Partly |
| Beyond the Java SDK | Not present | `jwks`, `login_realm`, `request_token_for_organization`, `sign_up_with_username`, `sign_up_with_phone_number` | `jwks` only |
| Sessions, cookies, state and nonce, JWT validation | Not provided | Not provided | Not applicable |

Test totals: 24 unit tests (7 lib, 14 SDK, 3 HTTP) and 1 live integration test against a real tenant.

## Features Missing Compared To The Java SDK

The crate reaches or exceeds the Java SDK on operation coverage, so the gaps are about typing and depth rather than breadth. The full engineering list is under Remaining Limitations below; these are the ones that change how you write code.

- **Management models are not field typed.** Java returns POJOs such as `User` and `Role`. This crate returns `serde_json::Value` wrappers, so field access is dynamic and typos surface at runtime instead of compile time.
- **Management parameters are untyped builders.** Java gives route-specific signatures. This crate uses `path_param`, `query`, `body`, and `multipart`, so required parameters are neither documented nor enforced.
- **Flat operation names instead of nested clients.** Java exposes `mgmt.users().get(id)`. This crate exposes `mgmt.users_get()`, generated as flattened snake case.
- **The Management API is untested against a live tenant.** All 866 methods are covered only by catalog and request-render tests. The live integration test exercises the Authentication API alone, because the client in `auth0-env.sh` has no Management API grant.
- **Client assertion signing is narrower than Java.** PEM RSA keys with RS256 and RS384 only.
- **Token refresh is not single flight.** Concurrent cache misses can request multiple replacement tokens.
- **No session handling.** No callback processing, state or nonce verification, cookie management, or JWT signature and claim validation. The Java SDK does not provide these either, so `sample-app/server` implements them by hand.

## Authenticate a User

The full authorization code flow. This is what `sample-app/server` runs.

```rust
use auth0_rust::AuthApi;

let auth = AuthApi::builder("tenant.auth0.com", "client_id")
    .client_secret("client_secret")
    .build()?;

let login_url = auth
    .authorize_url("http://localhost:3000/callback")
    .scope("openid profile email")
    .state("random-state")
    .build();

let tokens = auth
    .exchange_code("code-from-callback", "http://localhost:3000/callback")
    .execute()?;

let id_token = tokens.body.as_ref().and_then(|b| b.get("id_token"));
let access_token = tokens.body.as_ref().and_then(|b| b.get("access_token"));

let profile = auth.user_info("access-token").execute()?;
```

Send the browser to `login_url`, and Auth0 returns an authorization code to the callback. The code exchange runs on the server, so the client secret never reaches the browser.

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
- Unit tests use loopback HTTP servers and local request fixtures. Only `tests/auth0_integration.rs` reaches a live Auth0 tenant, and it covers the Authentication API alone.

## Tests

Run:

```bash
cargo test
cargo clippy --all-targets -- -D warnings
cargo doc --no-deps
```
