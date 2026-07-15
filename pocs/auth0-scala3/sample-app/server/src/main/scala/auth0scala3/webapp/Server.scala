package auth0scala3.webapp

import auth0scala3.Auth0
import com.fasterxml.jackson.databind.ObjectMapper
import com.sun.net.httpserver.HttpExchange
import com.sun.net.httpserver.HttpHandler
import com.sun.net.httpserver.HttpServer
import java.net.InetSocketAddress
import java.nio.charset.StandardCharsets.UTF_8
import java.nio.file.Files
import java.nio.file.Path
import java.security.SecureRandom
import java.util.Base64
import scala.collection.concurrent.TrieMap
import scala.jdk.CollectionConverters.*

object Server:
  private val sdk = "scala 3"
  private val port = 3000
  private val callbackPath = "/login/oauth2/code/okta"

  private val domain = env("AUTH0_DOMAIN")
  private val clientId = env("AUTH0_CLIENT_ID")
  private val clientSecret = env("AUTH0_CLIENT_SECRET")
  private val redirectUri = s"http://localhost:$port$callbackPath"
  private val returnTo = s"http://localhost:$port/"

  private val auth = Auth0.auth(domain, clientId, clientSecret)
  private val mapper = ObjectMapper()
  private val random = SecureRandom()
  private val sessions = TrieMap.empty[String, Map[String, Object]]
  private val states = TrieMap.empty[String, Boolean]
  private val assets = Path.of("webapp", "dist")

  def main(args: Array[String]): Unit =
    if !Files.isDirectory(assets) then
      System.err.println(s"missing ${assets.toAbsolutePath}, run: cd webapp && npm install && npm run build")
      System.exit(1)

    val server = HttpServer.create(InetSocketAddress(port), 0)
    server.createContext("/", handler(serveStatic))
    server.createContext("/login", handler(startLogin))
    server.createContext(callbackPath, handler(completeLogin))
    server.createContext("/api/me", handler(session))
    server.createContext("/logout", handler(logout))
    server.start()
    println(s"$sdk webapp on http://localhost:$port")

  private def handler(action: HttpExchange => Unit): HttpHandler =
    exchange =>
      try action(exchange)
      catch
        case error: Throwable =>
          error.printStackTrace()
          send(exchange, 500, "text/plain", s"error: ${error.getMessage}".getBytes(UTF_8))
      finally exchange.close()

  private def startLogin(exchange: HttpExchange): Unit =
    val state = token()
    states.put(state, true)
    val url = auth.authorizationUrl(
      redirectUri,
      config =>
        config.scope("openid profile email")
        config.state(state)
    )
    redirect(exchange, url)

  private def completeLogin(exchange: HttpExchange): Unit =
    val params = query(exchange)
    val code = params.get("code")
    val state = params.get("state")

    if code.isEmpty || state.isEmpty || states.remove(state.get).isEmpty then
      redirect(exchange, "/")
    else
      val tokens = auth.java.exchangeCode(code.get, redirectUri).execute().getBody
      val claims = decodeClaims(tokens.getIdToken)
      val sid = token()
      sessions.put(sid, claims)
      exchange.getResponseHeaders.add("Set-Cookie", s"sid=$sid; Path=/; HttpOnly; SameSite=Lax")
      redirect(exchange, "/profile")

  private def session(exchange: HttpExchange): Unit =
    val claims = cookie(exchange).flatMap(sessions.get)
    val body = Map[String, Object](
      "authenticated" -> Boolean.box(claims.isDefined),
      "sdk" -> sdk,
      "claims" -> claims.getOrElse(Map.empty).asJava
    )
    send(exchange, 200, "application/json", mapper.writeValueAsBytes(body.asJava))

  private def logout(exchange: HttpExchange): Unit =
    cookie(exchange).foreach(sessions.remove)
    exchange.getResponseHeaders.add("Set-Cookie", "sid=; Path=/; HttpOnly; Max-Age=0")
    redirect(exchange, auth.logoutUrl(returnTo))

  private def serveStatic(exchange: HttpExchange): Unit =
    val path = exchange.getRequestURI.getPath
    val requested = if path == "/" || path == "/profile" then "/index.html" else path
    val file = assets.resolve(requested.stripPrefix("/")).normalize()

    if !file.startsWith(assets) || !Files.isRegularFile(file) then
      send(exchange, 404, "text/plain", "not found".getBytes(UTF_8))
    else
      send(exchange, 200, contentType(file), Files.readAllBytes(file))

  private def decodeClaims(idToken: String): Map[String, Object] =
    val payload = idToken.split("\\.")(1)
    val json = String(Base64.getUrlDecoder.decode(payload), UTF_8)
    mapper
      .readValue(json, classOf[java.util.Map[String, Object]])
      .asScala
      .toMap

  private def contentType(file: Path): String =
    val name = file.getFileName.toString
    if name.endsWith(".html") then "text/html"
    else if name.endsWith(".js") then "text/javascript"
    else if name.endsWith(".css") then "text/css"
    else if name.endsWith(".svg") then "image/svg+xml"
    else "application/octet-stream"

  private def query(exchange: HttpExchange): Map[String, String] =
    Option(exchange.getRequestURI.getQuery).getOrElse("").split("&").flatMap { pair =>
      pair.split("=", 2) match
        case Array(key, value) => Some(key -> java.net.URLDecoder.decode(value, UTF_8))
        case _                 => None
    }.toMap

  private def cookie(exchange: HttpExchange): Option[String] =
    Option(exchange.getRequestHeaders.getFirst("Cookie"))
      .flatMap(_.split(";").map(_.trim).find(_.startsWith("sid=")))
      .map(_.drop(4))

  private def redirect(exchange: HttpExchange, url: String): Unit =
    exchange.getResponseHeaders.add("Location", url)
    exchange.sendResponseHeaders(302, -1)

  private def send(exchange: HttpExchange, status: Int, contentType: String, body: Array[Byte]): Unit =
    exchange.getResponseHeaders.add("Content-Type", contentType)
    exchange.sendResponseHeaders(status, body.length)
    exchange.getResponseBody.write(body)

  private def token(): String =
    val bytes = new Array[Byte](32)
    random.nextBytes(bytes)
    Base64.getUrlEncoder.withoutPadding().encodeToString(bytes)

  private def env(name: String): String =
    sys.env.getOrElse(name, throw IllegalStateException(s"$name must be set, source ./auth0-env.sh"))
