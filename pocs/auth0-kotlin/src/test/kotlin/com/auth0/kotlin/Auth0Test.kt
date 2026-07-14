package com.auth0.kotlin

import com.auth0.client.auth.AuthAPI
import com.auth0.client.mgmt.ManagementApi
import com.auth0.net.Request
import com.auth0.net.Response
import java.util.concurrent.CompletableFuture
import kotlinx.coroutines.runBlocking
import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertNotNull
import kotlin.test.assertTrue

class Auth0Test {
    @Test
    fun buildsAuthenticationUrlsWithJavaSdk() {
        val auth = Auth0.auth("tenant.auth0.com", "client-id", "secret")

        val authorizationUrl = auth.authorizationUrl("https://app.test/callback") {
            withAudience("https://api.test")
            withScope("openid profile")
            withState("state")
        }
        val logoutUrl = auth.logoutUrl("https://app.test/out") {
            useFederated(true)
        }

        assertTrue(authorizationUrl.startsWith("https://tenant.auth0.com/authorize?"))
        assertTrue(authorizationUrl.contains("client_id=client-id"))
        assertTrue(authorizationUrl.contains("audience=https%3A%2F%2Fapi.test"))
        assertTrue(authorizationUrl.contains("scope=openid%20profile"))
        assertTrue(logoutUrl.startsWith("https://tenant.auth0.com/v2/logout?"))
        assertTrue(logoutUrl.contains("federated="))
    }

    @Test
    fun exposesOfficialJavaClientsAndModels() {
        val management = Auth0.managementWithToken("tenant.auth0.com", "token") {
            timeout(3)
            maxRetries(1)
            customDomain("login.test")
            header("x-test", "value")
        }

        assertEquals("com.auth0.client.mgmt.UsersClient", management.users.javaClass.name)
        assertEquals("com.auth0.client.mgmt.ClientsClient", management.clients.javaClass.name)
        assertEquals("com.auth0.client.mgmt.RateLimitPoliciesClient", management.rateLimitPolicies.javaClass.name)
        assertEquals("com.auth0.client.mgmt.users.RolesClient", management.users.roles().javaClass.name)
        assertNotNull(management.organizations)
    }

    @Test
    fun buildsAutomaticClientCredentialsClient() {
        val management = Auth0.managementWithClientCredentials(
            "tenant.auth0.com",
            "client-id",
            "client-secret"
        )

        assertNotNull(management.users)
    }

    @Test
    fun keepsJavaSdkReachable() {
        val auth = Auth0.auth("tenant.auth0.com", "client-id")
        val management = Auth0.managementWithToken("tenant.auth0.com", "token")

        assertEquals(AuthAPI::class.java, auth.java.javaClass)
        assertEquals(ManagementApi::class.java, management.java.javaClass)
        assertNotNull(auth.use { authorizeUrl("https://app.test/callback") })
        assertNotNull(management.use { users() })
    }

    @Test
    fun adaptsJavaRequestsToKotlin() = runBlocking {
        val request = StaticRequest("body")
        val response = request.await()

        assertEquals("body", request.executeBody())
        assertEquals("body", request.awaitBody())
        assertEquals(200, response.statusCode)
        assertEquals("yes", response.headers["x-ok"])
    }

    @Test
    fun wrapsExistingJavaClients() {
        val javaAuth = AuthAPI.newBuilder("tenant.auth0.com", "client-id").build()
        val javaManagement = ManagementApi.builder().domain("tenant.auth0.com").token("token").build()

        assertEquals(javaAuth, javaAuth.asKotlin().java)
        assertEquals(javaManagement, javaManagement.asKotlin().java)
    }
}

private class StaticRequest<T>(private val value: T) : Request<T> {
    override fun execute(): Response<T> =
        StaticResponse(value)

    override fun executeAsync(): CompletableFuture<Response<T>> =
        CompletableFuture.completedFuture(StaticResponse(value))

    override fun addHeader(name: String, value: String): Request<T> =
        this

    override fun addParameter(name: String, value: Any): Request<T> =
        this

    override fun setBody(body: Any): Request<T> =
        this
}

private class StaticResponse<T>(private val value: T) : Response<T> {
    override fun getHeaders(): Map<String, String> =
        mapOf("x-ok" to "yes")

    override fun getBody(): T =
        value

    override fun getStatusCode(): Int =
        200
}
