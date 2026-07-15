package auth0scala3

import com.auth0.client.auth.AuthAPI
import com.auth0.client.mgmt.ManagementApi
import com.auth0.net.Request
import com.auth0.net.Response
import java.util.concurrent.CompletableFuture
import scala.concurrent.ExecutionContext.Implicits.global

final class Auth0Suite extends munit.FunSuite:
  test("authentication URLs use Java SDK") {
    val auth = Auth0.auth("tenant.auth0.com", "client-id", "secret")
    val authorizationUrl = auth.authorizationUrl(
      "https://app.test/callback",
      config =>
        config.audience("https://api.test")
        config.scope("openid profile")
        config.state("state")
    )
    val logoutUrl = auth.logoutUrl(
      "https://app.test/out",
      true,
      _.federated()
    )

    assert(authorizationUrl.startsWith("https://tenant.auth0.com/authorize?"))
    assert(authorizationUrl.contains("client_id=client-id"))
    assert(authorizationUrl.contains("audience=https%3A%2F%2Fapi.test"))
    assert(authorizationUrl.contains("scope=openid%20profile"))
    assert(logoutUrl.startsWith("https://tenant.auth0.com/v2/logout?"))
    assert(logoutUrl.contains("federated="))
  }

  test("management resources are typed Java clients") {
    val management = Auth0.managementWithToken(
      "tenant.auth0.com",
      "token",
      config =>
        config.timeout(3)
        config.maxRetries(1)
        config.customDomain("login.test")
        config.header("x-test", "value")
    )

    assertEquals(management.users().getClass.getName, "com.auth0.client.mgmt.UsersClient")
    assertEquals(management.clients().getClass.getName, "com.auth0.client.mgmt.ClientsClient")
    assertEquals(management.rateLimitPolicies().getClass.getName, "com.auth0.client.mgmt.RateLimitPoliciesClient")
    assertEquals(management.users().roles().getClass.getName, "com.auth0.client.mgmt.users.RolesClient")
    assert(management.organizations() != null)
  }

  test("client credentials build Java client") {
    val management = Auth0.managementWithClientCredentials(
      "tenant.auth0.com",
      "client-id",
      "client-secret"
    )

    assert(management.users() != null)
  }

  test("Java SDK remains reachable") {
    val auth = Auth0.auth("tenant.auth0.com", "client-id")
    val management = Auth0.managementWithToken("tenant.auth0.com", "token")

    assert(auth.java.isInstanceOf[AuthAPI])
    assert(management.java.isInstanceOf[ManagementApi])
    assert(auth.use(_.authorizeUrl("https://app.test/callback")) != null)
    assert(management.use(_.users()) != null)
  }

  test("Java requests adapt to Scala") {
    val request = new StaticRequest("body")

    for
      response <- request.future
      body <- request.bodyFuture
    yield
      assertEquals(request.executeBody(), "body")
      assertEquals(body, "body")
      assertEquals(response.statusCode, 200)
      assertEquals(response.bodyOption, Some("body"))
  }

  test("existing Java clients can be wrapped") {
    val javaAuth = AuthAPI.newBuilder("tenant.auth0.com", "client-id").build()
    val javaManagement = ManagementApi.builder().domain("tenant.auth0.com").token("token").build()

    assertEquals(javaAuth.asScala.java, javaAuth)
    assertEquals(javaManagement.asScala.java, javaManagement)
  }

  private final class StaticRequest[A](value: A) extends Request[A]:
    override def execute(): Response[A] =
      new StaticResponse(value)

    override def executeAsync(): CompletableFuture[Response[A]] =
      CompletableFuture.completedFuture[Response[A]](new StaticResponse(value))

    override def addHeader(name: String, value: String): Request[A] =
      this

    override def addParameter(name: String, value: Object): Request[A] =
      this

    override def setBody(body: Object): Request[A] =
      this

  private final class StaticResponse[A](value: A) extends Response[A]:
    override def getHeaders(): java.util.Map[String, String] =
      java.util.Map.of("x-ok", "yes")

    override def getBody(): A =
      value

    override def getStatusCode(): Int =
      200
