# auth0-scala-3x

An idiomatic Scala 3 wrapper over the official Auth0 Java SDK 3.10.0.

This artifact does not implement Auth0 protocols. Authentication requests, Management API operations, models, transport, token handling, retries, and errors come from `com.auth0:auth0`.

This wrapper adds:

* Scala configuration types for authentication, authorization URLs, logout, and Management API clients
* Typed access to every top-level Management API client
* `Future` adapters for Java SDK requests
* Direct access to every Java SDK API through `java` and `use`

This artifact is not an official Auth0 SDK. The wrapped Java dependency is the supported SDK.

## Idiomatic Scala 3 Design

The public API follows Scala 3 conventions:

* Immutable `val` references hold the wrapped Java clients, while configuration functions isolate the mutation required by Java builders.
* Default function arguments keep optional configuration concise without nullable sentinels.
* Extension methods add `Future`, body, status, and wrapping operations to Java SDK types without inheritance.
* `Option` represents a possibly absent Java response body, and `Future` represents asynchronous request completion.
* A deliberate `export` clause exposes the selected Management API surface without repetitive forwarding methods or widening the wrapper to every Java member.
* MUnit suites use normal sbt test discovery, structured assertions, and native `Future` test completion.

## Idiomatic Audit Changes

The audit replaced the manually forwarded Management API methods with a Scala 3 `export` facade, replaced the standalone test entry point with a discovered MUnit suite, removed blocking waits from asynchronous tests, and enabled deprecation, feature, unchecked, and unused-code compiler warnings.

## Status

This wrapper is usable for server-side Scala 3 applications that want small Scala adapters over the official Java SDK. Management API calls retain the synchronous behavior of the Java SDK.

## What Is Missing

* `Future` adapters for synchronous Management API operations
* Handwritten return types for the exported Management API methods when stronger binary compatibility than the Java facade provides is required
* HTTP-level tests for error mapping, retries, and automatic token renewal
* Session creation, callback processing, state and nonce verification, cookie handling, and JWT validation
* Scala-native models and Management API method signatures independent of the Java SDK

The wrapper intentionally remains thin. Use `java` or `use` whenever the Scala surface does not expose an operation directly.

## Scala 3 Compared To The Java SDK

Side by side against `com.auth0:auth0:3.10.0`. The Java column is the official API. The Scala column is what this wrapper puts in front of it. Counts come from `javap` on the published jar and from running the suites.

| Capability | Java SDK 3.10.0 | Scala 3 wrapper | Tested |
| --- | --- | --- | --- |
| Client construction | `AuthAPI.newBuilder(domain, clientId)`, with secret, with `ClientAssertionSigner` | `Auth0.auth(domain, clientId)`, with secret, with signer | Unit |
| Authorization URL | `authorizeUrl(uri)` returns a mutable builder | `authorizationUrl(uri, config)` returns the built `String` | Unit and live tenant |
| Logout URL | `logoutUrl(uri, includeClientId)` | `logoutUrl(uri, includeClientId, config)` | Unit |
| Remaining 27 auth operations | 27 typed methods | No Scala signature, reach through `auth.java` or `auth.use` | `exchangeCode` on live tenant |
| Management API clients | 45 accessor methods | 45 through the `export` facade | Unit |
| Management construction | `ManagementApiBuilder` | `ManagementConfig` | Unit |
| Async request | `executeAsync()` returns `CompletableFuture` | `future` and `bodyFuture` return `Future` | Unit |
| Response body | `getBody()` may return `null` | `bodyOption` returns `Option` | Unit |
| Response status | `getStatusCode()` | `statusCode` | Unit |
| Models, transport, retries, token renewal, error mapping | Implemented | Delegated unchanged | Not tested here |
| Sessions, cookies, state and nonce, JWT validation | Not provided | Not provided | Not applicable |

Test totals: 6 unit tests and 1 live integration test against a real tenant.

## Features Missing In The Scala 3 Wrapper

The wrapper is deliberately thin. These gaps are real and worth knowing before adopting it.

- **No Scala overloads for 27 of the 30 `AuthAPI` methods.** `login`, `exchangeCode`, `userInfo`, `renewAuth`, `revokeToken`, `signUp`, `resetPassword`, MFA, passwordless, PAR, and JAR have no Scala signature. Call them through `auth.java` or `auth.use` and handle Java nulls and exceptions yourself.
- **No `Future` adapters for Management API calls.** Every management operation blocks the calling thread. `future` and `bodyFuture` only apply to `Request` values, which the Authentication API returns.
- **No Scala-native models.** Requests and responses stay Java types such as `TokenHolder`, `UserInfo`, and `User`, so `null` and Java collections cross into Scala code.
- **Errors are Java exceptions.** `APIException` and `OAuthTokenException` are thrown rather than returned as `Either` or `Try`.
- **No session handling.** No callback processing, state or nonce verification, cookie management, or JWT signature and claim validation. The Java SDK does not provide these either, so `sample-app/server` implements them by hand.
- **Thin coverage of the wrapped surface.** The 6 unit tests cover the Scala adapters only. Error mapping, retries, and automatic token renewal are delegated to the Java SDK and are not tested here. The live integration test exercises only the authorization URL and the code exchange, so `login`, MFA, passwordless, PAR, and JAR are unverified against a real tenant.

