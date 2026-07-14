package com.auth0.kotlin

import com.auth0.client.auth.AuthAPI
import com.auth0.client.auth.AuthorizeUrlBuilder
import com.auth0.client.auth.ClientAssertionSigner
import com.auth0.client.mgmt.ManagementApi
import com.auth0.client.mgmt.core.Environment
import com.auth0.client.mgmt.core.LogConfig
import com.auth0.json.ObjectMapperProvider
import com.auth0.net.Request
import com.auth0.net.Response
import com.auth0.net.client.Auth0HttpClient
import com.fasterxml.jackson.module.kotlin.readValue
import com.fasterxml.jackson.module.kotlin.registerKotlinModule
import java.lang.reflect.Method
import java.lang.reflect.Modifier
import java.util.concurrent.CompletableFuture
import kotlin.coroutines.resume
import kotlin.coroutines.resumeWithException
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.suspendCancellableCoroutine
import kotlinx.coroutines.withContext
import okhttp3.OkHttpClient

object Auth0 {
    @JvmStatic
    fun auth(domain: String, clientId: String, block: AuthApiConfig.() -> Unit = {}): AuthAPI =
        AuthApiConfig().apply(block).build(domain, clientId)

    @JvmStatic
    fun auth(domain: String, clientId: String, clientSecret: String, block: AuthApiConfig.() -> Unit = {}): AuthAPI =
        AuthApiConfig().apply {
            this.clientSecret = clientSecret
            block()
        }.build(domain, clientId)

    @JvmStatic
    fun management(block: ManagementApiConfig.() -> Unit): ManagementApi =
        ManagementApiConfig().apply(block).build()

    @JvmStatic
    fun managementWithToken(
        domain: String,
        token: String,
        block: ManagementApiConfig.() -> Unit = {}
    ): ManagementApi =
        ManagementApiConfig().apply {
            this.domain = domain
            this.token = token
            block()
        }.build()

    @JvmStatic
    fun managementWithClientCredentials(
        domain: String,
        clientId: String,
        clientSecret: String,
        block: ManagementApiConfig.() -> Unit = {}
    ): ManagementApi =
        ManagementApiConfig().apply {
            this.domain = domain
            this.clientId = clientId
            this.clientSecret = clientSecret
            block()
        }.build()

    @JvmStatic
    fun kotlinAuth(domain: String, clientId: String, block: AuthApiConfig.() -> Unit = {}): KotlinAuthApi =
        AuthApiConfig().apply(block).build(domain, clientId).kotlin

    @JvmStatic
    fun kotlinAuth(
        domain: String,
        clientId: String,
        clientSecret: String,
        block: AuthApiConfig.() -> Unit = {}
    ): KotlinAuthApi =
        AuthApiConfig().apply {
            this.clientSecret = clientSecret
            block()
        }.build(domain, clientId).kotlin

    @JvmStatic
    fun kotlinManagement(block: ManagementApiConfig.() -> Unit): KotlinManagementApi =
        ManagementApiConfig().apply(block).build().kotlin

    @JvmStatic
    fun kotlinManagementWithToken(
        domain: String,
        token: String,
        block: ManagementApiConfig.() -> Unit = {}
    ): KotlinManagementApi =
        managementWithToken(domain, token, block).kotlin

    @JvmStatic
    fun kotlinManagementWithClientCredentials(
        domain: String,
        clientId: String,
        clientSecret: String,
        block: ManagementApiConfig.() -> Unit = {}
    ): KotlinManagementApi =
        managementWithClientCredentials(domain, clientId, clientSecret, block).kotlin
}

class AuthApiConfig {
    var clientSecret: String? = null
    var clientAssertionSigner: ClientAssertionSigner? = null
    var httpClient: Auth0HttpClient? = null

