# auth0-kotlin-2x

Kotlin SDK layer for Auth0 Java SDK 3.x.

Generated management clients and generated request and response models remain Java-backed for binary compatibility. The Kotlin layer keeps the full Java SDK reachable and adds Kotlin entry points, endpoint wrappers, coroutine helpers, Jackson Kotlin support, null-safe response helpers, and DSL blocks.

## Feature parity

This SDK has full feature access to the Auth0 Java SDK through the Java SDK backend. It does not have 100% Kotlin-native feature parity.

The remaining gaps are:

- Endpoint methods are not generated as named Kotlin functions for every operation. Every Java client method is reachable through `call("methodName", ...)`, `request("methodName", ...)`, or the original Java client.

- Nested resources are not all generated as named Kotlin functions. Top-level management resources have Kotlin functions, while nested generated resources are reachable through `resource("resourceName")`.

- Auth operations are not all generated as named Kotlin functions. They are reachable through `request("methodName", ...)`, `call("methodName", ...)`, or the original `AuthAPI`.

- Kotlin nullability is not hand-modeled for every generated Java model property. The Kotlin layer provides non-null and nullable response helpers, but generated model fields remain Java platform types.

- Java tests are not ported line by line to Kotlin. The included Java SDK test suite runs as part of `check`, and the Kotlin layer has parity tests around its own API.

## Kotlin entry points

```kotlin
val auth = Auth0.kotlinAuth("tenant.auth0.com", "client-id", "secret")

val loginUrl = auth.authorizeUrl("https://app.test/callback") {
    withScope("openid profile email")
    withAudience("https://api.test")
}

val management = Auth0.kotlinManagementWithToken("tenant.auth0.com", "token")
val users = management.users()
val roles = users.resource("roles")
```

## Endpoint calls

```kotlin
val user = management.users().call<Any>("get", "auth0|123")
val job = management.jobs().call<Any>("get", "job-id")
```

## Coroutines

```kotlin
val body = auth.request<Any>("userInfo", "access-token").awaitBody()
val response = auth.request<Any>("userInfo", "access-token").awaitResponse()
```

## JSON

```kotlin
val json = Auth0Json.encode(value)
val value = Auth0Json.decode<MyType>(json)
```

## Build

```bash
JAVA_HOME=/Users/diegopacheco/.sdkman/candidates/java/17.0.16-amzn GRADLE_USER_HOME=/private/tmp/auth0-kotlin-gradle ./gradlew --no-daemon clean build
```
