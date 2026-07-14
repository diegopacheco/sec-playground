package auth0scala3

import java.security.KeyPairGenerator
import java.util.Base64
import scala.collection.mutable.ArrayBuffer
import scala.concurrent.Await
import scala.concurrent.ExecutionContext.Implicits.global
import scala.concurrent.duration.DurationInt

object Auth0Spec:
  def main(args: Array[String]): Unit =
    val tests = Vector(
      "authorize url" -> (() => testAuthorizeUrl()),
      "logout url" -> (() => testLogoutUrl()),
      "token request" -> (() => testTokenRequest()),
      "code verifier" -> (() => testCodeVerifier()),
      "mfa oob validation" -> (() => testMfaOobValidation()),
      "management users request" -> (() => testManagementUsersRequest()),
      "nested management resource" -> (() => testNestedManagementResource()),
      "json roundtrip" -> (() => testJsonRoundtrip()),
      "raw response" -> (() => testRawResponse()),
      "endpoint errors" -> (() => testEndpointErrors()),
      "async request" -> (() => testAsyncRequest()),
      "request options" -> (() => testRequestOptions()),
      "custom domain whitelist" -> (() => testCustomDomainWhitelist()),
      "rsa assertion signer" -> (() => testRsaAssertionSigner()),
      "multipart upload" -> (() => testMultipartUpload()),
      "sse parser" -> (() => testSseParser()),
      "nested generated paths" -> (() => testNestedGeneratedPaths())
    )
    tests.foreach { case (name, test) =>
      test()
      println(name + " passed")
    }

  final class RecordingTransport(response: ApiResponse = ApiResponse(200, Map.empty, """{"access_token":"abc"}""")) extends Transport:
    val requests = ArrayBuffer.empty[ApiRequest]
    def send(request: ApiRequest): ApiResponse =
      requests += request
      response

  private def assertEquals[A](actual: A, expected: A): Unit =
    if actual != expected then throw AssertionError("expected " + expected + " got " + actual)

  private def assertTrue(value: Boolean): Unit =
    if !value then throw AssertionError("assertion failed")

  private def testAuthorizeUrl(): Unit =
    val auth = Auth0.auth("tenant.auth0.com", "cid")
    val url = auth.authorizeUrl("https://app/callback")
      .withConnection("github")
      .withAudience("https://api")
      .withScope("openid profile")
      .withState("state 1")
      .withCodeChallenge("challenge")
      .build
    assertTrue(url.startsWith("https://tenant.auth0.com/authorize?"))
    assertTrue(url.contains("client_id=cid"))
    assertTrue(url.contains("redirect_uri=https%3A%2F%2Fapp%2Fcallback"))
    assertTrue(url.contains("connection=github"))
    assertTrue(url.contains("scope=openid%20profile"))
    assertTrue(url.contains("code_challenge_method=S256"))

  private def testLogoutUrl(): Unit =
    val auth = Auth0.auth("https://tenant.auth0.com", "cid")
    val url = auth.logoutUrl("https://app/home", setClientId = true).useFederated(true).build
    assertTrue(url.startsWith("https://tenant.auth0.com/v2/logout?"))
    assertTrue(url.contains("returnTo=https%3A%2F%2Fapp%2Fhome"))
    assertTrue(url.contains("client_id=cid"))
    assertTrue(url.contains("federated="))

  private def testTokenRequest(): Unit =
    val transport = RecordingTransport()
    val auth = Auth0.auth("tenant.auth0.com", "cid", Some("secret"), transport)
    val request = auth.requestToken("https://tenant.auth0.com/api/v2/").request
    assertEquals(request.method, Method.POST)
    assertEquals(request.url, "https://tenant.auth0.com/oauth/token")
    val values = request.body.collect { case RequestBody.Form(v) => v }.get
    assertEquals(values("grant_type"), "client_credentials")
    assertEquals(values("client_id"), "cid")
    assertEquals(values("client_secret"), "secret")
    assertEquals(values("audience"), "https://tenant.auth0.com/api/v2/")

  private def testCodeVerifier(): Unit =
    val auth = Auth0.auth("tenant.auth0.com", "cid")
    val values = auth.exchangeCodeWithVerifier("code", "verifier", "https://app/callback").request.body.collect { case RequestBody.Form(v) => v }.get
    assertEquals(values("grant_type"), "authorization_code")
    assertEquals(values("code_verifier"), "verifier")
    assertTrue(!values.contains("client_secret"))

  private def testMfaOobValidation(): Unit =
    val auth = Auth0.auth("tenant.auth0.com", "cid")
    try
      auth.addOobAuthenticator("token", Vector("sms"))
      throw AssertionError("validation failed")
    catch
      case _: IllegalArgumentException => ()

  private def testManagementUsersRequest(): Unit =
    val transport = RecordingTransport()
    val api = ManagementApi.withToken("tenant.auth0.com", "token", transport)
    api.users.get(Map("page" -> Vector("0"), "per_page" -> Vector("50"))).send()
    val request = transport.requests.last
    assertEquals(request.method, Method.GET)
    assertEquals(request.url, "https://tenant.auth0.com/api/v2/users")
    assertEquals(request.headers("Authorization"), "Bearer token")
    assertEquals(request.query("page"), Vector("0"))
    assertEquals(request.query("per_page"), Vector("50"))

  private def testNestedManagementResource(): Unit =
    val transport = RecordingTransport()
    val api = ManagementApi.withToken("tenant.auth0.com", "token", transport)
    val request = api.users.at("auth0|1", "roles").post(Json.obj("roles" -> Json.arr(Json.str("rol_1")))).request
    assertEquals(request.url, "https://tenant.auth0.com/api/v2/users/auth0%7C1/roles")
    assertEquals(request.method, Method.POST)
    assertEquals(request.body.map(_.contentType), Some("application/json"))

  private def testJsonRoundtrip(): Unit =
    val json = Json.obj("name" -> Json.str("Ada"), "active" -> Json.bool(true), "roles" -> Json.arr(Json.str("admin"), Json.num(2)))
    assertEquals(Json.parse(Json.stringify(json)), json)

  private def testRawResponse(): Unit =
    val transport = RecordingTransport(ApiResponse(200, Map("x-rate-limit" -> Vector("10")), """{"ok":true}"""))
    val api = ManagementApi.withToken("tenant.auth0.com", "token", transport)
    val response = api.users.get().raw[Json].send()
    assertEquals(response.status, 200)
    assertEquals(response.headers("x-rate-limit"), Vector("10"))
    assertEquals(response.body, Json.obj("ok" -> Json.bool(true)))

  private def testEndpointErrors(): Unit =
    val transport = RecordingTransport(ApiResponse(401, Map.empty, """{"error":"bad"}"""))
    val api = ManagementApi.withToken("tenant.auth0.com", "token", transport)
    try
      api.users.get().send()
      throw AssertionError("error mapping failed")
    catch
      case e: UnauthorizedError =>
        assertEquals(e.status, 401)

  private def testAsyncRequest(): Unit =
    val transport = RecordingTransport(ApiResponse(200, Map.empty, """{"ok":true}"""))
    val api = ManagementApi.withToken("tenant.auth0.com", "token", transport)
    val response = Await.result(api.users.get().sendAsync(), 1.second)
    assertEquals(response.status, 200)

  private def testRequestOptions(): Unit =
    val transport = RecordingTransport()
    val api = ManagementApi.withToken("tenant.auth0.com", "token", transport)
    val options = RequestOptions(Map("X-Test" -> "1"), Map("q" -> Vector("email:*")), Some(1.second))
    val request = api.users.get(options = options).request
    assertEquals(request.headers("X-Test"), "1")
    assertEquals(request.query("q"), Vector("email:*"))
    assertEquals(request.timeout, Some(1.second))

  private def testCustomDomainWhitelist(): Unit =
    val api = ManagementApi(ManagementConfig(domain = Some("tenant.auth0.com"), token = Some("token"), customDomain = Some("login.acme.com"), transport = RecordingTransport()))
    val whitelisted = api.users.at("auth0|1").patch(Json.obj()).request
    val blocked = api.logs.get().request
    assertEquals(whitelisted.headers("Auth0-Custom-Domain"), "login.acme.com")
    assertTrue(!blocked.headers.contains("Auth0-Custom-Domain"))

  private def testRsaAssertionSigner(): Unit =
    val generator = KeyPairGenerator.getInstance("RSA")
    generator.initialize(2048)
    val key = generator.generateKeyPair().getPrivate.asInstanceOf[java.security.interfaces.RSAPrivateKey]
    val signer = RSAClientAssertionSigner(key)
    val assertion = signer.createSignedClientAssertion("cid", "https://tenant.auth0.com", "cid")
    val parts = assertion.split("\\.")
    assertEquals(parts.length, 3)
    val header = String(Base64.getUrlDecoder.decode(parts(0)), java.nio.charset.StandardCharsets.UTF_8)
    assertEquals(Json.parse(header), Json.obj("alg" -> Json.str("RS256"), "typ" -> Json.str("JWT")))

  private def testMultipartUpload(): Unit =
    val api = ManagementApi.withToken("tenant.auth0.com", "token", RecordingTransport())
    val request = api.jobs.usersImports.upload("users", "users.json", "application/json", """[]""".getBytes(java.nio.charset.StandardCharsets.UTF_8), Map("connection_id" -> "con")).request
    assertEquals(request.method, Method.POST)
    assertTrue(request.body.exists(_.contentType.startsWith("multipart/form-data")))
    assertTrue(request.body.exists(_.text.contains("users.json")))

  private def testSseParser(): Unit =
    val events = Sse.parse("id: 1\nevent: created\ndata: a\ndata: b\nretry: 10\n\n")
    assertEquals(events.size, 1)
    assertEquals(events.head.id, Some("1"))
    assertEquals(events.head.event, Some("created"))
    assertEquals(events.head.data, "a\nb")
    assertEquals(events.head.retry, Some(10))

  private def testNestedGeneratedPaths(): Unit =
    val api = ManagementApi.withToken("tenant.auth0.com", "token", RecordingTransport())
    assertEquals(api.users.roles("auth0|1").path, "users/auth0%7C1/roles")
    assertEquals(api.organizations.members("org_1").roles("auth0|1").path, "organizations/org_1/members/auth0%7C1/roles")
    assertEquals(api.guardian.factors.sms.path, "guardian/factors/sms")
    assertEquals(api.jobs.usersExports.path, "jobs/users-exports")