    internal fun build(domain: String, clientId: String): AuthAPI {
        val builder = when {
            clientAssertionSigner != null -> AuthAPI.newBuilder(domain, clientId, clientAssertionSigner)
            clientSecret != null -> AuthAPI.newBuilder(domain, clientId, clientSecret)
            else -> AuthAPI.newBuilder(domain, clientId)
        }
        httpClient?.let(builder::withHttpClient)
        return builder.build()
    }
}

class ManagementApiConfig {
    var domain: String? = null
    var url: String? = null
    var environment: Environment? = null
    var token: String? = null
    var clientId: String? = null
    var clientSecret: String? = null
    var audience: String? = null
    var timeout: Int? = null
    var maxRetries: Int? = null
    var customDomain: String? = null
    var httpClient: OkHttpClient? = null
    var logging: LogConfig? = null

    private val headers = linkedMapOf<String, String>()

    fun clientCredentials(clientId: String, clientSecret: String) {
        this.clientId = clientId
        this.clientSecret = clientSecret
    }

    fun header(name: String, value: String) {
        headers[name] = value
    }

    internal fun build(): ManagementApi {
        val builder = ManagementApi.builder()
        domain?.let(builder::domain)
        url?.let(builder::url)
        environment?.let(builder::environment)
        token?.let(builder::token)
        audience?.let(builder::audience)
        timeout?.let(builder::timeout)
        maxRetries?.let(builder::maxRetries)
        customDomain?.let(builder::customDomain)
        httpClient?.let(builder::httpClient)
        logging?.let(builder::logging)
        headers.forEach(builder::addHeader)
        if (clientId != null && clientSecret != null) {
            builder.clientCredentials(clientId, clientSecret)
        }
        return builder.build()
    }
}

data class Auth0Response<T>(
    val statusCode: Int,
    val headers: Map<String, String>,
    val body: T
)

data class Auth0NullableResponse<T>(
    val statusCode: Int,
    val headers: Map<String, String>,
    val body: T?
)

val AuthAPI.kotlin: KotlinAuthApi
    get() = KotlinAuthApi(this)

val ManagementApi.kotlin: KotlinManagementApi
    get() = KotlinManagementApi(this)

class KotlinAuthApi(val java: AuthAPI) {
    fun authorizeUrl(redirectUri: String, block: AuthorizeUrlBuilder.() -> Unit = {}): String =
        java.authorizeUrl(redirectUri).apply(block).build()

    fun logoutUrl(returnToUrl: String, setClientId: Boolean = true, federated: Boolean = false): String =
        java.logoutUrl(returnToUrl, setClientId).useFederated(federated).build()

    @Suppress("UNCHECKED_CAST")
    fun <T : Any> request(name: String, vararg args: Any?): Request<T> =
        ReflectionCall.invoke(java, name, args) as Request<T>

    suspend fun <T : Any> await(name: String, vararg args: Any?): T =
        request<T>(name, *args).awaitBody()

    @Suppress("UNCHECKED_CAST")
    fun <T : Any> call(name: String, vararg args: Any?): T =
        ReflectionCall.invoke(java, name, args) as T

    fun operation(name: String, vararg args: Any?, block: KotlinOperation.() -> Unit = {}): Any? =
        KotlinOperation(ReflectionCall.invoke(java, name, args)).apply(block).value
}

