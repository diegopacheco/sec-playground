# auth0-scala-3x

An idiomatic Scala 3 wrapper over the official Auth0 Java SDK 3.10.0.

This artifact does not implement Auth0 protocols. Authentication requests, Management API operations, models, transport, token handling, retries, and errors come from `com.auth0:auth0`.

This wrapper adds:

* Scala configuration types for authentication, authorization URLs, logout, and Management API clients
* Typed access to every top-level Management API client
* `Future` adapters for Java SDK requests
* Direct access to every Java SDK API through `java` and `use`

This artifact is not an official Auth0 SDK. The wrapped Java dependency is the supported SDK.

## Status

This wrapper is usable for server-side Scala 3 applications that want small Scala adapters over the official Java SDK. Management API calls retain the synchronous behavior of the Java SDK.

## What Is Missing

* `Future` adapters for synchronous Management API operations
* Standard test discovery and structured test reporting; the current build runs a standalone test entry point
* Explicit return types on the public wrapper surface for stronger API compatibility
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
val user = management.users.get("auth0|123")
val roles = management.users.roles()
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
