package com.auth0.kotlin

import com.auth0.exception.APIException
import com.auth0.client.mgmt.core.OAuthTokenException
import kotlin.test.Test
import kotlin.test.assertContains
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith

class Auth0IntegrationTest {
    private val domain: String = requireNotNull(System.getenv("AUTH0_DOMAIN"))

    @Test
    fun reachesLiveTenantThroughAuthorizationApi() {
        val auth = Auth0.auth(domain, "integration-invalid-client")
        val error = assertFailsWith<APIException> {
            auth.java.userInfo("integration-invalid-token").execute()
        }

        assertEquals(401, error.statusCode)
    }

    @Test
    fun rejectsInvalidAuthenticationClientCredentials() {
        val auth = Auth0.auth(
            domain,
            "integration-invalid-client",
            "integration-invalid-secret",
        )
        val error = assertFailsWith<APIException> {
            auth.java.requestToken("https://$domain/api/v2/").execute()
        }

        assertEquals(401, error.statusCode)
    }

    @Test
    fun rejectsInvalidManagementClientCredentials() {
        val management = Auth0.managementWithClientCredentials(
            domain,
            "integration-invalid-client",
            "integration-invalid-secret",
        )
        val error = assertFailsWith<OAuthTokenException> {
            management.users.get("auth0|integration-invalid-user")
        }

        assertContains(error.message.orEmpty(), "HTTP 401")
    }
}