class KotlinManagementApi(val java: ManagementApi) {
    fun actions(block: KotlinEndpointClient.() -> Unit = {}): KotlinEndpointClient = resource("actions", block)
    fun branding(block: KotlinEndpointClient.() -> Unit = {}): KotlinEndpointClient = resource("branding", block)
    fun clientGrants(block: KotlinEndpointClient.() -> Unit = {}): KotlinEndpointClient = resource("clientGrants", block)
    fun clients(block: KotlinEndpointClient.() -> Unit = {}): KotlinEndpointClient = resource("clients", block)
    fun connectionProfiles(block: KotlinEndpointClient.() -> Unit = {}): KotlinEndpointClient = resource("connectionProfiles", block)
    fun connections(block: KotlinEndpointClient.() -> Unit = {}): KotlinEndpointClient = resource("connections", block)
    fun customDomains(block: KotlinEndpointClient.() -> Unit = {}): KotlinEndpointClient = resource("customDomains", block)
    fun deviceCredentials(block: KotlinEndpointClient.() -> Unit = {}): KotlinEndpointClient = resource("deviceCredentials", block)
    fun emailTemplates(block: KotlinEndpointClient.() -> Unit = {}): KotlinEndpointClient = resource("emailTemplates", block)
    fun eventStreams(block: KotlinEndpointClient.() -> Unit = {}): KotlinEndpointClient = resource("eventStreams", block)
    fun events(block: KotlinEndpointClient.() -> Unit = {}): KotlinEndpointClient = resource("events", block)
    fun flows(block: KotlinEndpointClient.() -> Unit = {}): KotlinEndpointClient = resource("flows", block)
    fun forms(block: KotlinEndpointClient.() -> Unit = {}): KotlinEndpointClient = resource("forms", block)
    fun userGrants(block: KotlinEndpointClient.() -> Unit = {}): KotlinEndpointClient = resource("userGrants", block)
    fun groups(block: KotlinEndpointClient.() -> Unit = {}): KotlinEndpointClient = resource("groups", block)
    fun hooks(block: KotlinEndpointClient.() -> Unit = {}): KotlinEndpointClient = resource("hooks", block)
    fun jobs(block: KotlinEndpointClient.() -> Unit = {}): KotlinEndpointClient = resource("jobs", block)
    fun logStreams(block: KotlinEndpointClient.() -> Unit = {}): KotlinEndpointClient = resource("logStreams", block)
    fun logs(block: KotlinEndpointClient.() -> Unit = {}): KotlinEndpointClient = resource("logs", block)
    fun networkAcls(block: KotlinEndpointClient.() -> Unit = {}): KotlinEndpointClient = resource("networkAcls", block)
    fun organizations(block: KotlinEndpointClient.() -> Unit = {}): KotlinEndpointClient = resource("organizations", block)
    fun prompts(block: KotlinEndpointClient.() -> Unit = {}): KotlinEndpointClient = resource("prompts", block)
    fun rateLimitPolicies(block: KotlinEndpointClient.() -> Unit = {}): KotlinEndpointClient = resource("rateLimitPolicies", block)
    fun refreshTokens(block: KotlinEndpointClient.() -> Unit = {}): KotlinEndpointClient = resource("refreshTokens", block)
    fun resourceServers(block: KotlinEndpointClient.() -> Unit = {}): KotlinEndpointClient = resource("resourceServers", block)
    fun roles(block: KotlinEndpointClient.() -> Unit = {}): KotlinEndpointClient = resource("roles", block)
    fun rules(block: KotlinEndpointClient.() -> Unit = {}): KotlinEndpointClient = resource("rules", block)
    fun rulesConfigs(block: KotlinEndpointClient.() -> Unit = {}): KotlinEndpointClient = resource("rulesConfigs", block)
    fun selfServiceProfiles(block: KotlinEndpointClient.() -> Unit = {}): KotlinEndpointClient = resource("selfServiceProfiles", block)
    fun sessions(block: KotlinEndpointClient.() -> Unit = {}): KotlinEndpointClient = resource("sessions", block)
    fun stats(block: KotlinEndpointClient.() -> Unit = {}): KotlinEndpointClient = resource("stats", block)
    fun supplementalSignals(block: KotlinEndpointClient.() -> Unit = {}): KotlinEndpointClient = resource("supplementalSignals", block)
    fun tickets(block: KotlinEndpointClient.() -> Unit = {}): KotlinEndpointClient = resource("tickets", block)
    fun tokenExchangeProfiles(block: KotlinEndpointClient.() -> Unit = {}): KotlinEndpointClient = resource("tokenExchangeProfiles", block)
    fun userAttributeProfiles(block: KotlinEndpointClient.() -> Unit = {}): KotlinEndpointClient = resource("userAttributeProfiles", block)
    fun userBlocks(block: KotlinEndpointClient.() -> Unit = {}): KotlinEndpointClient = resource("userBlocks", block)
    fun users(block: KotlinEndpointClient.() -> Unit = {}): KotlinEndpointClient = resource("users", block)
    fun anomaly(block: KotlinEndpointClient.() -> Unit = {}): KotlinEndpointClient = resource("anomaly", block)
    fun attackProtection(block: KotlinEndpointClient.() -> Unit = {}): KotlinEndpointClient = resource("attackProtection", block)
    fun emails(block: KotlinEndpointClient.() -> Unit = {}): KotlinEndpointClient = resource("emails", block)
    fun guardian(block: KotlinEndpointClient.() -> Unit = {}): KotlinEndpointClient = resource("guardian", block)
    fun keys(block: KotlinEndpointClient.() -> Unit = {}): KotlinEndpointClient = resource("keys", block)
    fun riskAssessments(block: KotlinEndpointClient.() -> Unit = {}): KotlinEndpointClient = resource("riskAssessments", block)
    fun tenants(block: KotlinEndpointClient.() -> Unit = {}): KotlinEndpointClient = resource("tenants", block)
    fun verifiableCredentials(block: KotlinEndpointClient.() -> Unit = {}): KotlinEndpointClient = resource("verifiableCredentials", block)

