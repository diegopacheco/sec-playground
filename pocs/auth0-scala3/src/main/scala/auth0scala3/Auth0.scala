package auth0scala3

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
import okhttp3.OkHttpClient
import scala.concurrent.ExecutionContext
import scala.concurrent.Future
import scala.jdk.FutureConverters.*

object Auth0:
  def auth(domain: String, clientId: String): AuthClient =
    auth(domain, clientId, _ => ())

  def auth(
      domain: String,
      clientId: String,
      configure: AuthenticationConfig => Unit
  ): AuthClient =
    configured(AuthAPI.newBuilder(domain, clientId), configure)

  def auth(domain: String, clientId: String, clientSecret: String): AuthClient =
    auth(domain, clientId, clientSecret, _ => ())

  def auth(
      domain: String,
      clientId: String,
      clientSecret: String,
      configure: AuthenticationConfig => Unit
  ): AuthClient =
    configured(AuthAPI.newBuilder(domain, clientId, clientSecret), configure)

  def auth(domain: String, clientId: String, signer: ClientAssertionSigner): AuthClient =
    auth(domain, clientId, signer, _ => ())

  def auth(
      domain: String,
      clientId: String,
      signer: ClientAssertionSigner,
      configure: AuthenticationConfig => Unit
  ): AuthClient =
    configured(AuthAPI.newBuilder(domain, clientId, signer), configure)

  def managementWithToken(
      domain: String,
      token: String,
      configure: ManagementConfig => Unit = _ => ()
  ): ManagementClient =
    configured(ManagementApi.builder().domain(domain).token(token), configure)

  def managementWithClientCredentials(
      domain: String,
      clientId: String,
      clientSecret: String,
      configure: ManagementConfig => Unit = _ => ()
  ): ManagementClient =
    configured(ManagementApi.builder().domain(domain).clientCredentials(clientId, clientSecret), configure)

  def management(configure: ManagementConfig => Unit): ManagementClient =
    configured(ManagementApi.builder(), configure)

  private def configured(
      builder: AuthAPI.Builder,
      configure: AuthenticationConfig => Unit
  ): AuthClient =
    configure(AuthenticationConfig(builder))
    AuthClient(builder.build())

  private def configured(
      builder: ManagementApiBuilder,
      configure: ManagementConfig => Unit
  ): ManagementClient =
    configure(ManagementConfig(builder))
    ManagementClient(builder.build())

final class AuthenticationConfig private (builder: AuthAPI.Builder):
  def httpClient(client: Auth0HttpClient): Unit =
    builder.withHttpClient(client)

object AuthenticationConfig:
  private[auth0scala3] def apply(builder: AuthAPI.Builder): AuthenticationConfig =
    new AuthenticationConfig(builder)

final class ManagementConfig private (builder: ManagementApiBuilder):
  def domain(value: String): Unit =
    builder.domain(value)

  def url(value: String): Unit =
    builder.url(value)

  def environment(value: Environment): Unit =
    builder.environment(value)

  def token(value: String): Unit =
    builder.token(value)

  def clientCredentials(clientId: String, clientSecret: String): Unit =
    builder.clientCredentials(clientId, clientSecret)

  def audience(value: String): Unit =
    builder.audience(value)

  def timeout(seconds: Int): Unit =
    builder.timeout(seconds)

  def maxRetries(value: Int): Unit =
    builder.maxRetries(value)

  def customDomain(value: String): Unit =
    builder.customDomain(value)

  def httpClient(client: OkHttpClient): Unit =
    builder.httpClient(client)

  def logging(config: LogConfig): Unit =
    builder.logging(config)

  def header(name: String, value: String): Unit =
    builder.addHeader(name, value)

object ManagementConfig:
  private[auth0scala3] def apply(builder: ManagementApiBuilder): ManagementConfig =
    new ManagementConfig(builder)

final class AuthorizationConfig private (builder: AuthorizeUrlBuilder):
  def connection(value: String): Unit =
    builder.withConnection(value)

  def audience(value: String): Unit =
    builder.withAudience(value)

  def state(value: String): Unit =
    builder.withState(value)

  def scope(value: String): Unit =
    builder.withScope(value)

  def responseType(value: String): Unit =
    builder.withResponseType(value)

  def organization(value: String): Unit =
    builder.withOrganization(value)

  def invitation(value: String): Unit =
    builder.withInvitation(value)

  def parameter(name: String, value: String): Unit =
    builder.withParameter(name, value)

  def codeChallenge(value: String): Unit =
    builder.withCodeChallenge(value)

object AuthorizationConfig:
  private[auth0scala3] def apply(builder: AuthorizeUrlBuilder): AuthorizationConfig =
    new AuthorizationConfig(builder)

final class LogoutConfig private (builder: LogoutUrlBuilder):
  def federated(value: Boolean = true): Unit =
    builder.useFederated(value)

object LogoutConfig:
  private[auth0scala3] def apply(builder: LogoutUrlBuilder): LogoutConfig =
    new LogoutConfig(builder)

final class AuthClient(val java: AuthAPI):
  def authorizationUrl(
      redirectUri: String,
      configure: AuthorizationConfig => Unit = _ => ()
  ): String =
    val builder = java.authorizeUrl(redirectUri)
    configure(AuthorizationConfig(builder))
    builder.build()

  def logoutUrl(
      returnToUrl: String,
      includeClientId: Boolean = true,
      configure: LogoutConfig => Unit = _ => ()
  ): String =
    val builder = java.logoutUrl(returnToUrl, includeClientId)
    configure(LogoutConfig(builder))
    builder.build()

  def use[A](operation: AuthAPI => A): A =
    operation(java)

object AuthClient:
  def apply(java: AuthAPI): AuthClient =
    new AuthClient(java)

final class ManagementClient(val java: ManagementApi):
  export java.{
    actions,
    anomaly,
    attackProtection,
    branding,
    clientGrants,
    clients,
    connectionProfiles,
    connections,
    customDomains,
    deviceCredentials,
    emailTemplates,
    emails,
    eventStreams,
    events,
    flows,
    forms,
    groups,
    guardian,
    hooks,
    jobs,
    keys,
    logStreams,
    logs,
    networkAcls,
    organizations,
    prompts,
    rateLimitPolicies,
    refreshTokens,
    resourceServers,
    riskAssessments,
    roles,
    rules,
    rulesConfigs,
    selfServiceProfiles,
    sessions,
    stats,
    supplementalSignals,
    tenants,
    tickets,
    tokenExchangeProfiles,
    userAttributeProfiles,
    userBlocks,
    users,
    verifiableCredentials
  }

  def use[A](operation: ManagementApi => A): A =
    operation(java)

object ManagementClient:
  def apply(java: ManagementApi): ManagementClient =
    new ManagementClient(java)

extension (java: AuthAPI)
  def asScala: AuthClient =
    AuthClient(java)

extension (java: ManagementApi)
  def asScala: ManagementClient =
    ManagementClient(java)

extension [A](request: Request[A])
  def executeBody(): A =
    request.execute().getBody

  def future: Future[Response[A]] =
    request.executeAsync().asScala

  def bodyFuture(using ExecutionContext): Future[A] =
    future.map(_.getBody)

extension [A](response: Response[A])
  def statusCode: Int =
    response.getStatusCode

  def bodyOption: Option[A] =
    Option(response.getBody)
