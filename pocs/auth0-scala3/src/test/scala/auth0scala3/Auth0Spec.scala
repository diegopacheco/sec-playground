package auth0scala3

import com.auth0.client.auth.AuthAPI
import com.auth0.client.mgmt.ManagementApi
import com.auth0.net.Request
import com.auth0.net.Response
import java.util.concurrent.CompletableFuture
import scala.concurrent.Await
import scala.concurrent.ExecutionContext.Implicits.global
import scala.concurrent.duration.DurationInt

object Auth0Spec:
  def main(args: Array[String]): Unit =
    val tests = Vector(
      "authentication URLs use Java SDK" -> (() => authenticationUrlsUseJavaSdk()),
      "management resources are typed Java clients" -> (() => managementResourcesAreTypedJavaClients()),
      "client credentials build Java client" -> (() => clientCredentialsBuildJavaClient()),
      "Java SDK remains reachable" -> (() => javaSdkRemainsReachable()),
      "Java requests adapt to Scala" -> (() => javaRequestsAdaptToScala()),
      "existing Java clients can be wrapped" -> (() => existingJavaClientsCanBeWrapped())
    )
    tests.foreach { case (name, test) =>
      test()
      println(name + " passed")
    }

  private def assertEquals[A](actual: A, expected: A): Unit =
    if actual != expected then throw AssertionError("expected " + expected + " got " + actual)

  private def assertTrue(value: Boolean): Unit =
    if !value then throw AssertionError("assertion failed")

  private def authenticationUrlsUseJavaSdk(): Unit =
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

    assertTrue(authorizationUrl.startsWith("https://tenant.auth0.com/authorize?"))
    assertTrue(authorizationUrl.contains("client_id=client-id"))
    assertTrue(authorizationUrl.contains("audience=https%3A%2F%2Fapi.test"))
    assertTrue(authorizationUrl.contains("scope=openid%20profile"))
    assertTrue(logoutUrl.startsWith("https://tenant.auth0.com/v2/logout?"))
    assertTrue(logoutUrl.contains("federated="))

  private def managementResourcesAreTypedJavaClients(): Unit =
    val management = Auth0.managementWithToken(
      "tenant.auth0.com",
      "token",
      config =>
        config.timeout(3)
        config.maxRetries(1)
        config.customDomain("login.test")
        config.header("x-test", "value")
    )

    assertEquals(management.users.getClass.getName, "com.auth0.client.mgmt.UsersClient")
    assertEquals(management.clients.getClass.getName, "com.auth0.client.mgmt.ClientsClient")
    assertEquals(management.rateLimitPolicies.getClass.getName, "com.auth0.client.mgmt.RateLimitPoliciesClient")
    assertEquals(management.users.roles().getClass.getName, "com.auth0.client.mgmt.users.RolesClient")
    assertTrue(management.organizations != null)

  private def clientCredentialsBuildJavaClient(): Unit =
    val management = Auth0.managementWithClientCredentials(
      "tenant.auth0.com",
      "client-id",
      "client-secret"
    )

    assertTrue(management.users != null)

  private def javaSdkRemainsReachable(): Unit =
    val auth = Auth0.auth("tenant.auth0.com", "client-id")
    val management = Auth0.managementWithToken("tenant.auth0.com", "token")

    assertTrue(auth.java.isInstanceOf[AuthAPI])
    assertTrue(management.java.isInstanceOf[ManagementApi])
    assertTrue(auth.use(_.authorizeUrl("https://app.test/callback")) != null)
    assertTrue(management.use(_.users()) != null)

  private def javaRequestsAdaptToScala(): Unit =
    val request = StaticRequest("body")
    val response = Await.result(request.future, 1.second)

    assertEquals(request.executeBody(), "body")
    assertEquals(Await.result(request.bodyFuture, 1.second), "body")
    assertEquals(response.statusCode, 200)
    assertEquals(response.bodyOption, Some("body"))

  private def existingJavaClientsCanBeWrapped(): Unit =
    val javaAuth = AuthAPI.newBuilder("tenant.auth0.com", "client-id").build()
    val javaManagement = ManagementApi.builder().domain("tenant.auth0.com").token("token").build()

    assertEquals(javaAuth.asScala.java, javaAuth)
    assertEquals(javaManagement.asScala.java, javaManagement)

  private final class StaticRequest[A](value: A) extends Request[A]:
    override def execute(): Response[A] =
      StaticResponse(value)

    override def executeAsync(): CompletableFuture[Response[A]] =
      CompletableFuture.completedFuture[Response[A]](StaticResponse(value))

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

  private object StaticRequest:
    def apply[A](value: A): StaticRequest[A] =
      new StaticRequest(value)

  private object StaticResponse:
    def apply[A](value: A): StaticResponse[A] =
      new StaticResponse(value)