## Authenticate a User

The full authorization code flow with the wrapper. This is what `sample-app/server` runs.

```scala
import auth0scala3.Auth0

val auth = Auth0.auth("tenant.auth0.com", "client-id", "client-secret")

val loginUrl = auth.authorizationUrl(
  "http://localhost:3000/callback",
  config =>
    config.scope("openid profile email")
    config.state("random-state")
)

val tokens = auth.java.exchangeCode("code-from-callback", "http://localhost:3000/callback").execute().getBody
val idToken = tokens.getIdToken
val accessToken = tokens.getAccessToken

val profile = auth.java.userInfo(accessToken).execute().getBody.getValues
```

Send the browser to `loginUrl`, and Auth0 returns an authorization code to the callback. The code exchange runs on the server, so the client secret never reaches the browser.

## Authentication

```scala
val auth = Auth0.auth("tenant.auth0.com", "client-id", "client-secret")

val url = auth.authorizationUrl(
  "https://app.test/callback",
  config =>
    config.audience("https://api.test")
    config.scope("openid profile email")
)

val userInfo = auth.java.userInfo("access-token").bodyFuture
```

## Management API

```scala
val management = Auth0.managementWithToken("tenant.auth0.com", "access-token")
val user = management.users().get("auth0|123")
val roles = management.users().roles()
```

Management configuration also uses Scala types:

```scala
val management = Auth0.managementWithToken("tenant.auth0.com", "access-token", config =>
  config.timeout(10)
  config.maxRetries(2)
  config.header("x-request-source", "scala-service")
)
```

Automatic client credentials use the Java SDK implementation:

```scala
val management = Auth0.managementWithClientCredentials(
  "tenant.auth0.com",
  "client-id",
  "client-secret"
)
```

## Prerequisites

### Toolchain

| Requirement | Version | Notes |
| --- | --- | --- |
| JDK | 17 or newer | The wrapped `com.auth0:auth0:3.10.0` runs on the JVM |
| sbt | 1.11.6 | Pinned by `project/build.properties`, the sbt launcher downloads this version |
| Scala | 3.7.2 | Pinned by `build.sbt`, downloaded by sbt, no manual install |
| Node.js | 20 or newer | Only for `sample-app/webapp`, which builds with Vite 8 and React 19 |

Verify:

```bash
java -version
sbt --version
node --version
```

### Auth0 Account

A free Auth0 account and one tenant are required. Nothing in this project runs against a live tenant without them.

1. Create an account at https://auth0.com and note the tenant domain, shaped `your-tenant.us.auth0.com`.
2. Create an application of type **Regular Web Application**.
3. Copy the **Domain**, **Client ID**, and **Client Secret** from the application settings.
4. Add `http://localhost:3000/login/oauth2/code/okta` to **Allowed Callback URLs**. The sample app server only accepts this path.
5. Add `http://localhost:3000` to **Allowed Logout URLs**.

The client secret is a server-side credential. The authorization code exchange runs on the server so the secret never reaches the browser.

### Environment Script

`auth0-env.sh` holds the tenant credentials and is listed in `.gitignore`, so a fresh clone does not contain it. Create it in this project root before running anything:

```bash
cat > auth0-env.sh <<'EOF'
export AUTH0_DOMAIN=your-tenant.us.auth0.com
export AUTH0_CLIENT_ID=your-client-id
export AUTH0_CLIENT_SECRET=your-client-secret
export AUTH0_REDIRECT_URI=http://localhost:3000/login/oauth2/code/okta
EOF
```

Replace every placeholder with the values from step 3. Keep the file untracked, it carries a real secret.

| Variable | Used by |
| --- | --- |
| `AUTH0_DOMAIN` | Client construction, authorization and logout URLs, Management API |
| `AUTH0_CLIENT_ID` | Client construction, authorization URL |
| `AUTH0_CLIENT_SECRET` | Authorization code exchange, client credentials |
| `AUTH0_REDIRECT_URI` | Authorization URL and code exchange, must match the callback configured in Auth0 |

Load it into the current shell:

```bash
source ./auth0-env.sh
```

`test-integration.sh` and `sample-app/start.sh` source this file themselves.

## Build

```bash
sbt clean test package
```

Unit tests need no tenant. The live integration test needs the environment script:

```bash
./test-integration.sh
```
