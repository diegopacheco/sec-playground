# auth0-scala-3x

Scala 3 SDK for Auth0 built on the JVM with an idiomatic Scala Authentication API and a compact Management API resource layer.

## Features

* Authentication API request builders for authorization URLs, logout URLs, token flows, signup, password reset, passwordless, MFA, PAR, JAR, CIBA, and authenticator management.
* Management API resources with typed nested path helpers for users, organizations, guardian, jobs, actions, branding, connections, roles, prompts, tenants, and verifiable credentials.
* Raw response wrappers with decoded body, status, headers, text, and bytes.
* Endpoint-specific errors for common Auth0 HTTP failure codes.
* Async request execution through `Future`.
* Request options for per-call headers, query parameters, and timeout.
* Retry, logging, and proxy-capable Java HTTP transport.
* Auth0 custom domain header whitelisting for supported Management API paths.
* RSA client assertion signing for private key JWT authentication.
* Multipart upload, byte download, and SSE parsing helpers.

## Compatibility

The Auth0 Java SDK contains a large generated Management API layer with 546 dedicated client classes and 2,369 generated request and response model files.

This project does not rewrite that generated layer into Scala classes. Keeping the generated Java Management clients and models available through normal JVM interop preserves binary compatibility and avoids creating a second generated API surface that can drift from Auth0 Java.

Scala can use those Java clients and model classes directly when strict typed parity is required. This SDK focuses Scala code on the handwritten surface: immutable request builders, Scala 3 data types, a small transport abstraction, and generic Management API resource access.

The same compatibility rule applies to generated request and response models. This SDK uses `Json` for generic Management API payloads, while Java model classes remain the typed source of truth for users who need the generated schema classes.

## Not 100% Java Parity

This SDK is not a 1:1 Scala rewrite of every Java SDK class.

What is still not 100%:

* There are no generated Scala equivalents for all 546 Java Management API client classes.
* There are no generated Scala equivalents for all 2,369 Java Management API request and response model classes.
* Management API endpoint reachability exists through named resources, nested helpers, and generic HTTP methods, but every Java generated method overload is not reproduced with the same exact Scala method signature.
* Java endpoint defaults from generated request classes are not all encoded as dedicated Scala request types.
* Java generated pagination types such as operation-specific page wrappers are not reproduced as dedicated Scala classes.
* Endpoint-specific Java error model payload classes are not all reproduced as Scala classes.
* The Scala error hierarchy covers common HTTP status categories, but not every generated Java error schema type.
* Async support is Scala `Future` based, not a generated async client class for every Java client.
* Transport behavior is implemented with Java `HttpClient`, not OkHttp, so retry, logging, timeout, and proxy behavior are functionally available but not implementation-identical to Java.
* Java binary compatibility is not provided for the Scala API surface.
* Java annotations and generated Fern metadata are not reproduced in Scala.
* Java builder classes for every generated Management API request model are not reproduced.
* Java raw client classes for every generated Management API group are not reproduced as dedicated Scala classes.
* Java typed file stream and response body wrapper classes are not reproduced with identical names and behavior.
* Java SSE helper names and internals are not reproduced exactly, though SSE parsing exists.
* Java custom domain interception is matched by whitelist behavior, but not by the same interceptor implementation.
* Java telemetry header behavior is not reproduced.
* Java test-only APIs are not reproduced.
* Java package names are not mirrored exactly.
* Java deprecation markers and generated documentation blocks are not copied into Scala.

For strict Java typed parity, use the Auth0 Java generated clients and model classes from Scala through JVM interop. This SDK keeps the Scala surface smaller and idiomatic while preserving endpoint reachability.

## Build

```bash
sbt test
```
