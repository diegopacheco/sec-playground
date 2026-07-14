package com.auth0.kotlin

import com.auth0.net.Request
import com.auth0.net.Response
import java.io.File
import java.lang.reflect.Modifier
import java.util.concurrent.CompletableFuture
import kotlinx.coroutines.runBlocking
import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertNotNull
import kotlin.test.assertTrue

class Auth0Test {
    @Test
    fun buildsAuthClientAndUrls() {
        val auth = Auth0.auth("tenant.auth0.com", "client-id", "secret")

        val authorizeUrl = auth.authorizeUrlFor("https://app.test/callback") {
            withAudience("https://api.test")
            withScope("openid profile")
            withState("state")
        }

        val logoutUrl = auth.logoutUrlFor("https://app.test/out", federated = true)

        assertTrue(authorizeUrl.startsWith("https://tenant.auth0.com/authorize?"))
        assertTrue(authorizeUrl.contains("client_id=client-id"))
        assertTrue(authorizeUrl.contains("audience=https%3A%2F%2Fapi.test"))
        assertTrue(authorizeUrl.contains("scope=openid%20profile"))
        assertTrue(logoutUrl.startsWith("https://tenant.auth0.com/v2/logout?"))
        assertTrue(logoutUrl.contains("federated="))
    }

    @Test
    fun buildsManagementClient() {
        val management = Auth0.managementWithToken("tenant.auth0.com", "token") {
            timeout = 3
            maxRetries = 1
            customDomain = "login.test"
            header("x-test", "value")
        }

        assertNotNull(management.users())
        assertNotNull(management.clients())
        assertNotNull(management.organizations())
    }

    @Test
    fun exposesKotlinManagementResources() {
        val javaNames = com.auth0.client.mgmt.ManagementApi::class.java.methods
            .filter { it.parameterCount == 0 }
            .filter { !Modifier.isStatic(it.modifiers) }
            .filter { it.declaringClass != Object::class.java }
            .map { it.name }
            .toSet() - "builder"

        val kotlinNames = KotlinManagementApi::class.java.methods
            .filter { it.parameterCount <= 1 }
            .map { it.name }
            .toSet()

        assertTrue(javaNames.isNotEmpty())
        assertEquals(emptySet(), javaNames - kotlinNames)
    }

    @Test
    fun reachesNestedManagementResources() {
        val management = Auth0.kotlinManagementWithToken("tenant.auth0.com", "token")

        val usersRoles = management.users().resource("roles")
        val jobsRaw = management.jobs().raw()

        assertEquals("com.auth0.client.mgmt.users.RolesClient", usersRoles.java.javaClass.name)
        assertEquals("com.auth0.client.mgmt.RawJobsClient", jobsRaw.java.javaClass.name)
    }

    @Test
    fun mapsRequestResults() {
        val request = StaticRequest("body")
        val response = request.response()
        val asyncBody = request.bodyAsync().get()
        val (status, headers, body) = request.execute()

        assertEquals(200, response.statusCode)
        assertEquals("yes", response.headers["x-ok"])
        assertEquals("body", response.body)
        assertEquals("body", asyncBody)
        assertEquals(200, status)
        assertEquals("yes", headers["x-ok"])
        assertEquals("body", body)
    }

    @Test
    fun mapsCoroutineRequestResults() = runBlocking {
        val request = StaticRequest("body")
        val response = request.awaitResponse()
        val body = request.awaitBody()

        assertEquals(200, response.statusCode)
        assertEquals("yes", response.headers["x-ok"])
        assertEquals("body", response.body)
        assertEquals("body", body)
    }

    @Test
    fun mapsNullableRequestResults() = runBlocking {
        val request = StaticRequest<String?>(null)
        val response = request.awaitNullableResponse()

        assertEquals(200, response.statusCode)
        assertEquals(null, response.body)
        assertEquals(null, request.nullableBody())
    }

    @Test
    fun serializesKotlinDataClasses() {
        val json = Auth0Json.encode(JsonUser("user-id", "email@test"))
        val user = Auth0Json.decode<JsonUser>(json)

        assertEquals(JsonUser("user-id", "email@test"), user)
    }

    @Test
    fun tracksJavaTestCoverage() {
        val tests = File("/Users/diegopacheco/git/misc/auth0-java/src/test/java")
            .walkTopDown()
            .filter { it.isFile && it.name.endsWith("Test.java") }
            .toList()

        val testMethods = tests.sumOf { file ->
            Regex("@Test").findAll(file.readText()).count()
        }

        assertTrue(tests.size >= 130)
        assertTrue(testMethods >= 500)
    }
}

private data class JsonUser(val id: String, val email: String)

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
