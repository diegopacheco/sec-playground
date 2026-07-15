package com.auth0.kotlin

import com.auth0.client.auth.AuthAPI
import com.auth0.client.auth.AuthorizeUrlBuilder
import com.auth0.client.auth.ClientAssertionSigner
import com.auth0.client.auth.LogoutUrlBuilder
import com.auth0.client.mgmt.ActionsClient
import com.auth0.client.mgmt.BrandingClient
import com.auth0.client.mgmt.ClientGrantsClient
import com.auth0.client.mgmt.ClientsClient
import com.auth0.client.mgmt.ConnectionProfilesClient
import com.auth0.client.mgmt.ConnectionsClient
import com.auth0.client.mgmt.CustomDomainsClient
import com.auth0.client.mgmt.DeviceCredentialsClient
import com.auth0.client.mgmt.EmailTemplatesClient
import com.auth0.client.mgmt.EventStreamsClient
import com.auth0.client.mgmt.EventsClient
import com.auth0.client.mgmt.FlowsClient
import com.auth0.client.mgmt.FormsClient
import com.auth0.client.mgmt.GroupsClient
import com.auth0.client.mgmt.HooksClient
import com.auth0.client.mgmt.JobsClient
import com.auth0.client.mgmt.LogStreamsClient
import com.auth0.client.mgmt.LogsClient
import com.auth0.client.mgmt.ManagementApi
import com.auth0.client.mgmt.ManagementApiBuilder
import com.auth0.client.mgmt.NetworkAclsClient
import com.auth0.client.mgmt.OrganizationsClient
import com.auth0.client.mgmt.PromptsClient
import com.auth0.client.mgmt.RateLimitPoliciesClient
import com.auth0.client.mgmt.RefreshTokensClient
import com.auth0.client.mgmt.ResourceServersClient
import com.auth0.client.mgmt.RolesClient
import com.auth0.client.mgmt.RulesClient
import com.auth0.client.mgmt.RulesConfigsClient
import com.auth0.client.mgmt.SelfServiceProfilesClient
import com.auth0.client.mgmt.SessionsClient
import com.auth0.client.mgmt.StatsClient
import com.auth0.client.mgmt.SupplementalSignalsClient
import com.auth0.client.mgmt.TicketsClient
import com.auth0.client.mgmt.TokenExchangeProfilesClient
import com.auth0.client.mgmt.UserAttributeProfilesClient
import com.auth0.client.mgmt.UserBlocksClient
import com.auth0.client.mgmt.UserGrantsClient
import com.auth0.client.mgmt.UsersClient
import com.auth0.client.mgmt.anomaly.AnomalyClient
import com.auth0.client.mgmt.attackprotection.AttackProtectionClient
import com.auth0.client.mgmt.core.Environment
import com.auth0.client.mgmt.core.LogConfig
import com.auth0.client.mgmt.emails.EmailsClient
import com.auth0.client.mgmt.guardian.GuardianClient
import com.auth0.client.mgmt.keys.KeysClient
import com.auth0.client.mgmt.riskassessments.RiskAssessmentsClient
import com.auth0.client.mgmt.tenants.TenantsClient
import com.auth0.client.mgmt.verifiablecredentials.VerifiableCredentialsClient
import com.auth0.net.Request
import com.auth0.net.Response
import com.auth0.net.client.Auth0HttpClient
import kotlin.coroutines.resume
import kotlin.coroutines.resumeWithException
import kotlinx.coroutines.suspendCancellableCoroutine
import okhttp3.OkHttpClient

@DslMarker
public annotation class Auth0Dsl

public object Auth0 {
    public fun auth(
        domain: String,
        clientId: String,
        configure: AuthenticationConfig.() -> Unit = {},
    ): AuthClient =
        AuthAPI.newBuilder(domain, clientId).configured(configure)

    public fun auth(
        domain: String,
        clientId: String,
        clientSecret: String,
        configure: AuthenticationConfig.() -> Unit = {},
    ): AuthClient =
        AuthAPI.newBuilder(domain, clientId, clientSecret).configured(configure)

    public fun auth(
        domain: String,
        clientId: String,
        signer: ClientAssertionSigner,
        configure: AuthenticationConfig.() -> Unit = {},
    ): AuthClient =
        AuthAPI.newBuilder(domain, clientId, signer).configured(configure)

    public fun managementWithToken(
        domain: String,
        token: String,
        configure: ManagementConfig.() -> Unit = {},
    ): ManagementClient =
        ManagementApi.builder().domain(domain).token(token).configured(configure)

    public fun managementWithClientCredentials(
        domain: String,
        clientId: String,
        clientSecret: String,
        configure: ManagementConfig.() -> Unit = {},
    ): ManagementClient =
        ManagementApi.builder().domain(domain).clientCredentials(clientId, clientSecret).configured(configure)

    public fun management(configure: ManagementConfig.() -> Unit): ManagementClient =
        ManagementApi.builder().configured(configure)
}

@Auth0Dsl
public class AuthenticationConfig internal constructor(private val builder: AuthAPI.Builder) {
    public fun httpClient(client: Auth0HttpClient): Unit {
        builder.withHttpClient(client)
    }
}

@Auth0Dsl
public class ManagementConfig internal constructor(private val builder: ManagementApiBuilder) {
    public fun domain(value: String): Unit {
        builder.domain(value)
    }

    public fun url(value: String): Unit {
        builder.url(value)
    }

