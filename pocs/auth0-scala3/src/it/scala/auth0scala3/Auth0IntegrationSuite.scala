package auth0scala3

import com.auth0.exception.APIException
import okhttp3.OkHttpClient
import okhttp3.Request as HttpRequest

final class Auth0IntegrationSuite extends munit.FunSuite:
  private val domain = sys.env("AUTH0_DOMAIN")
  private val clientId = sys.env("AUTH0_CLIENT_ID")
  private val clientSecret = sys.env("AUTH0_CLIENT_SECRET")
  private val redirectUri = sys.env("AUTH0_REDIRECT_URI")

  test("live tenant accepts the credentials and the authorization request") {
    val auth = Auth0.auth(domain, clientId, clientSecret)
    val url = auth.authorizationUrl(
      redirectUri,
      config =>
        config.scope("openid profile email")
        config.state("integration-state")
    )

    val client = OkHttpClient.Builder().followRedirects(false).build()
    val response = client.newCall(HttpRequest.Builder().url(url).build()).execute()

    try
      val location = Option(response.header("location")).getOrElse("")
      assertEquals(response.code, 302, s"tenant rejected the authorization request: $location")
      assert(
        location.startsWith("/u/login"),
        s"expected redirect to Universal Login, got $location"
      )
    finally response.close()

    val error = intercept[APIException] {
      auth.java.exchangeCode("integration-invalid-code", redirectUri).execute()
    }

    assertNotEquals(
      error.getStatusCode,
      401,
      s"Auth0 refused the client secret, so the credentials are not valid: ${error.getDescription}"
    )
    assertEquals(
      error.getError,
      "invalid_grant",
      s"expected Auth0 to authenticate the client and reject only the code, got ${error.getDescription}"
    )
  }
