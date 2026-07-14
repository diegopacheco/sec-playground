# auth0-kotlin

An idiomatic Kotlin wrapper over the official Auth0 Java SDK 3.10.0.

This artifact does not implement Auth0 protocols. Authentication requests, Management API operations, models, transport, token handling, retries, and errors come from `com.auth0:auth0`.

This wrapper adds:

* Kotlin configuration blocks
* Typed properties for every top-level Management API client
* Coroutine adapters for Java SDK requests
* Direct access to every Java SDK API through `java` and `use`

This artifact is not an official Auth0 SDK. The wrapped Java dependency is the supported SDK.

## Status

This wrapper is usable for server-side Kotlin applications that want Kotlin configuration blocks and coroutine adapters for Authentication API requests. Management API calls retain the synchronous behavior of the Java SDK.

## What Is Missing

* Coroutine adapters for synchronous Management API operations
* HTTP-level tests for cancellation, error mapping, retries, and automatic token renewal
* Compatibility with running Gradle directly on Java 25; use Java 17 for the current Gradle wrapper
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

## Build

```bash
./gradlew clean build
```
