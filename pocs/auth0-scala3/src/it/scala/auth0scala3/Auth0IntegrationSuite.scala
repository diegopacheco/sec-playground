package auth0scala3

import com.auth0.exception.APIException
import com.auth0.client.mgmt.core.OAuthTokenException

final class Auth0IntegrationSuite extends munit.FunSuite:
  private val domain = sys.env.getOrElse("AUTH0_DOMAIN", fail("AUTH0_DOMAIN must be set"))

  test("reaches live tenant through authorization API") {
    val auth = Auth0.auth(domain, "integration-invalid-client")
    val error = intercept[APIException] {
      auth.java.userInfo("integration-invalid-token").execute()
    }

    assertEquals(error.getStatusCode, 401)
  }

  test("rejects invalid authentication client credentials") {
    val auth = Auth0.auth(
      domain,
      "integration-invalid-client",
      "integration-invalid-secret"
    )
    val error = intercept[APIException] {
      auth.java.requestToken(s"https://$domain/api/v2/").execute()
    }

    assertEquals(error.getStatusCode, 401)
  }

  test("rejects invalid management client credentials") {
    val management = Auth0.managementWithClientCredentials(
      domain,
      "integration-invalid-client",
      "integration-invalid-secret"
    )
    val error = intercept[OAuthTokenException] {
      management.users().get("auth0|integration-invalid-user")
    }

    assert(error.getMessage.contains("HTTP 401"))
  }