    fun resource(name: String, block: KotlinEndpointClient.() -> Unit = {}): KotlinEndpointClient =
        KotlinEndpointClient(requireNotNull(ReflectionCall.invoke(java, name, emptyArray()))).apply(block)
}

class KotlinEndpointClient(val java: Any) {
    fun raw(block: KotlinEndpointClient.() -> Unit = {}): KotlinEndpointClient =
        resource("withRawResponse", block)

    fun resource(name: String, block: KotlinEndpointClient.() -> Unit = {}): KotlinEndpointClient =
        KotlinEndpointClient(requireNotNull(ReflectionCall.invoke(java, name, emptyArray()))).apply(block)

    @Suppress("UNCHECKED_CAST")
    fun <T : Any> call(name: String, vararg args: Any?): T =
        ReflectionCall.invoke(java, name, args) as T

    fun nullableCall(name: String, vararg args: Any?): Any? =
        ReflectionCall.invoke(java, name, args)

    suspend fun <T : Any> callSuspending(name: String, vararg args: Any?): T =
        withContext(Dispatchers.IO) { call(name, *args) }

    suspend fun nullableCallSuspending(name: String, vararg args: Any?): Any? =
        withContext(Dispatchers.IO) { nullableCall(name, *args) }

    fun operation(name: String, vararg args: Any?, block: KotlinOperation.() -> Unit = {}): Any? =
        KotlinOperation(ReflectionCall.invoke(java, name, args)).apply(block).value
}

class KotlinOperation(val value: Any?) {
    fun endpoint(block: KotlinEndpointClient.() -> Unit): KotlinEndpointClient =
        KotlinEndpointClient(requireNotNull(value)).apply(block)

    @Suppress("UNCHECKED_CAST")
    fun <T : Any> body(): T =
        when (value) {
            is Request<*> -> (value as Request<T>).body()
            else -> requireNotNull(value) as T
        }

    @Suppress("UNCHECKED_CAST")
    suspend fun <T : Any> awaitBody(): T =
        when (value) {
            is Request<*> -> (value as Request<T>).awaitBody()
            else -> requireNotNull(value) as T
        }
}

object Auth0Json {
    val mapper = ObjectMapperProvider.getMapper().copy().registerKotlinModule()

    fun encode(value: Any): String =
        mapper.writeValueAsString(value)

    inline fun <reified T : Any> decode(json: String): T =
        mapper.readValue(json)
}

fun AuthAPI.authorizeUrlFor(
    redirectUri: String,
    block: AuthorizeUrlBuilder.() -> Unit = {}
): String =
    authorizeUrl(redirectUri).apply(block).build()

