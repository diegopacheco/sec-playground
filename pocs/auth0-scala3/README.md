# auth0-scala-3x

An idiomatic Scala 3 wrapper over the official Auth0 Java SDK 3.10.0.

This artifact does not implement Auth0 protocols. Authentication requests, Management API operations, models, transport, token handling, retries, and errors come from `com.auth0:auth0`.

This wrapper adds:

* Scala configuration functions
* Typed access to every top-level Management API client
* `Future` adapters for Java SDK requests
* Direct access to every Java SDK API through `java` and `use`

This artifact is not an official Auth0 SDK. The wrapped Java dependency is the supported SDK.

## Authentication

```scala
val auth = Auth0.auth("tenant.auth0.com", "client-id", "client-secret")

val url = auth.authorizationUrl(
  "https://app.test/callback",
  _.withAudience("https://api.test").withScope("openid profile email")
)

val userInfo = auth.java.userInfo("access-token").bodyFuture
```

## Management API

```scala
val management = Auth0.managementWithToken("tenant.auth0.com", "access-token")
val user = management.users.get("auth0|123")
val roles = management.users.roles()
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
