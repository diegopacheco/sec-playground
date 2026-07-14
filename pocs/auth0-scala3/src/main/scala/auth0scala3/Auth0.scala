package auth0scala3

import com.auth0.client.auth.AuthAPI
import com.auth0.client.auth.AuthorizeUrlBuilder
import com.auth0.client.auth.ClientAssertionSigner
import com.auth0.client.auth.LogoutUrlBuilder
import com.auth0.client.mgmt.ManagementApi
import com.auth0.client.mgmt.ManagementApiBuilder
import com.auth0.net.Request
import com.auth0.net.Response
import scala.concurrent.ExecutionContext
import scala.concurrent.Future
import scala.jdk.FutureConverters.*

object Auth0:
  def auth(domain: String, clientId: String): AuthClient =
    auth(domain, clientId, identity[AuthAPI.Builder])

  def auth(
      domain: String,
      clientId: String,
      configure: AuthAPI.Builder => AuthAPI.Builder
  ): AuthClient =
    AuthClient(configure(AuthAPI.newBuilder(domain, clientId)).build())

  def auth(domain: String, clientId: String, clientSecret: String): AuthClient =
    auth(domain, clientId, clientSecret, identity[AuthAPI.Builder])

  def auth(
      domain: String,
      clientId: String,
      clientSecret: String,
      configure: AuthAPI.Builder => AuthAPI.Builder
  ): AuthClient =
    AuthClient(configure(AuthAPI.newBuilder(domain, clientId, clientSecret)).build())

  def auth(domain: String, clientId: String, signer: ClientAssertionSigner): AuthClient =
    auth(domain, clientId, signer, identity[AuthAPI.Builder])

  def auth(
      domain: String,
      clientId: String,
      signer: ClientAssertionSigner,
      configure: AuthAPI.Builder => AuthAPI.Builder
  ): AuthClient =
    AuthClient(configure(AuthAPI.newBuilder(domain, clientId, signer)).build())

  def managementWithToken(
      domain: String,
      token: String,
      configure: ManagementApiBuilder => ManagementApiBuilder = identity[ManagementApiBuilder]
  ): ManagementClient =
    ManagementClient(configure(ManagementApi.builder().domain(domain).token(token)).build())

  def managementWithClientCredentials(
      domain: String,
      clientId: String,
      clientSecret: String,
      configure: ManagementApiBuilder => ManagementApiBuilder = identity[ManagementApiBuilder]
  ): ManagementClient =
    ManagementClient(configure(ManagementApi.builder().domain(domain).clientCredentials(clientId, clientSecret)).build())

  def management(configure: ManagementApiBuilder => ManagementApiBuilder): ManagementClient =
    ManagementClient(configure(ManagementApi.builder()).build())

final class AuthClient(val java: AuthAPI):
  def authorizationUrl(
      redirectUri: String,
      configure: AuthorizeUrlBuilder => AuthorizeUrlBuilder = identity[AuthorizeUrlBuilder]
  ): String =
    configure(java.authorizeUrl(redirectUri)).build()

  def logoutUrl(
      returnToUrl: String,
      includeClientId: Boolean = true,
      configure: LogoutUrlBuilder => LogoutUrlBuilder = identity[LogoutUrlBuilder]
  ): String =
    configure(java.logoutUrl(returnToUrl, includeClientId)).build()

  def use[A](operation: AuthAPI => A): A =
    operation(java)

object AuthClient:
  def apply(java: AuthAPI): AuthClient =
    new AuthClient(java)

final class ManagementClient(val java: ManagementApi):
  def actions = java.actions()
  def branding = java.branding()
  def clientGrants = java.clientGrants()
  def clients = java.clients()
  def connectionProfiles = java.connectionProfiles()
  def connections = java.connections()
  def customDomains = java.customDomains()
  def deviceCredentials = java.deviceCredentials()
  def emailTemplates = java.emailTemplates()
  def eventStreams = java.eventStreams()
  def events = java.events()
  def flows = java.flows()
  def forms = java.forms()
  def userGrants = java.userGrants()
  def groups = java.groups()
  def hooks = java.hooks()
  def jobs = java.jobs()
  def logStreams = java.logStreams()
  def logs = java.logs()
  def networkAcls = java.networkAcls()
  def organizations = java.organizations()
  def prompts = java.prompts()
  def rateLimitPolicies = java.rateLimitPolicies()
  def refreshTokens = java.refreshTokens()
  def resourceServers = java.resourceServers()
  def roles = java.roles()
  def rules = java.rules()
  def rulesConfigs = java.rulesConfigs()
  def selfServiceProfiles = java.selfServiceProfiles()
  def sessions = java.sessions()
  def stats = java.stats()
  def supplementalSignals = java.supplementalSignals()
  def tickets = java.tickets()
  def tokenExchangeProfiles = java.tokenExchangeProfiles()
  def userAttributeProfiles = java.userAttributeProfiles()
  def userBlocks = java.userBlocks()
  def users = java.users()
  def anomaly = java.anomaly()
  def attackProtection = java.attackProtection()
  def emails = java.emails()
  def guardian = java.guardian()
  def keys = java.keys()
  def riskAssessments = java.riskAssessments()
  def tenants = java.tenants()
  def verifiableCredentials = java.verifiableCredentials()

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
