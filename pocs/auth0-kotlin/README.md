# auth0-kotlin

An idiomatic Kotlin wrapper over the official Auth0 Java SDK 3.10.0.

This artifact does not implement Auth0 protocols. Authentication requests, Management API operations, models, transport, token handling, retries, and errors come from `com.auth0:auth0`.

This wrapper adds:

* Kotlin configuration blocks
* Typed properties for every top-level Management API client
* Coroutine adapters for Java SDK requests
* Direct access to every Java SDK API through `java` and `use`

This artifact is not an official Auth0 SDK. The wrapped Java dependency is the supported SDK.

## Idiomatic Kotlin Design

The public API follows Kotlin library conventions:

* `explicitApi()` makes every public visibility and return type deliberate and prevents accidental API changes caused by inferred Java types.
* Receiver lambdas provide type-safe configuration blocks, default lambda arguments keep the common path concise, and named arguments remain available at call sites.
* `@Auth0Dsl` marks configuration receivers as a Kotlin DSL and prevents accidental receiver mixing when configuration blocks are nested.
* Wrapper state is exposed through immutable `val` properties, and Management API clients have explicit Kotlin property types.
* Extension functions adapt existing Java clients and requests without inheritance or changes to the official SDK.
* `await` uses cancellable suspension, propagates Java failures, and cancels the underlying future when the coroutine is cancelled.

## Idiomatic Audit Changes

The audit enabled explicit API mode, added explicit public signatures across the wrapper, added the DSL marker to authentication and management configuration receivers, and normalized multiline Kotlin declarations with trailing commas. The protocol and transport behavior remain delegated to the official Java SDK.

## Status

This wrapper is usable for server-side Kotlin applications that want Kotlin configuration blocks and coroutine adapters for Authentication API requests. Management API calls retain the synchronous behavior of the Java SDK.

## What Is Missing

* Coroutine adapters for synchronous Management API operations
* HTTP-level tests for cancellation, error mapping, retries, and automatic token renewal
* Session creation, callback processing, state and nonce verification, cookie handling, and JWT validation
* Kotlin-native models and Management API method signatures independent of the Java SDK

The wrapper intentionally remains thin. Use `java` or `use` whenever the Kotlin surface does not expose an operation directly.

## Authentication

```kotlin
val auth = Auth0.auth("tenant.auth0.com", "client-id", "client-secret")

val url = auth.authorizationUrl("https://app.test/callback") {
    withAudience("https://api.test")
    withScope("openid profile email")
}

val userInfo = auth.java.userInfo("access-token").awaitBody()
```

## Kotlin Compared To The Java SDK

Side by side against `com.auth0:auth0:3.10.0`. The Java column is the official API. The Kotlin column is what this wrapper puts in front of it. Counts come from `javap` on the published jar and from running the suites.

| Capability | Java SDK 3.10.0 | Kotlin wrapper | Tested |
| --- | --- | --- | --- |
| Client construction | `AuthAPI.newBuilder(domain, clientId)`, with secret, with `ClientAssertionSigner` | `Auth0.auth(domain, clientId)`, with secret, with signer | Unit |
| Authorization URL | `authorizeUrl(uri)` returns a mutable builder | `authorizationUrl(uri) { }` receiver lambda returns the built `String` | Unit and live tenant |
| Logout URL | `logoutUrl(uri, includeClientId)` | `logoutUrl(uri, includeClientId) { }` | Unit |
| Remaining 27 auth operations | 27 typed methods | No Kotlin signature, reach through `auth.java` or `auth.use` | `exchangeCode` on live tenant |
| Management API clients | 45 accessor methods | 45 typed `val` properties | Unit |
| Management construction | `ManagementApiBuilder` | `ManagementConfig` DSL | Unit |
| Async request | `executeAsync()` returns `CompletableFuture` | `await` and `awaitBody` are cancellable suspend functions | Unit |
| Response body | `getBody()` may return `null` | `executeBody()` returns a nullable body | Unit |
| Nullability | Platform types leak into callers | `explicitApi()` makes every public type deliberate | Compile time |
| Models, transport, retries, token renewal, error mapping | Implemented | Delegated unchanged | Not tested here |
| Sessions, cookies, state and nonce, JWT validation | Not provided | Not provided | Not applicable |

Test totals: 6 unit tests and 1 live integration test against a real tenant.

## Features Missing In The Kotlin Wrapper

The wrapper is deliberately thin. These gaps are real and worth knowing before adopting it.

- **No Kotlin overloads for 27 of the 30 `AuthAPI` methods.** `login`, `exchangeCode`, `userInfo`, `renewAuth`, `revokeToken`, `signUp`, `resetPassword`, MFA, passwordless, PAR, and JAR have no Kotlin signature. Call them through `auth.java` or `auth.use` and handle Java nulls and exceptions yourself.
- **No coroutine adapters for Management API calls.** Every management operation blocks the calling thread. `await` and `awaitBody` only apply to `Request` values, which the Authentication API returns.
- **No Kotlin-native models.** Requests and responses stay Java types such as `TokenHolder`, `UserInfo`, and `User`, so platform types and Java collections cross into Kotlin code despite `explicitApi()`.
- **Errors are Java exceptions.** `APIException` and `OAuthTokenException` are thrown rather than returned as `Result`.
- **No session handling.** No callback processing, state or nonce verification, cookie management, or JWT signature and claim validation. The Java SDK does not provide these either, so `sample-app/server` implements them by hand.
- **Thin coverage of the wrapped surface.** The 6 unit tests cover the Kotlin adapters only. Error mapping, retries, automatic token renewal, and coroutine cancellation against a real server are not tested here. The live integration test exercises only the authorization URL and the code exchange, so `login`, MFA, passwordless, PAR, and JAR are unverified against a real tenant.

## Authenticate a User

The full authorization code flow with the wrapper. This is what `sample-app/server` runs.

```kotlin
import com.auth0.kotlin.Auth0

val auth = Auth0.auth("tenant.auth0.com", "client-id", "client-secret")

val loginUrl = auth.authorizationUrl("http://localhost:3000/callback") {
    withScope("openid profile email")
    withState("random-state")
}

val tokens = auth.java.exchangeCode("code-from-callback", "http://localhost:3000/callback").execute().body
val idToken = tokens.idToken
val accessToken = tokens.accessToken

val profile = auth.java.userInfo(accessToken).execute().body.values
```

Send the browser to `loginUrl`, and Auth0 returns an authorization code to the callback. The code exchange runs on the server, so the client secret never reaches the browser.

## Management API

```kotlin
val management = Auth0.managementWithToken("tenant.auth0.com", "access-token")
val user = management.users.get("auth0|123")
val roles = management.users.roles()
```

Automatic client credentials use the Java SDK implementation:

```kotlin
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
| JDK | 17 | Pinned by `jvmToolchain(17)` and the Java toolchain in `build.gradle.kts` |
| Gradle | 9.0.0 | Supplied by the wrapper, use `./gradlew`, no manual install |
| Kotlin | 2.2.21 | Pinned by the `kotlin("jvm")` plugin, resolved by Gradle |
| Node.js | 20 or newer | Only for `sample-app/webapp`, which builds with Vite 8 and React 19 |

Verify:

```bash
java -version
./gradlew --version
node --version
```

Gradle provisions the JDK 17 toolchain when a matching JDK is discoverable. If it is not, install JDK 17 and let Gradle detect it.

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
./gradlew clean build
```

Unit tests need no tenant. The live integration test needs the environment script:

```bash
./test-integration.sh
```