fun AuthAPI.logoutUrlFor(
    returnToUrl: String,
    setClientId: Boolean = true,
    federated: Boolean = false
): String =
    logoutUrl(returnToUrl, setClientId).useFederated(federated).build()

fun <T> Request<T>.response(): Auth0Response<T> =
    execute().toAuth0Response()

fun <T> Request<T>.body(): T =
    execute().body

fun <T> Request<T>.responseAsync(): CompletableFuture<Auth0Response<T>> =
    executeAsync().thenApply { it.toAuth0Response() }

fun <T> Request<T>.bodyAsync(): CompletableFuture<T> =
    executeAsync().thenApply { it.body }

fun <T> Request<T>.nullableResponse(): Auth0NullableResponse<T> =
    execute().toNullableAuth0Response()

fun <T> Request<T>.nullableBody(): T? =
    execute().body

suspend fun <T : Any> Request<T>.awaitResponse(): Auth0Response<T> =
    suspendCancellableCoroutine { continuation ->
        val future = executeAsync().thenApply { it.toAuth0Response() }
        continuation.invokeOnCancellation { future.cancel(true) }
        future.whenComplete { response, throwable ->
            if (throwable == null) {
                continuation.resume(response)
            } else {
                continuation.resumeWithException(throwable)
            }
        }
    }

suspend fun <T : Any> Request<T>.awaitBody(): T =
    awaitResponse().body

suspend fun <T> Request<T>.awaitNullableResponse(): Auth0NullableResponse<T> =
    suspendCancellableCoroutine { continuation ->
        val future = executeAsync().thenApply { it.toNullableAuth0Response() }
        continuation.invokeOnCancellation { future.cancel(true) }
        future.whenComplete { response, throwable ->
            if (throwable == null) {
                continuation.resume(response)
            } else {
                continuation.resumeWithException(throwable)
            }
        }
    }

suspend fun <T> Request<T>.awaitNullableBody(): T? =
    awaitNullableResponse().body

fun <T> Response<T>.toAuth0Response(): Auth0Response<T> =
    Auth0Response(statusCode, headers, requireNotNull(body))

fun <T> Response<T>.toNullableAuth0Response(): Auth0NullableResponse<T> =
    Auth0NullableResponse(statusCode, headers, body)

operator fun <T> Response<T>.component1(): Int =
    statusCode

operator fun <T> Response<T>.component2(): Map<String, String> =
    headers

operator fun <T> Response<T>.component3(): T =
    body

private object ReflectionCall {
    fun invoke(target: Any, name: String, args: Array<out Any?>): Any? {
        val method = target.javaClass.methods
            .filter { Modifier.isPublic(it.modifiers) }
            .filter { it.name == name }
            .filter { it.parameterCount == args.size }
            .firstOrNull { accepts(it, args) }
            ?: throw IllegalArgumentException("No public method ${target.javaClass.name}.$name with ${args.size} arguments")
        return method.invoke(target, *args)
    }

    private fun accepts(method: Method, args: Array<out Any?>): Boolean =
        method.parameterTypes.zip(args).all { (type, arg) ->
            arg == null && !type.isPrimitive || arg != null && boxed(type).isAssignableFrom(arg.javaClass)
        }

    private fun boxed(type: Class<*>): Class<*> =
        when (type) {
            java.lang.Boolean.TYPE -> java.lang.Boolean::class.java
            java.lang.Byte.TYPE -> java.lang.Byte::class.java
            java.lang.Character.TYPE -> java.lang.Character::class.java
            java.lang.Short.TYPE -> java.lang.Short::class.java
            java.lang.Integer.TYPE -> java.lang.Integer::class.java
            java.lang.Long.TYPE -> java.lang.Long::class.java
            java.lang.Float.TYPE -> java.lang.Float::class.java
            java.lang.Double.TYPE -> java.lang.Double::class.java
            java.lang.Void.TYPE -> java.lang.Void::class.java
            else -> type
        }
}
