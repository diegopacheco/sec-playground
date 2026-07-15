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

## Build

```bash
sbt clean test package
```
