package com.auth0.kotlin.webapp

import com.auth0.kotlin.Auth0
import com.fasterxml.jackson.databind.ObjectMapper
import com.sun.net.httpserver.HttpExchange
import com.sun.net.httpserver.HttpHandler
import com.sun.net.httpserver.HttpServer
import java.net.InetSocketAddress
import java.net.URLDecoder
import java.nio.charset.StandardCharsets.UTF_8
import java.nio.file.Files
import java.nio.file.Path
import java.security.SecureRandom
import java.util.Base64
import java.util.concurrent.ConcurrentHashMap
import kotlin.system.exitProcess

private const val SDK = "kotlin"
private const val PORT = 3000
private const val CALLBACK_PATH = "/login/oauth2/code/okta"

private val domain = env("AUTH0_DOMAIN")
private val clientId = env("AUTH0_CLIENT_ID")
private val clientSecret = env("AUTH0_CLIENT_SECRET")
private val redirectUri = "http://localhost:$PORT$CALLBACK_PATH"
private val returnTo = "http://localhost:$PORT/"

private val auth = Auth0.auth(domain, clientId, clientSecret)
private val mapper = ObjectMapper()
private val random = SecureRandom()
private val sessions = ConcurrentHashMap<String, Map<String, Any?>>()
private val states = ConcurrentHashMap<String, Boolean>()
private val assets: Path = Path.of("webapp", "dist")

fun main() {
    if (!Files.isDirectory(assets)) {
        System.err.println("missing ${assets.toAbsolutePath()}, run: cd webapp && npm install && npm run build")
        exitProcess(1)
    }

    val server = HttpServer.create(InetSocketAddress(PORT), 0)
    server.createContext("/", handler(::serveStatic))
    server.createContext("/login", handler(::startLogin))
    server.createContext(CALLBACK_PATH, handler(::completeLogin))
    server.createContext("/api/me", handler(::session))
    server.createContext("/logout", handler(::logout))
    server.start()
    println("$SDK webapp on http://localhost:$PORT")
}

private fun handler(action: (HttpExchange) -> Unit): HttpHandler =
    HttpHandler { exchange ->
        try {
            action(exchange)
        } catch (error: Throwable) {
            error.printStackTrace()
            send(exchange, 500, "text/plain", "error: ${error.message}".toByteArray())
        } finally {
            exchange.close()
        }
    }

private fun startLogin(exchange: HttpExchange) {
    val state = token()
    states[state] = true
    val url = auth.authorizationUrl(redirectUri) {
        withScope("openid profile email")
        withState(state)
    }
    redirect(exchange, url)
}

private fun completeLogin(exchange: HttpExchange) {
    val params = query(exchange)
    val code = params["code"]
    val state = params["state"]

    if (code == null || state == null || states.remove(state) == null) {
        redirect(exchange, "/")
        return
    }

    val tokens = auth.java.exchangeCode(code, redirectUri).execute().body
    val claims = decodeClaims(tokens.idToken)
    val sid = token()
    sessions[sid] = claims
    exchange.responseHeaders.add("Set-Cookie", "sid=$sid; Path=/; HttpOnly; SameSite=Lax")
    redirect(exchange, "/profile")
}

private fun session(exchange: HttpExchange) {
    val claims = cookie(exchange)?.let { sessions[it] }
    val body = mapOf(
        "authenticated" to (claims != null),
        "sdk" to SDK,
        "claims" to (claims ?: emptyMap()),
    )
    send(exchange, 200, "application/json", mapper.writeValueAsBytes(body))
}

private fun logout(exchange: HttpExchange) {
    cookie(exchange)?.let(sessions::remove)
    exchange.responseHeaders.add("Set-Cookie", "sid=; Path=/; HttpOnly; Max-Age=0")
    redirect(exchange, auth.logoutUrl(returnTo))
}

private fun serveStatic(exchange: HttpExchange) {
    val path = exchange.requestURI.path
    val requested = if (path == "/" || path == "/profile") "/index.html" else path
    val file = assets.resolve(requested.removePrefix("/")).normalize()

    if (!file.startsWith(assets) || !Files.isRegularFile(file)) {
        send(exchange, 404, "text/plain", "not found".toByteArray())
    } else {
        send(exchange, 200, contentType(file), Files.readAllBytes(file))
    }
}

private fun decodeClaims(idToken: String): Map<String, Any?> {
    val payload = idToken.split(".")[1]
    val json = String(Base64.getUrlDecoder().decode(payload), UTF_8)
    @Suppress("UNCHECKED_CAST")
    return mapper.readValue(json, Map::class.java) as Map<String, Any?>
}

private fun contentType(file: Path): String {
    val name = file.fileName.toString()
    return when {
        name.endsWith(".html") -> "text/html"
        name.endsWith(".js") -> "text/javascript"
        name.endsWith(".css") -> "text/css"
        name.endsWith(".svg") -> "image/svg+xml"
        else -> "application/octet-stream"
    }
}

private fun query(exchange: HttpExchange): Map<String, String> =
    (exchange.requestURI.query ?: "")
        .split("&")
        .mapNotNull { pair ->
            val parts = pair.split("=", limit = 2)
            if (parts.size == 2) parts[0] to URLDecoder.decode(parts[1], UTF_8) else null
        }
        .toMap()

private fun cookie(exchange: HttpExchange): String? =
    exchange.requestHeaders.getFirst("Cookie")
        ?.split(";")
        ?.map(String::trim)
        ?.firstOrNull { it.startsWith("sid=") }
        ?.removePrefix("sid=")

private fun redirect(exchange: HttpExchange, url: String) {
    exchange.responseHeaders.add("Location", url)
    exchange.sendResponseHeaders(302, -1)
}

private fun send(exchange: HttpExchange, status: Int, contentType: String, body: ByteArray) {
    exchange.responseHeaders.add("Content-Type", contentType)
    exchange.sendResponseHeaders(status, body.size.toLong())
    exchange.responseBody.write(body)
}

private fun token(): String {
    val bytes = ByteArray(32)
    random.nextBytes(bytes)
    return Base64.getUrlEncoder().withoutPadding().encodeToString(bytes)
}

private fun env(name: String): String =
    System.getenv(name) ?: error("$name must be set, source ./auth0-env.sh")
