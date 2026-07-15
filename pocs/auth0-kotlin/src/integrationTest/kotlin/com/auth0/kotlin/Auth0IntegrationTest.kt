package com.auth0.kotlin

import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertTrue
import okhttp3.OkHttpClient
import okhttp3.Request as HttpRequest

class Auth0IntegrationTest {
    private val domain: String = requireNotNull(System.getenv("AUTH0_DOMAIN"))
    private val clientId: String = requireNotNull(System.getenv("AUTH0_CLIENT_ID"))
    private val clientSecret: String = requireNotNull(System.getenv("AUTH0_CLIENT_SECRET"))
    private val redirectUri: String = requireNotNull(System.getenv("AUTH0_REDIRECT_URI"))

    @Test
    fun liveTenantAcceptsAuthorizationRequestBuiltBySdk() {
        val auth = Auth0.auth(domain, clientId, clientSecret)
        val url = auth.authorizationUrl(redirectUri) {
            withScope("openid profile email")
            withState("integration-state")
        }

        assertTrue(url.startsWith("https://$domain/authorize"))
        assertTrue(url.contains("client_id=$clientId"))

        val client = OkHttpClient.Builder().followRedirects(false).build()
        val response = client.newCall(HttpRequest.Builder().url(url).build()).execute()

        response.use {
            val location = it.header("location").orEmpty()
            assertEquals(302, it.code, "tenant rejected the authorization request: $location")
            assertTrue(
                location.startsWith("/u/login"),
                "expected redirect to Universal Login, got $location",
            )
        }
    }
}
