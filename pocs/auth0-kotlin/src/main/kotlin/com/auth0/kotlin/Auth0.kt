package com.auth0.kotlin

import com.auth0.client.auth.AuthAPI
import com.auth0.client.auth.AuthorizeUrlBuilder
import com.auth0.client.auth.ClientAssertionSigner
import com.auth0.client.auth.LogoutUrlBuilder
import com.auth0.client.mgmt.ManagementApi
import com.auth0.client.mgmt.ManagementApiBuilder
import com.auth0.client.mgmt.core.Environment
import com.auth0.client.mgmt.core.LogConfig
import com.auth0.net.Request
import com.auth0.net.Response
import com.auth0.net.client.Auth0HttpClient
import kotlin.coroutines.resume
import kotlin.coroutines.resumeWithException
import kotlinx.coroutines.suspendCancellableCoroutine
import okhttp3.OkHttpClient

object Auth0 {
    fun auth(
        domain: String,
        clientId: String,
        configure: AuthenticationConfig.() -> Unit = {}
    ): AuthClient =
        AuthAPI.newBuilder(domain, clientId).configured(configure)

    fun auth(
        domain: String,
        clientId: String,
        clientSecret: String,
        configure: AuthenticationConfig.() -> Unit = {}
    ): AuthClient =
        AuthAPI.newBuilder(domain, clientId, clientSecret).configured(configure)

    fun auth(
        domain: String,
        clientId: String,
        signer: ClientAssertionSigner,
        configure: AuthenticationConfig.() -> Unit = {}
    ): AuthClient =
        AuthAPI.newBuilder(domain, clientId, signer).configured(configure)

    fun managementWithToken(
        domain: String,
        token: String,
        configure: ManagementConfig.() -> Unit = {}
    ): ManagementClient =
        ManagementApi.builder().domain(domain).token(token).configured(configure)

    fun managementWithClientCredentials(
        domain: String,
        clientId: String,
        clientSecret: String,
        configure: ManagementConfig.() -> Unit = {}
    ): ManagementClient =
        ManagementApi.builder().domain(domain).clientCredentials(clientId, clientSecret).configured(configure)

    fun management(configure: ManagementConfig.() -> Unit): ManagementClient =
        ManagementApi.builder().configured(configure)
}

class AuthenticationConfig internal constructor(private val builder: AuthAPI.Builder) {
    fun httpClient(client: Auth0HttpClient) {
        builder.withHttpClient(client)
    }
}

class ManagementConfig internal constructor(private val builder: ManagementApiBuilder) {
    fun domain(value: String) {
        builder.domain(value)
    }

    fun url(value: String) {
        builder.url(value)
    }

    fun environment(value: Environment) {
        builder.environment(value)
    }

    fun token(value: String) {
        builder.token(value)
    }

    fun clientCredentials(clientId: String, clientSecret: String) {
        builder.clientCredentials(clientId, clientSecret)
    }

    fun audience(value: String) {
        builder.audience(value)
    }

    fun timeout(seconds: Int) {
        builder.timeout(seconds)
    }

    fun maxRetries(value: Int) {
        builder.maxRetries(value)
    }

    fun customDomain(value: String) {
        builder.customDomain(value)
    }

    fun httpClient(client: OkHttpClient) {
        builder.httpClient(client)
    }

    fun logging(config: LogConfig) {
        builder.logging(config)
    }

    fun header(name: String, value: String) {
        builder.addHeader(name, value)
    }
}

class AuthClient(val java: AuthAPI) {
    fun authorizationUrl(
        redirectUri: String,
        configure: AuthorizeUrlBuilder.() -> Unit = {}
    ): String =
        java.authorizeUrl(redirectUri).apply(configure).build()

    fun logoutUrl(
        returnToUrl: String,
        includeClientId: Boolean = true,
        configure: LogoutUrlBuilder.() -> Unit = {}
    ): String =
        java.logoutUrl(returnToUrl, includeClientId).apply(configure).build()

    inline fun <T> use(block: AuthAPI.() -> T): T =
        java.block()
}

class ManagementClient(val java: ManagementApi) {
    val actions get() = java.actions()
    val branding get() = java.branding()
    val clientGrants get() = java.clientGrants()
    val clients get() = java.clients()
    val connectionProfiles get() = java.connectionProfiles()
    val connections get() = java.connections()
    val customDomains get() = java.customDomains()
    val deviceCredentials get() = java.deviceCredentials()
    val emailTemplates get() = java.emailTemplates()
    val eventStreams get() = java.eventStreams()
    val events get() = java.events()
    val flows get() = java.flows()
    val forms get() = java.forms()
    val userGrants get() = java.userGrants()
    val groups get() = java.groups()
    val hooks get() = java.hooks()
    val jobs get() = java.jobs()
    val logStreams get() = java.logStreams()
    val logs get() = java.logs()
    val networkAcls get() = java.networkAcls()
    val organizations get() = java.organizations()
    val prompts get() = java.prompts()
    val rateLimitPolicies get() = java.rateLimitPolicies()
    val refreshTokens get() = java.refreshTokens()
    val resourceServers get() = java.resourceServers()
    val roles get() = java.roles()
    val rules get() = java.rules()
    val rulesConfigs get() = java.rulesConfigs()
    val selfServiceProfiles get() = java.selfServiceProfiles()
    val sessions get() = java.sessions()
    val stats get() = java.stats()
    val supplementalSignals get() = java.supplementalSignals()
    val tickets get() = java.tickets()
    val tokenExchangeProfiles get() = java.tokenExchangeProfiles()
    val userAttributeProfiles get() = java.userAttributeProfiles()
    val userBlocks get() = java.userBlocks()
    val users get() = java.users()
    val anomaly get() = java.anomaly()
    val attackProtection get() = java.attackProtection()
    val emails get() = java.emails()
    val guardian get() = java.guardian()
    val keys get() = java.keys()
    val riskAssessments get() = java.riskAssessments()
    val tenants get() = java.tenants()
    val verifiableCredentials get() = java.verifiableCredentials()

    inline fun <T> use(block: ManagementApi.() -> T): T =
        java.block()
}

fun AuthAPI.asKotlin(): AuthClient =
    AuthClient(this)

fun ManagementApi.asKotlin(): ManagementClient =
    ManagementClient(this)

fun <T> Request<T>.executeBody(): T? =
    execute().body

suspend fun <T> Request<T>.await(): Response<T> =
    suspendCancellableCoroutine { continuation ->
        val future = executeAsync()
        continuation.invokeOnCancellation { future.cancel(true) }
        future.whenComplete { response, error ->
            if (error == null) {
                continuation.resume(response)
            } else {
                continuation.resumeWithException(error.cause ?: error)
            }
        }
    }

suspend fun <T> Request<T>.awaitBody(): T? =
    await().body

private fun AuthAPI.Builder.configured(configure: AuthenticationConfig.() -> Unit): AuthClient {
    AuthenticationConfig(this).configure()
    return AuthClient(build())
}

private fun ManagementApiBuilder.configured(configure: ManagementConfig.() -> Unit): ManagementClient {
    ManagementConfig(this).configure()
    return ManagementClient(build())
}