    public fun environment(value: Environment): Unit {
        builder.environment(value)
    }

    public fun token(value: String): Unit {
        builder.token(value)
    }

    public fun clientCredentials(clientId: String, clientSecret: String): Unit {
        builder.clientCredentials(clientId, clientSecret)
    }

    public fun audience(value: String): Unit {
        builder.audience(value)
    }

    public fun timeout(seconds: Int): Unit {
        builder.timeout(seconds)
    }

    public fun maxRetries(value: Int): Unit {
        builder.maxRetries(value)
    }

    public fun customDomain(value: String): Unit {
        builder.customDomain(value)
    }

    public fun httpClient(client: OkHttpClient): Unit {
        builder.httpClient(client)
    }

    public fun logging(config: LogConfig): Unit {
        builder.logging(config)
    }

    public fun header(name: String, value: String): Unit {
        builder.addHeader(name, value)
    }
}

public class AuthClient(public val java: AuthAPI) {
    public fun authorizationUrl(
        redirectUri: String,
        configure: AuthorizeUrlBuilder.() -> Unit = {},
    ): String =
        java.authorizeUrl(redirectUri).apply(configure).build()

    public fun logoutUrl(
        returnToUrl: String,
        includeClientId: Boolean = true,
        configure: LogoutUrlBuilder.() -> Unit = {},
    ): String =
        java.logoutUrl(returnToUrl, includeClientId).apply(configure).build()

    public inline fun <T> use(block: AuthAPI.() -> T): T =
        java.block()
}

public class ManagementClient(public val java: ManagementApi) {
    public val actions: ActionsClient get() = java.actions()
    public val branding: BrandingClient get() = java.branding()
    public val clientGrants: ClientGrantsClient get() = java.clientGrants()
    public val clients: ClientsClient get() = java.clients()
    public val connectionProfiles: ConnectionProfilesClient get() = java.connectionProfiles()
    public val connections: ConnectionsClient get() = java.connections()
    public val customDomains: CustomDomainsClient get() = java.customDomains()
    public val deviceCredentials: DeviceCredentialsClient get() = java.deviceCredentials()
    public val emailTemplates: EmailTemplatesClient get() = java.emailTemplates()
    public val eventStreams: EventStreamsClient get() = java.eventStreams()
    public val events: EventsClient get() = java.events()
    public val flows: FlowsClient get() = java.flows()
    public val forms: FormsClient get() = java.forms()
    public val userGrants: UserGrantsClient get() = java.userGrants()
    public val groups: GroupsClient get() = java.groups()
    public val hooks: HooksClient get() = java.hooks()
    public val jobs: JobsClient get() = java.jobs()
    public val logStreams: LogStreamsClient get() = java.logStreams()
    public val logs: LogsClient get() = java.logs()
    public val networkAcls: NetworkAclsClient get() = java.networkAcls()
    public val organizations: OrganizationsClient get() = java.organizations()
    public val prompts: PromptsClient get() = java.prompts()
    public val rateLimitPolicies: RateLimitPoliciesClient get() = java.rateLimitPolicies()
    public val refreshTokens: RefreshTokensClient get() = java.refreshTokens()
    public val resourceServers: ResourceServersClient get() = java.resourceServers()
    public val roles: RolesClient get() = java.roles()
    public val rules: RulesClient get() = java.rules()
    public val rulesConfigs: RulesConfigsClient get() = java.rulesConfigs()
    public val selfServiceProfiles: SelfServiceProfilesClient get() = java.selfServiceProfiles()
    public val sessions: SessionsClient get() = java.sessions()
    public val stats: StatsClient get() = java.stats()
    public val supplementalSignals: SupplementalSignalsClient get() = java.supplementalSignals()
    public val tickets: TicketsClient get() = java.tickets()
    public val tokenExchangeProfiles: TokenExchangeProfilesClient get() = java.tokenExchangeProfiles()
    public val userAttributeProfiles: UserAttributeProfilesClient get() = java.userAttributeProfiles()
    public val userBlocks: UserBlocksClient get() = java.userBlocks()
    public val users: UsersClient get() = java.users()
    public val anomaly: AnomalyClient get() = java.anomaly()
    public val attackProtection: AttackProtectionClient get() = java.attackProtection()
    public val emails: EmailsClient get() = java.emails()
    public val guardian: GuardianClient get() = java.guardian()
    public val keys: KeysClient get() = java.keys()
    public val riskAssessments: RiskAssessmentsClient get() = java.riskAssessments()
    public val tenants: TenantsClient get() = java.tenants()
    public val verifiableCredentials: VerifiableCredentialsClient get() = java.verifiableCredentials()

    public inline fun <T> use(block: ManagementApi.() -> T): T =
        java.block()
}

public fun AuthAPI.asKotlin(): AuthClient =
    AuthClient(this)

public fun ManagementApi.asKotlin(): ManagementClient =
    ManagementClient(this)

public fun <T> Request<T>.executeBody(): T? =
    execute().body

public suspend fun <T> Request<T>.await(): Response<T> =
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

public suspend fun <T> Request<T>.awaitBody(): T? =
    await().body

private fun AuthAPI.Builder.configured(configure: AuthenticationConfig.() -> Unit): AuthClient {
    AuthenticationConfig(this).configure()
    return AuthClient(build())
}

private fun ManagementApiBuilder.configured(configure: ManagementConfig.() -> Unit): ManagementClient {
    ManagementConfig(this).configure()
    return ManagementClient(build())
}
