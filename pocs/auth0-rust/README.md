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

## What Is Still Not 100%

This crate is broad and builds, but it is not a byte-for-byte or class-for-class Rust port of the Java SDK. Known gaps:

- Management API models are nominal wrappers around `serde_json::Value`, not fully field-typed Rust structs for every property in every Java DTO.
- Management API methods are one per route, but their parameters are provided through `path_param`, `query`, `body`, and `multipart` builders instead of strongly typed Rust method parameters for every route.
- Generated Management API operation names are flattened snake case names, not nested client structs for every Java package path.
- Java overloads are represented by builder methods and suffixed generated route methods where routes collide, not by overloaded functions because Rust does not support overloads.
- Async support exists for request execution, but automatic Management API client-credentials token acquisition still uses the shared synchronous token cache path.
- RSA client assertion signing supports RS256 and RS384 from PEM keys, but it does not expose every Java key object construction path.
- Proxy support maps to reqwest proxy URLs and basic auth, not Java `Proxy` objects.
- Logging support writes through simple stderr logging, not OkHttp interceptor-compatible logging.
- Custom HTTP client injection accepts reqwest blocking and async clients, not arbitrary user-defined transport traits.
- Retry support retries selected HTTP status codes and can respect `Retry-After`, but it does not clone Java interceptor internals exactly.
- Rate-limit parity is limited to retry behavior and error structs; it does not expose Java token quota helper classes as stateful objects.
- Pagination supports offset paging from common Auth0 response shapes, but it is not generated per endpoint with route-specific item types.
- Multipart support accepts in-memory parts; it does not yet provide file-path convenience helpers.
- Auth API typed responses cover the common Java auth response classes, but not every nested field is strongly typed.
- Auth API request builders cover the Java Authentication API operations inspected from `AuthAPI`, but method naming is Rust-style and not Java-style.
- Request options support headers, query parameters, and timeout; they do not model every Java request option object one-to-one.
- HTTP status error structs exist for the Java Management API status classes, but endpoint-specific error body schemas are not field-typed.
- Tests include route coverage, request construction, Java auth fixture decoding, multipart, request options, and error mapping, but they do not replicate the full Java wire test suite.
- The endpoint catalog was generated from Java raw client source parsing, not from the original Fern/OpenAPI source.
- The crate does not include publishing metadata parity with the Java project.

## Tests

Run:

```bash
cargo test
```
