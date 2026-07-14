package auth0scala3

import java.net.URI
import java.net.InetSocketAddress
import java.net.ProxySelector
import java.net.URLEncoder
import java.net.http.HttpClient
import java.net.http.HttpRequest
import java.net.http.HttpResponse
import java.nio.charset.StandardCharsets
import java.security.Signature
import java.security.interfaces.RSAPrivateKey
import java.time.Duration
import java.time.Instant
import java.util.Base64
import java.util.UUID
import scala.collection.mutable
import scala.concurrent.ExecutionContext
import scala.concurrent.Future
import scala.concurrent.blocking
import scala.concurrent.duration.FiniteDuration
import scala.concurrent.duration.SECONDS

enum Method:
  case GET, POST, PUT, PATCH, DELETE

enum Json:
  case Null
  case Bool(value: Boolean)
  case Num(value: BigDecimal)
  case Str(value: String)
  case Arr(value: Vector[Json])
  case Obj(value: Map[String, Json])

object Json:
  def obj(values: (String, Json)*): Json = Obj(values.toMap)
  def arr(values: Json*): Json = Arr(values.toVector)
  def str(value: String): Json = Str(value)
  def num(value: BigDecimal): Json = Num(value)
  def bool(value: Boolean): Json = Bool(value)
  def from(value: Any): Json =
    value match
      case null => Null
      case v: Json => v
      case v: String => Str(v)
      case v: Boolean => Bool(v)
      case v: Int => Num(BigDecimal(v))
      case v: Long => Num(BigDecimal(v))
      case v: Double => Num(BigDecimal.decimal(v))
      case v: Float => Num(BigDecimal.decimal(v.toDouble))
      case v: BigDecimal => Num(v)
      case v: Map[?, ?] => Obj(v.toVector.map { case (k, x) => k.toString -> from(x) }.toMap)
      case v: Iterable[?] => Arr(v.toVector.map(from))
      case v: Product => Str(v.toString)
      case v => Str(v.toString)
  def stringify(json: Json): String =
    json match
      case Null => "null"
      case Bool(value) => value.toString
      case Num(value) => value.bigDecimal.stripTrailingZeros.toPlainString
      case Str(value) => "\"" + escape(value) + "\""
      case Arr(value) => value.map(stringify).mkString("[", ",", "]")
      case Obj(value) => value.toVector.map { case (k, v) => stringify(Str(k)) + ":" + stringify(v) }.mkString("{", ",", "}")
  def parse(input: String): Json =
    Parser(input).parse()
  private def escape(value: String): String =
    val out = new StringBuilder
    value.foreach {
      case '"' => out.append("\\\"")
      case '\\' => out.append("\\\\")
      case '\b' => out.append("\\b")
      case '\f' => out.append("\\f")
      case '\n' => out.append("\\n")
      case '\r' => out.append("\\r")
      case '\t' => out.append("\\t")
      case c if c < ' ' => out.append("\\u%04x".format(c.toInt))
      case c => out.append(c)
    }
    out.toString
  private final class Parser(input: String):
    private var i = 0
    def parse(): Json =
      val value = parseValue()
      spaces()
      if i != input.length then fail("trailing data")
      value
    private def parseValue(): Json =
      spaces()
      if i >= input.length then fail("unexpected end")
      input(i) match
        case 'n' => token("null", Null)
        case 't' => token("true", Bool(true))
        case 'f' => token("false", Bool(false))
        case '"' => Str(parseString())
        case '[' => parseArray()
        case '{' => parseObject()
        case c if c == '-' || c.isDigit => parseNumber()
        case _ => fail("unexpected character")
    private def parseObject(): Json =
      i += 1
      val fields = mutable.LinkedHashMap.empty[String, Json]
      spaces()
      if peek('}') then
        i += 1
        Obj(fields.toMap)
      else
        var done = false
        while !done do
          spaces()
          if !peek('"') then fail("object key expected")
          val key = parseString()
          spaces()
          if !peek(':') then fail("colon expected")
          i += 1
          fields += key -> parseValue()
          spaces()
          if peek(',') then i += 1
          else if peek('}') then
            i += 1
            done = true
          else fail("object separator expected")
        Obj(fields.toMap)
    private def parseArray(): Json =
      i += 1
      val values = Vector.newBuilder[Json]
      spaces()
      if peek(']') then
        i += 1
        Arr(values.result())
      else
        var done = false
        while !done do
          values += parseValue()
          spaces()
          if peek(',') then i += 1
          else if peek(']') then
            i += 1
            done = true
          else fail("array separator expected")
        Arr(values.result())
    private def parseString(): String =
      i += 1
      val out = new StringBuilder
      var done = false
      while !done do
        if i >= input.length then fail("unterminated string")
        val c = input(i)
        i += 1
        c match
          case '"' => done = true
          case '\\' =>
            if i >= input.length then fail("unterminated escape")
            val e = input(i)
            i += 1
            e match
              case '"' => out.append('"')
              case '\\' => out.append('\\')
              case '/' => out.append('/')
              case 'b' => out.append('\b')
              case 'f' => out.append('\f')
              case 'n' => out.append('\n')
              case 'r' => out.append('\r')
              case 't' => out.append('\t')
              case 'u' =>
                if i + 4 > input.length then fail("bad unicode escape")
                out.append(Integer.parseInt(input.substring(i, i + 4), 16).toChar)
                i += 4
              case _ => fail("bad escape")
          case x => out.append(x)
      out.toString
    private def parseNumber(): Json =
      val start = i
      if peek('-') then i += 1
      digits()
      if peek('.') then
        i += 1
        digits()
      if i < input.length && (input(i) == 'e' || input(i) == 'E') then
        i += 1
        if i < input.length && (input(i) == '+' || input(i) == '-') then i += 1
        digits()
      Num(BigDecimal(input.substring(start, i)))
    private def digits(): Unit =
      val start = i
      while i < input.length && input(i).isDigit do i += 1
      if start == i then fail("digit expected")
    private def token(text: String, value: Json): Json =
      if input.startsWith(text, i) then
        i += text.length
        value
      else fail(text + " expected")
    private def spaces(): Unit =
      while i < input.length && input(i).isWhitespace do i += 1
    private def peek(c: Char): Boolean = i < input.length && input(i) == c
    private def fail(message: String): Nothing = throw IllegalArgumentException(message + " at " + i)

sealed trait RequestBody:
  def contentType: String
  def text: String
  def bytes: Array[Byte] = text.getBytes(StandardCharsets.UTF_8)

object RequestBody:
  final case class Form(values: Map[String, String]) extends RequestBody:
    val contentType = "application/x-www-form-urlencoded"
    val text = Http.encodePairs(values.view.mapValues(Vector(_)).toMap)
  final case class JsonBody(value: Json) extends RequestBody:
    val contentType = "application/json"
    val text = Json.stringify(value)
  final case class Text(value: String, contentType: String = "text/plain") extends RequestBody:
    val text = value
  final case class Binary(value: Array[Byte], contentType: String) extends RequestBody:
    val text = String(value, StandardCharsets.UTF_8)
    override val bytes: Array[Byte] = value
  final case class Multipart(parts: Vector[MultipartPart], boundary: String = "auth0scala3-" + UUID.randomUUID().toString) extends RequestBody:
    val contentType = "multipart/form-data; boundary=" + boundary
    lazy val text = String(bytes, StandardCharsets.ISO_8859_1)
    override lazy val bytes: Array[Byte] =
      val out = java.io.ByteArrayOutputStream()
      parts.foreach { part =>
        out.write(("--" + boundary + "\r\n").getBytes(StandardCharsets.ISO_8859_1))
        out.write(("Content-Disposition: form-data; name=\"" + part.name + "\"" + part.fileName.map(v => "; filename=\"" + v + "\"").getOrElse("") + "\r\n").getBytes(StandardCharsets.ISO_8859_1))
        out.write(("Content-Type: " + part.contentType + "\r\n\r\n").getBytes(StandardCharsets.ISO_8859_1))
        out.write(part.bytes)
        out.write("\r\n".getBytes(StandardCharsets.ISO_8859_1))
      }
      out.write(("--" + boundary + "--\r\n").getBytes(StandardCharsets.ISO_8859_1))
      out.toByteArray

final case class MultipartPart(name: String, fileName: Option[String], contentType: String, bytes: Array[Byte])

final case class RequestOptions(
    headers: Map[String, String] = Map.empty,
    query: Map[String, Vector[String]] = Map.empty,
    timeout: Option[FiniteDuration] = None
)

final case class ApiRequest(
    method: Method,
    url: String,
    headers: Map[String, String] = Map.empty,
    query: Map[String, Vector[String]] = Map.empty,
    body: Option[RequestBody] = None,
    timeout: Option[FiniteDuration] = None
):
  def withHeader(name: String, value: String): ApiRequest = copy(headers = headers + (name -> value))
  def withQuery(name: String, value: String): ApiRequest = copy(query = query + (name -> (query.getOrElse(name, Vector.empty) :+ value)))
  def withHeaders(values: Map[String, String]): ApiRequest = copy(headers = headers ++ values)
  def withOptions(options: RequestOptions): ApiRequest =
    copy(headers = headers ++ options.headers, query = Http.mergeQuery(query, options.query), timeout = options.timeout.orElse(timeout))

final case class ApiResponse(status: Int, headers: Map[String, Vector[String]], body: String, bodyBytes: Array[Byte] = Array.emptyByteArray):
  def isSuccess: Boolean = status >= 200 && status < 300
  def json: Json = Json.parse(if body.trim.isEmpty then "null" else body)

final case class RawResponse[A](body: A, status: Int, headers: Map[String, Vector[String]], text: String, bytes: Array[Byte])

trait ResponseDecoder[A]:
  def decode(response: ApiResponse): A

object ResponseDecoder:
  given ResponseDecoder[ApiResponse] with
    def decode(response: ApiResponse): ApiResponse = response
  given ResponseDecoder[Json] with
    def decode(response: ApiResponse): Json = response.json
  given ResponseDecoder[String] with
    def decode(response: ApiResponse): String = response.body
  given ResponseDecoder[Array[Byte]] with
    def decode(response: ApiResponse): Array[Byte] = if response.bodyBytes.isEmpty then response.body.getBytes(StandardCharsets.UTF_8) else response.bodyBytes
  given ResponseDecoder[Unit] with
    def decode(response: ApiResponse): Unit = ()

open class ApiException(val status: Int, val responseBody: String, val request: ApiRequest) extends RuntimeException("Auth0 request failed with status " + status)
final class BadRequestError(body: String, request: ApiRequest) extends ApiException(400, body, request)
final class UnauthorizedError(body: String, request: ApiRequest) extends ApiException(401, body, request)
final class ForbiddenError(body: String, request: ApiRequest) extends ApiException(403, body, request)
final class NotFoundError(body: String, request: ApiRequest) extends ApiException(404, body, request)
final class ConflictError(body: String, request: ApiRequest) extends ApiException(409, body, request)
final class TooManyRequestsError(body: String, request: ApiRequest) extends ApiException(429, body, request)
final class ServiceUnavailableError(body: String, request: ApiRequest) extends ApiException(503, body, request)

object ApiException:
  def from(response: ApiResponse, request: ApiRequest): ApiException =
    response.status match
      case 400 => BadRequestError(response.body, request)
      case 401 => UnauthorizedError(response.body, request)
      case 403 => ForbiddenError(response.body, request)
      case 404 => NotFoundError(response.body, request)
      case 409 => ConflictError(response.body, request)
      case 429 => TooManyRequestsError(response.body, request)
      case 503 => ServiceUnavailableError(response.body, request)
      case code => new ApiException(code, response.body, request)

trait Transport:
  def send(request: ApiRequest): ApiResponse

trait Logger:
  def log(value: String): Unit

object ConsoleLogger extends Logger:
  def log(value: String): Unit = println(value)

final case class LoggingOptions(enabled: Boolean = false, logger: Logger = ConsoleLogger, redactHeaders: Set[String] = Set("Authorization"))
final case class ProxyOptions(host: String, port: Int)
final case class TransportOptions(proxy: Option[ProxyOptions] = None, retry: Int = 0, logging: LoggingOptions = LoggingOptions())

final class JavaNetTransport(client: HttpClient = HttpClient.newHttpClient(), options: TransportOptions = TransportOptions()) extends Transport:
  def send(request: ApiRequest): ApiResponse =
    sendWithRetry(request, 0)
  private def sendWithRetry(request: ApiRequest, attempt: Int): ApiResponse =
    val response =
      try sendOnce(request)
      catch
        case e: RuntimeException if attempt < options.retry =>
          return sendWithRetry(request, attempt + 1)
    if (response.status == 429 || response.status >= 500) && attempt < options.retry then sendWithRetry(request, attempt + 1) else response
  private def sendOnce(request: ApiRequest): ApiResponse =
    val body = request.body.map(_.bytes).getOrElse(Array.emptyByteArray)
    val builder = HttpRequest.newBuilder(URI.create(Http.urlWithQuery(request.url, request.query)))
    request.timeout.foreach(t => builder.timeout(Duration.ofNanos(t.toNanos)))
    request.body.foreach(b => builder.header("Content-Type", b.contentType))
    request.headers.foreach { case (name, value) => builder.header(name, value) }
    if options.logging.enabled then options.logging.logger.log(request.method.toString + " " + request.url + " " + redact(request.headers))
    request.method match
      case Method.GET => builder.GET()
      case Method.DELETE =>
        if body.isEmpty then builder.DELETE() else builder.method("DELETE", HttpRequest.BodyPublishers.ofByteArray(body))
      case other =>
        builder.method(other.toString, HttpRequest.BodyPublishers.ofByteArray(body))
    val response = client.send(builder.build(), HttpResponse.BodyHandlers.ofByteArray())
    val headers = response.headers().map().entrySet().toArray.toVector.map { entry =>
      val e = entry.asInstanceOf[java.util.Map.Entry[String, java.util.List[String]]]
      e.getKey -> e.getValue.toArray.toVector.map(_.toString)
    }.toMap
    val bytes = response.body()
    val text = String(bytes, StandardCharsets.UTF_8)
    val apiResponse = ApiResponse(response.statusCode(), headers, text, bytes)
    if options.logging.enabled then options.logging.logger.log(response.statusCode().toString + " " + request.url)
    apiResponse
  private def redact(headers: Map[String, String]): Map[String, String] =
    headers.map { case (k, v) => if options.logging.redactHeaders.exists(_.equalsIgnoreCase(k)) then k -> "[redacted]" else k -> v }

object JavaNetTransport:
  def from(options: TransportOptions): JavaNetTransport =
    val builder = HttpClient.newBuilder()
    options.proxy.foreach(p => builder.proxy(ProxySelector.of(InetSocketAddress(p.host, p.port))))
    JavaNetTransport(builder.build(), options)

object Http:
  def encode(value: String): String = URLEncoder.encode(value, StandardCharsets.UTF_8).replace("+", "%20")
  def encodePairs(values: Map[String, Vector[String]]): String =
    values.toVector.flatMap { case (k, vs) => vs.map(v => encode(k) + "=" + encode(v)) }.mkString("&")
  def urlWithQuery(url: String, query: Map[String, Vector[String]]): String =
    if query.isEmpty then url
    else
      val separator = if url.contains("?") then "&" else "?"
      url + separator + encodePairs(query)
  def mergeQuery(left: Map[String, Vector[String]], right: Map[String, Vector[String]]): Map[String, Vector[String]] =
    right.foldLeft(left) { case (all, (key, values)) => all + (key -> (all.getOrElse(key, Vector.empty) ++ values)) }

final case class ExecutableRequest(request: ApiRequest, transport: Transport):
  def withHeader(name: String, value: String): ExecutableRequest = copy(request = request.withHeader(name, value))
  def withQuery(name: String, value: String): ExecutableRequest = copy(request = request.withQuery(name, value))
  def withOptions(options: RequestOptions): ExecutableRequest = copy(request = request.withOptions(options))
  def withForm(name: String, value: String): ExecutableRequest =
    val next = request.body match
      case Some(RequestBody.Form(values)) => RequestBody.Form(values + (name -> value))
      case None => RequestBody.Form(Map(name -> value))
      case _ => throw IllegalStateException("request body is not form data")
    copy(request = request.copy(body = Some(next)))
  def send(): ApiResponse =
    val response = transport.send(request)
    if response.isSuccess then response else throw ApiException.from(response, request)
  def sendJson(): Json = send().json
  def raw[A](using decoder: ResponseDecoder[A]): RawRequest[A] = RawRequest(this, decoder)
  def sendAsync()(using ExecutionContext): Future[ApiResponse] = Future(blocking(send()))
  def sendJsonAsync()(using ExecutionContext): Future[Json] = Future(blocking(sendJson()))
  def download(): Array[Byte] = summon[ResponseDecoder[Array[Byte]]].decode(send())
  def stream(): Vector[SseEvent] = Sse.parse(send().body)
  def sendUnit(): Unit =
    send()
    ()

final case class RawRequest[A](value: ExecutableRequest, decoder: ResponseDecoder[A]):
  def send(): RawResponse[A] =
    val response = value.send()
    RawResponse(decoder.decode(response), response.status, response.headers, response.body, response.bodyBytes)
  def sendAsync()(using ExecutionContext): Future[RawResponse[A]] = Future(blocking(send()))

final case class SseEvent(event: Option[String], data: String, id: Option[String], retry: Option[Int])

object Sse:
  def parse(input: String): Vector[SseEvent] =
    input.split("\\r?\\n\\r?\\n").toVector.filter(_.trim.nonEmpty).map(parseEvent)
  private def parseEvent(input: String): SseEvent =
    val fields = input.split("\\r?\\n").toVector.flatMap { line =>
      val index = line.indexOf(":")
      if index < 0 then Some(line -> "") else Some(line.substring(0, index) -> line.substring(index + 1).stripPrefix(" "))
    }
    val data = fields.collect { case ("data", value) => value }.mkString("\n")
    SseEvent(fields.collectFirst { case ("event", value) => value }, data, fields.collectFirst { case ("id", value) => value }, fields.collectFirst { case ("retry", value) => value.toInt })

final case class TokenRequest(value: ExecutableRequest):
  def withRealm(realm: String): TokenRequest = copy(value.withForm("realm", realm))
  def withAudience(audience: String): TokenRequest = copy(value.withForm("audience", audience))
  def withScope(scope: String): TokenRequest = copy(value.withForm("scope", scope))
  def request: ApiRequest = value.request
  def send(): Json = value.sendJson()
  def sendAsync()(using ExecutionContext): Future[Json] = value.sendJsonAsync()

final case class SignUpRequest(value: ExecutableRequest):
  def withCustomFields(fields: Map[String, String]): SignUpRequest =
    val next = value.request.body match
      case Some(RequestBody.Form(values)) => RequestBody.JsonBody(Json.Obj(values.view.mapValues(Json.str).toMap + ("user_metadata" -> Json.from(fields))))
      case Some(RequestBody.JsonBody(Json.Obj(values))) => RequestBody.JsonBody(Json.Obj(values + ("user_metadata" -> Json.from(fields))))
      case _ => RequestBody.JsonBody(Json.obj("user_metadata" -> Json.from(fields)))
    copy(value.copy(request = value.request.copy(body = Some(next))))
  def request: ApiRequest = value.request
  def send(): Json = value.sendJson()
  def sendAsync()(using ExecutionContext): Future[Json] = value.sendJsonAsync()

trait ClientAssertionSigner:
  def createSignedClientAssertion(issuer: String, audience: String, subject: String): String

final class ClientAssertionSigningException(message: String, cause: Throwable | Null = null) extends RuntimeException(message, cause)

enum RSASigningAlgorithm(val jcaName: String, val jwtName: String):
  case RSA256 extends RSASigningAlgorithm("SHA256withRSA", "RS256")
  case RSA384 extends RSASigningAlgorithm("SHA384withRSA", "RS384")

final class RSAClientAssertionSigner(key: RSAPrivateKey, algorithm: RSASigningAlgorithm = RSASigningAlgorithm.RSA256) extends ClientAssertionSigner:
  def createSignedClientAssertion(issuer: String, audience: String, subject: String): String =
    try
      val now = Instant.now().getEpochSecond
      val header = Json.stringify(Json.obj("alg" -> Json.str(algorithm.jwtName), "typ" -> Json.str("JWT")))
      val payload = Json.stringify(Json.obj(
        "iss" -> Json.str(issuer),
        "aud" -> Json.str(audience),
        "sub" -> Json.str(subject),
        "iat" -> Json.num(now),
        "exp" -> Json.num(now + 180),
        "jti" -> Json.str(UUID.randomUUID().toString)
      ))
      val unsigned = Jwt.base64(header.getBytes(StandardCharsets.UTF_8)) + "." + Jwt.base64(payload.getBytes(StandardCharsets.UTF_8))
      unsigned + "." + Jwt.base64(sign(unsigned.getBytes(StandardCharsets.UTF_8)))
    catch
      case e: Exception => throw ClientAssertionSigningException("failed to create client assertion", e)
  private def sign(value: Array[Byte]): Array[Byte] =
    val signature = Signature.getInstance(algorithm.jcaName)
    signature.initSign(key)
    signature.update(value)
    signature.sign()

object Jwt:
  def base64(value: Array[Byte]): String = Base64.getUrlEncoder.withoutPadding().encodeToString(value)

enum PasswordlessEmailType(val value: String):
  case Code extends PasswordlessEmailType("code")
  case Link extends PasswordlessEmailType("link")

final case class AuthConfig(
    domain: String,
    clientId: String,
    clientSecret: Option[String] = None,
    clientAssertionSigner: Option[ClientAssertionSigner] = None,
    transport: Transport = JavaNetTransport()
)

final class AuthenticationApi(config: AuthConfig):
  private val baseUrl = Urls.baseDomain(config.domain)
  private val tokenUrl = baseUrl + "/oauth/token"
  def authorizeUrl(redirectUri: String): AuthorizeUrlBuilder =
    AuthorizeUrlBuilder(baseUrl, config.clientId, redirectUri)
  def authorizeUrlWithPAR(requestUri: String): String =
    Http.urlWithQuery(baseUrl + "/authorize", Map("client_id" -> Vector(config.clientId), "request_uri" -> Vector(requestUri)))
  def authorizeUrlWithJAR(request: String): String =
    Http.urlWithQuery(baseUrl + "/authorize", Map("client_id" -> Vector(config.clientId), "request" -> Vector(request)))
  def logoutUrl(returnToUrl: String, setClientId: Boolean): LogoutUrlBuilder =
    LogoutUrlBuilder(baseUrl, Option.when(setClientId)(config.clientId), returnToUrl)
  def raw(method: Method, path: String, body: Option[RequestBody] = None, headers: Map[String, String] = Map.empty): ExecutableRequest =
    ExecutableRequest(ApiRequest(method, Urls.join(baseUrl, path), headers, body = body), config.transport)
  def authorizeBackChannel(scope: String, bindingMessage: String, loginHint: Json, audience: Option[String] = None, requestExpiry: Option[Int] = None): ExecutableRequest =
    form(Method.POST, "/bc-authorize", clientAuth(required = false) ++ Map("client_id" -> config.clientId, "scope" -> scope, "binding_message" -> bindingMessage, "login_hint" -> Json.stringify(loginHint)) ++ audience.map("audience" -> _).toMap ++ requestExpiry.map(v => "requested_expiry" -> v.toString).toMap)
  def getBackChannelLoginStatus(authReqId: String, grantType: String): ExecutableRequest =
    form(Method.POST, "/oauth/token", clientAuth(required = false) ++ Map("client_id" -> config.clientId, "auth_req_id" -> authReqId, "grant_type" -> grantType))
  def pushedAuthorizationRequest(redirectUri: String, responseType: String, params: Map[String, String] = Map.empty, authorizationDetails: Option[Json] = None): ExecutableRequest =
    val values = clientAuth(required = false) ++ Map("client_id" -> config.clientId, "redirect_uri" -> redirectUri, "response_type" -> responseType) ++ params ++ authorizationDetails.map(v => "authorization_details" -> Json.stringify(v)).toMap
    form(Method.POST, "/oauth/par", values)
  def pushedAuthorizationRequestWithJAR(request: String, authorizationDetails: Option[Json] = None): ExecutableRequest =
    val values = clientAuth(required = false) ++ Map("client_id" -> config.clientId, "request" -> request) ++ authorizationDetails.map(v => "authorization_details" -> Json.stringify(v)).toMap
    form(Method.POST, "/oauth/par", values)
  def userInfo(accessToken: String): ExecutableRequest =
    raw(Method.GET, "/userinfo", headers = Map("Authorization" -> bearer(accessToken)))
  def resetPassword(email: String, connection: String, clientId: String = config.clientId, organization: Option[String] = None): ExecutableRequest =
    form(Method.POST, "/dbconnections/change_password", Map("client_id" -> clientId, "email" -> email, "connection" -> connection) ++ organization.map("organization" -> _).toMap)
  def signUp(email: String, password: String, connection: String, username: Option[String] = None, phoneNumber: Option[String] = None): SignUpRequest =
    val values = Map("client_id" -> config.clientId, "email" -> email, "password" -> password, "connection" -> connection) ++ username.map("username" -> _).toMap ++ phoneNumber.map("phone_number" -> _).toMap
    SignUpRequest(form(Method.POST, "/dbconnections/signup", values))
  def login(emailOrUsername: String, password: String): TokenRequest =
    token(Map("grant_type" -> "password", "username" -> emailOrUsername, "password" -> password), required = true)
  def login(emailOrUsername: String, password: String, realm: String): TokenRequest =
    token(Map("grant_type" -> "http://auth0.com/oauth/grant-type/password-realm", "username" -> emailOrUsername, "password" -> password, "realm" -> realm), required = true)
  def exchangePasswordlessOtp(emailOrPhone: String, realm: String, otp: String): TokenRequest =
    token(Map("grant_type" -> "http://auth0.com/oauth/grant-type/passwordless/otp", "username" -> emailOrPhone, "realm" -> realm, "otp" -> otp), required = false)
  def requestToken(audience: String, organization: Option[String] = None): TokenRequest =
    token(Map("grant_type" -> "client_credentials", "audience" -> audience) ++ organization.map("organization" -> _).toMap, required = true)
  def exchangeToken(subjectToken: String, subjectTokenType: String): TokenRequest =
    token(Map("grant_type" -> "urn:ietf:params:oauth:grant-type:token-exchange", "subject_token" -> subjectToken, "subject_token_type" -> subjectTokenType), required = true)
  def revokeToken(refreshToken: String): ExecutableRequest =
    form(Method.POST, "/oauth/revoke", clientAuth(required = false) ++ Map("client_id" -> config.clientId, "token" -> refreshToken))
  def renewAuth(refreshToken: String): TokenRequest =
    token(Map("grant_type" -> "refresh_token", "refresh_token" -> refreshToken), required = false)
  def exchangeCode(code: String, redirectUri: String): TokenRequest =
    token(Map("grant_type" -> "authorization_code", "code" -> code, "redirect_uri" -> redirectUri), required = true)
  def exchangeCodeWithVerifier(code: String, verifier: String, redirectUri: String): TokenRequest =
    token(Map("grant_type" -> "authorization_code", "code" -> code, "redirect_uri" -> redirectUri, "code_verifier" -> verifier), required = false)
  def startPasswordlessEmailFlow(email: String, kind: PasswordlessEmailType): ExecutableRequest =
    form(Method.POST, "/passwordless/start", clientAuth(required = false) ++ Map("client_id" -> config.clientId, "connection" -> "email", "email" -> email, "send" -> kind.value))
  def startPasswordlessSmsFlow(phoneNumber: String): ExecutableRequest =
    form(Method.POST, "/passwordless/start", clientAuth(required = false) ++ Map("client_id" -> config.clientId, "connection" -> "sms", "phone_number" -> phoneNumber))
  def exchangeMfaOtp(mfaToken: String, otp: String): TokenRequest =
    token(Map("grant_type" -> "http://auth0.com/oauth/grant-type/mfa-otp", "mfa_token" -> mfaToken, "otp" -> otp), required = false)
  def exchangeMfaOob(mfaToken: String, oobCode: String, bindingCode: Option[String] = None): TokenRequest =
    token(Map("grant_type" -> "http://auth0.com/oauth/grant-type/mfa-oob", "mfa_token" -> mfaToken, "oob_code" -> oobCode) ++ bindingCode.filter(_.nonEmpty).map("binding_code" -> _).toMap, required = false)
  def exchangeMfaRecoveryCode(mfaToken: String, recoveryCode: String): TokenRequest =
    token(Map("grant_type" -> "http://auth0.com/oauth/grant-type/mfa-recovery-code", "mfa_token" -> mfaToken, "recovery_code" -> recoveryCode), required = false)
  def mfaChallengeRequest(mfaToken: String, challengeType: Option[String] = None, authenticatorId: Option[String] = None): ExecutableRequest =
    form(Method.POST, "/mfa/challenge", clientAuth(required = false) ++ Map("client_id" -> config.clientId, "mfa_token" -> mfaToken) ++ challengeType.map("challenge_type" -> _).toMap ++ authenticatorId.map("authenticator_id" -> _).toMap)
  def addOtpAuthenticator(mfaToken: String): ExecutableRequest =
    raw(Method.POST, "/mfa/associate", Some(RequestBody.JsonBody(Json.obj("client_id" -> Json.str(config.clientId), "authenticator_types" -> Json.arr(Json.str("otp"))))), Map("Authorization" -> bearer(mfaToken)))  
  def addOobAuthenticator(mfaToken: String, oobChannels: Vector[String], phoneNumber: Option[String] = None, emailAddress: Option[String] = None): ExecutableRequest =
    if (oobChannels.contains("sms") || oobChannels.contains("voice")) && phoneNumber.isEmpty then throw IllegalArgumentException("phoneNumber is required")
    if oobChannels.contains("email") && emailAddress.isEmpty then throw IllegalArgumentException("emailAddress is required")
    val body = Json.obj(
      "client_id" -> Json.str(config.clientId),
      "authenticator_types" -> Json.arr(Json.str("oob")),
      "oob_channels" -> Json.Arr(oobChannels.map(Json.str)),
      "phone_number" -> phoneNumber.map(Json.str).getOrElse(Json.Null),
      "email" -> emailAddress.map(Json.str).getOrElse(Json.Null)
    )
    raw(Method.POST, "/mfa/associate", Some(RequestBody.JsonBody(body)), Map("Authorization" -> bearer(mfaToken)))
  def listAuthenticators(accessToken: String): ExecutableRequest =
    raw(Method.GET, "/mfa/authenticators", headers = Map("Authorization" -> bearer(accessToken)))
  def deleteAuthenticator(accessToken: String, authenticatorId: String): ExecutableRequest =
    raw(Method.DELETE, "/mfa/authenticators/" + Http.encode(authenticatorId), headers = Map("Authorization" -> bearer(accessToken)))
  private def token(values: Map[String, String], required: Boolean): TokenRequest =
    TokenRequest(form(Method.POST, "/oauth/token", clientAuth(required) ++ Map("client_id" -> config.clientId) ++ values))
  private def form(method: Method, path: String, values: Map[String, String]): ExecutableRequest =
    raw(method, path, Some(RequestBody.Form(values)))
  private def clientAuth(required: Boolean): Map[String, String] =
    if required && config.clientSecret.isEmpty && config.clientAssertionSigner.isEmpty then throw IllegalStateException("client authentication is required")
    config.clientAssertionSigner.map(s => Map("client_assertion" -> s.createSignedClientAssertion(config.clientId, baseUrl, config.clientId), "client_assertion_type" -> "urn:ietf:params:oauth:client-assertion-type:jwt-bearer")).orElse(config.clientSecret.map(v => Map("client_secret" -> v))).getOrElse(Map.empty)
  private def bearer(token: String): String = "Bearer " + token

final case class AuthorizeUrlBuilder(baseUrl: String, clientId: String, redirectUri: String, parameters: Map[String, String] = Map("response_type" -> "code")):
  def withConnection(value: String): AuthorizeUrlBuilder = withParameter("connection", value)
  def withAudience(value: String): AuthorizeUrlBuilder = withParameter("audience", value)
  def withState(value: String): AuthorizeUrlBuilder = withParameter("state", value)
  def withScope(value: String): AuthorizeUrlBuilder = withParameter("scope", value)
  def withResponseType(value: String): AuthorizeUrlBuilder = withParameter("response_type", value)
  def withOrganization(value: String): AuthorizeUrlBuilder = withParameter("organization", value)
  def withInvitation(value: String): AuthorizeUrlBuilder = withParameter("invitation", value)
  def withCodeChallenge(value: String): AuthorizeUrlBuilder = copy(parameters = parameters + ("code_challenge" -> value) + ("code_challenge_method" -> "S256"))
  def withParameter(name: String, value: String): AuthorizeUrlBuilder = copy(parameters = parameters + (name -> value))
  def build: String = Http.urlWithQuery(baseUrl + "/authorize", Map("redirect_uri" -> Vector(redirectUri), "client_id" -> Vector(clientId)) ++ parameters.view.mapValues(Vector(_)).toMap)

final case class LogoutUrlBuilder(baseUrl: String, clientId: Option[String], returnToUrl: String, federated: Boolean = false):
  def useFederated(value: Boolean): LogoutUrlBuilder = copy(federated = value)
  def build: String =
    val base = Map("returnTo" -> Vector(returnToUrl)) ++ clientId.map(v => "client_id" -> Vector(v)).toMap
    val all = if federated then base + ("federated" -> Vector("")) else base
    Http.urlWithQuery(baseUrl + "/v2/logout", all)

final case class ClientCredentials(clientId: String, clientSecret: String)

final case class ManagementConfig(
    domain: Option[String] = None,
    url: Option[String] = None,
    token: Option[String] = None,
    clientCredentials: Option[ClientCredentials] = None,
    audience: Option[String] = None,
    headers: Map[String, String] = Map.empty,
    customDomain: Option[String] = None,
    timeout: FiniteDuration = FiniteDuration(60, SECONDS),
    transport: Transport = JavaNetTransport()
)

final class ManagementApi(config: ManagementConfig):
  private val baseUrl = config.url.getOrElse(config.domain.map(d => Urls.baseDomain(d) + "/api/v2").getOrElse("https://auth0.com/api/v2")).stripSuffix("/")
  private val tokenSupplier = config.clientCredentials.map(c => CachedToken(config.domain.map(Urls.baseDomain).getOrElse(baseUrl.stripSuffix("/api/v2")), c, config.audience.getOrElse(config.domain.map(d => Urls.baseDomain(d) + "/api/v2/").getOrElse(baseUrl + "/")), config.transport))
  def raw(method: Method, path: String, query: Map[String, Vector[String]] = Map.empty, body: Option[RequestBody] = None, headers: Map[String, String] = Map.empty, options: RequestOptions = RequestOptions()): ExecutableRequest =
    val authHeader = config.token.map(bearer).orElse(tokenSupplier.map(t => bearer(t.get))).map("Authorization" -> _).toMap
    val custom = config.customDomain.filter(_ => CustomDomainWhitelist.allows(path)).map("Auth0-Custom-Domain" -> _).toMap
    val allQuery = Http.mergeQuery(query, options.query)
    val allHeaders = config.headers ++ custom ++ authHeader ++ headers ++ options.headers
    ExecutableRequest(ApiRequest(method, Urls.join(baseUrl, path), allHeaders, allQuery, body, options.timeout.orElse(Some(config.timeout))), config.transport)
  def async: AsyncManagementApi = AsyncManagementApi(this)
  def resource(path: String): ManagementResource = ManagementResource(this, Vector(path))
  def actions: ActionsResource = ActionsResource(resource("actions"))
  def anomaly: ManagementResource = resource("anomaly")
  def attackProtection: AttackProtectionResource = AttackProtectionResource(resource("attack-protection"))
  def branding: BrandingResource = BrandingResource(resource("branding"))
  def clientGrants: ClientGrantsResource = ClientGrantsResource(resource("client-grants"))
  def clients: ClientsResource = ClientsResource(resource("clients"))
  def connectionProfiles: ManagementResource = resource("connection-profiles")
  def connections: ConnectionsResource = ConnectionsResource(resource("connections"))
  def customDomains: ManagementResource = resource("custom-domains")
  def deviceCredentials: ManagementResource = resource("device-credentials")
  def emails: EmailsResource = EmailsResource(resource("emails"))
  def emailTemplates: ManagementResource = resource("email-templates")
  def eventStreams: EventStreamsResource = EventStreamsResource(resource("event-streams"))
  def events: ManagementResource = resource("events")
  def flows: FlowsResource = FlowsResource(resource("flows"))
  def forms: ManagementResource = resource("forms")
  def groups: GroupsResource = GroupsResource(resource("groups"))
  def guardian: GuardianResource = GuardianResource(resource("guardian"))
  def hooks: HooksResource = HooksResource(resource("hooks"))
  def jobs: JobsResource = JobsResource(resource("jobs"))
  def keys: KeysResource = KeysResource(resource("keys"))
  def logStreams: ManagementResource = resource("log-streams")
  def logs: ManagementResource = resource("logs")
  def networkAcls: ManagementResource = resource("network-acls")
  def organizations: OrganizationsResource = OrganizationsResource(resource("organizations"))
  def prompts: PromptsResource = PromptsResource(resource("prompts"))
  def rateLimitPolicies: ManagementResource = resource("rate-limit-policies")
  def refreshTokens: ManagementResource = resource("refresh-tokens")
  def resourceServers: ManagementResource = resource("resource-servers")
  def riskAssessments: RiskAssessmentsResource = RiskAssessmentsResource(resource("risk-assessments"))
  def roles: RolesResource = RolesResource(resource("roles"))
  def rules: ManagementResource = resource("rules")
  def rulesConfigs: ManagementResource = resource("rules-configs")
  def selfServiceProfiles: SelfServiceProfilesResource = SelfServiceProfilesResource(resource("self-service-profiles"))
  def sessions: ManagementResource = resource("sessions")
  def stats: ManagementResource = resource("stats")
  def supplementalSignals: ManagementResource = resource("supplemental-signals")
  def tenants: TenantsResource = TenantsResource(resource("tenants"))
  def tickets: ManagementResource = resource("tickets")
  def tokenExchangeProfiles: ManagementResource = resource("token-exchange-profiles")
  def userAttributeProfiles: ManagementResource = resource("user-attribute-profiles")
  def userBlocks: ManagementResource = resource("user-blocks")
  def userGrants: ManagementResource = resource("user-grants")
  def users: UsersResource = UsersResource(resource("users"))
  def verifiableCredentials: VerifiableCredentialsResource = VerifiableCredentialsResource(resource("verifiable-credentials"))
  private def bearer(token: String): String = "Bearer " + token

object ManagementApi:
  def withToken(domain: String, token: String, transport: Transport = JavaNetTransport()): ManagementApi =
    ManagementApi(ManagementConfig(domain = Some(domain), token = Some(token), transport = transport))
  def withClientCredentials(domain: String, clientId: String, clientSecret: String, transport: Transport = JavaNetTransport()): ManagementApi =
    ManagementApi(ManagementConfig(domain = Some(domain), clientCredentials = Some(ClientCredentials(clientId, clientSecret)), transport = transport))

object CustomDomainWhitelist:
  def allows(path: String): Boolean =
    val value = path.stripPrefix("/")
    value == "jobs/verification-email" ||
      value == "tickets/email-verification" ||
      value == "tickets/password-change" ||
      value == "users" ||
      value.startsWith("users/") ||
      value.matches("organizations/[^/]+/invitations") ||
      value == "guardian/enrollments/ticket" ||
      value.matches("self-service-profiles/[^/]+/sso-ticket")

trait ResourceOps:
  def resource: ManagementResource
  def at(values: String*): ManagementResource = resource.at(values*)
  def get(query: Map[String, Vector[String]] = Map.empty, headers: Map[String, String] = Map.empty, options: RequestOptions = RequestOptions()): ExecutableRequest = resource.get(query, headers, options)
  def post(json: Json, query: Map[String, Vector[String]] = Map.empty, headers: Map[String, String] = Map.empty, options: RequestOptions = RequestOptions()): ExecutableRequest = resource.post(json, query, headers, options)
  def put(json: Json, query: Map[String, Vector[String]] = Map.empty, headers: Map[String, String] = Map.empty, options: RequestOptions = RequestOptions()): ExecutableRequest = resource.put(json, query, headers, options)
  def patch(json: Json, query: Map[String, Vector[String]] = Map.empty, headers: Map[String, String] = Map.empty, options: RequestOptions = RequestOptions()): ExecutableRequest = resource.patch(json, query, headers, options)
  def delete(query: Map[String, Vector[String]] = Map.empty, headers: Map[String, String] = Map.empty, options: RequestOptions = RequestOptions()): ExecutableRequest = resource.delete(query, headers, options)
  def upload(name: String, fileName: String, contentType: String, bytes: Array[Byte], fields: Map[String, String] = Map.empty, options: RequestOptions = RequestOptions()): ExecutableRequest = resource.upload(name, fileName, contentType, bytes, fields, options)

final case class ManagementResource(api: ManagementApi, segments: Vector[String]) extends ResourceOps:
  def resource: ManagementResource = this
  override def at(values: String*): ManagementResource = copy(segments = segments ++ values)
  override def get(query: Map[String, Vector[String]] = Map.empty, headers: Map[String, String] = Map.empty, options: RequestOptions = RequestOptions()): ExecutableRequest =
    api.raw(Method.GET, path, query, headers = headers, options = options)
  override def post(json: Json, query: Map[String, Vector[String]] = Map.empty, headers: Map[String, String] = Map.empty, options: RequestOptions = RequestOptions()): ExecutableRequest =
    api.raw(Method.POST, path, query, Some(RequestBody.JsonBody(json)), headers, options)
  override def put(json: Json, query: Map[String, Vector[String]] = Map.empty, headers: Map[String, String] = Map.empty, options: RequestOptions = RequestOptions()): ExecutableRequest =
    api.raw(Method.PUT, path, query, Some(RequestBody.JsonBody(json)), headers, options)
  override def patch(json: Json, query: Map[String, Vector[String]] = Map.empty, headers: Map[String, String] = Map.empty, options: RequestOptions = RequestOptions()): ExecutableRequest =
    api.raw(Method.PATCH, path, query, Some(RequestBody.JsonBody(json)), headers, options)
  override def delete(query: Map[String, Vector[String]] = Map.empty, headers: Map[String, String] = Map.empty, options: RequestOptions = RequestOptions()): ExecutableRequest =
    api.raw(Method.DELETE, path, query, headers = headers, options = options)
  override def upload(name: String, fileName: String, contentType: String, bytes: Array[Byte], fields: Map[String, String] = Map.empty, options: RequestOptions = RequestOptions()): ExecutableRequest =
    val parts = fields.toVector.map { case (k, v) => MultipartPart(k, None, "text/plain", v.getBytes(StandardCharsets.UTF_8)) } :+ MultipartPart(name, Some(fileName), contentType, bytes)
    api.raw(Method.POST, path, body = Some(RequestBody.Multipart(parts)), options = options)
  def path: String = segments.map(Http.encode).mkString("/")

final case class AsyncManagementApi(api: ManagementApi):
  def resource(path: String): AsyncManagementResource = AsyncManagementResource(api.resource(path))
  def users: AsyncManagementResource = AsyncManagementResource(api.users.resource)
  def organizations: AsyncManagementResource = AsyncManagementResource(api.organizations.resource)
  def guardian: AsyncManagementResource = AsyncManagementResource(api.guardian.resource)
  def jobs: AsyncManagementResource = AsyncManagementResource(api.jobs.resource)

final case class AsyncManagementResource(resource: ManagementResource):
  def at(values: String*): AsyncManagementResource = copy(resource.at(values*))
  def get(query: Map[String, Vector[String]] = Map.empty, headers: Map[String, String] = Map.empty, options: RequestOptions = RequestOptions())(using ExecutionContext): Future[ApiResponse] =
    resource.get(query, headers, options).sendAsync()
  def post(json: Json, query: Map[String, Vector[String]] = Map.empty, headers: Map[String, String] = Map.empty, options: RequestOptions = RequestOptions())(using ExecutionContext): Future[ApiResponse] =
    resource.post(json, query, headers, options).sendAsync()
  def put(json: Json, query: Map[String, Vector[String]] = Map.empty, headers: Map[String, String] = Map.empty, options: RequestOptions = RequestOptions())(using ExecutionContext): Future[ApiResponse] =
    resource.put(json, query, headers, options).sendAsync()
  def patch(json: Json, query: Map[String, Vector[String]] = Map.empty, headers: Map[String, String] = Map.empty, options: RequestOptions = RequestOptions())(using ExecutionContext): Future[ApiResponse] =
    resource.patch(json, query, headers, options).sendAsync()
  def delete(query: Map[String, Vector[String]] = Map.empty, headers: Map[String, String] = Map.empty, options: RequestOptions = RequestOptions())(using ExecutionContext): Future[ApiResponse] =
    resource.delete(query, headers, options).sendAsync()

final case class ActionsResource(resource: ManagementResource) extends ResourceOps:
  def executions(actionId: String): ManagementResource = resource.at(actionId, "executions")
  def modules(actionId: String): ManagementResource = resource.at(actionId, "modules")
  def triggers: ManagementResource = resource.at("triggers")
  def triggerBindings(triggerId: String): ManagementResource = resource.at("triggers", triggerId, "bindings")
  def versions(actionId: String): ManagementResource = resource.at(actionId, "versions")

final case class AttackProtectionResource(resource: ManagementResource) extends ResourceOps:
  def botDetection: ManagementResource = resource.at("bot-detection")
  def breachedPasswordDetection: ManagementResource = resource.at("breached-password-detection")
  def bruteForceProtection: ManagementResource = resource.at("brute-force-protection")
  def captcha: ManagementResource = resource.at("captcha")
  def phoneProviderProtection: ManagementResource = resource.at("phone-provider-protection")
  def suspiciousIpThrottling: ManagementResource = resource.at("suspicious-ip-throttling")

final case class BrandingResource(resource: ManagementResource) extends ResourceOps:
  def templates: ManagementResource = resource.at("templates")
  def themes: ManagementResource = resource.at("themes")
  def phone: BrandingPhoneResource = BrandingPhoneResource(resource.at("phone"))

final case class BrandingPhoneResource(resource: ManagementResource) extends ResourceOps:
  def providers: ManagementResource = resource.at("providers")
  def templates: ManagementResource = resource.at("templates")

final case class ClientGrantsResource(resource: ManagementResource) extends ResourceOps:
  def organizations(clientGrantId: String): ManagementResource = resource.at(clientGrantId, "organizations")

final case class ClientsResource(resource: ManagementResource) extends ResourceOps:
  def connections(clientId: String): ManagementResource = resource.at(clientId, "connections")
  def credentials(clientId: String): ManagementResource = resource.at(clientId, "credentials")

final case class ConnectionsResource(resource: ManagementResource) extends ResourceOps:
  def clients(connectionId: String): ManagementResource = resource.at(connectionId, "clients")
  def directoryProvisioning(connectionId: String): ConnectionDirectoryProvisioningResource = ConnectionDirectoryProvisioningResource(resource.at(connectionId, "directory-provisioning"))
  def keys(connectionId: String): ManagementResource = resource.at(connectionId, "keys")
  def scimConfiguration(connectionId: String): ConnectionScimConfigurationResource = ConnectionScimConfigurationResource(resource.at(connectionId, "scim-configuration"))
  def users(connectionId: String): ManagementResource = resource.at(connectionId, "users")

final case class ConnectionDirectoryProvisioningResource(resource: ManagementResource) extends ResourceOps:
  def synchronizations: ManagementResource = resource.at("synchronizations")

final case class ConnectionScimConfigurationResource(resource: ManagementResource) extends ResourceOps:
  def tokens: ManagementResource = resource.at("tokens")

final case class EmailsResource(resource: ManagementResource) extends ResourceOps:
  def provider: ManagementResource = resource.at("provider")

final case class EventStreamsResource(resource: ManagementResource) extends ResourceOps:
  def deliveries(eventStreamId: String): ManagementResource = resource.at(eventStreamId, "deliveries")
  def redeliveries(eventStreamId: String): ManagementResource = resource.at(eventStreamId, "redeliveries")

final case class FlowsResource(resource: ManagementResource) extends ResourceOps:
  def executions(flowId: String): ManagementResource = resource.at(flowId, "executions")
  def vault: FlowsVaultResource = FlowsVaultResource(resource.at("vault"))

final case class FlowsVaultResource(resource: ManagementResource) extends ResourceOps:
  def connections: ManagementResource = resource.at("connections")

final case class GroupsResource(resource: ManagementResource) extends ResourceOps:
  def members(groupId: String): ManagementResource = resource.at(groupId, "members")
  def roles(groupId: String): ManagementResource = resource.at(groupId, "roles")

final case class GuardianResource(resource: ManagementResource) extends ResourceOps:
  def enrollments: ManagementResource = resource.at("enrollments")
  def enrollmentTicket: ManagementResource = resource.at("enrollments", "ticket")
  def factors: GuardianFactorsResource = GuardianFactorsResource(resource.at("factors"))
  def policies: ManagementResource = resource.at("policies")

final case class GuardianFactorsResource(resource: ManagementResource) extends ResourceOps:
  def duo: GuardianDuoResource = GuardianDuoResource(resource.at("duo"))
  def phone: ManagementResource = resource.at("phone")
  def pushNotification: ManagementResource = resource.at("push-notification")
  def sms: ManagementResource = resource.at("sms")

final case class GuardianDuoResource(resource: ManagementResource) extends ResourceOps:
  def settings: ManagementResource = resource.at("settings")

final case class HooksResource(resource: ManagementResource) extends ResourceOps:
  def secrets(hookId: String): ManagementResource = resource.at(hookId, "secrets")

final case class JobsResource(resource: ManagementResource) extends ResourceOps:
  def errors(jobId: String): ManagementResource = resource.at(jobId, "errors")
  def usersExports: ManagementResource = resource.at("users-exports")
  def usersImports: ManagementResource = resource.at("users-imports")
  def verificationEmail: ManagementResource = resource.at("verification-email")

final case class KeysResource(resource: ManagementResource) extends ResourceOps:
  def customSigning: ManagementResource = resource.at("custom-signing")
  def encryption: ManagementResource = resource.at("encryption")
  def signing: ManagementResource = resource.at("signing")

final case class OrganizationsResource(resource: ManagementResource) extends ResourceOps:
  def clientGrants(organizationId: String): ManagementResource = resource.at(organizationId, "client-grants")
  def connections(organizationId: String): ManagementResource = resource.at(organizationId, "connections")
  def discoveryDomains(organizationId: String): ManagementResource = resource.at(organizationId, "discovery-domains")
  def enabledConnections(organizationId: String): ManagementResource = resource.at(organizationId, "enabled-connections")
  def groups(organizationId: String): OrganizationGroupsResource = OrganizationGroupsResource(resource.at(organizationId, "groups"))
  def invitations(organizationId: String): ManagementResource = resource.at(organizationId, "invitations")
  def members(organizationId: String): OrganizationMembersResource = OrganizationMembersResource(resource.at(organizationId, "members"))

final case class OrganizationGroupsResource(resource: ManagementResource) extends ResourceOps:
  def roles(groupId: String): ManagementResource = resource.at(groupId, "roles")

final case class OrganizationMembersResource(resource: ManagementResource) extends ResourceOps:
  def effectiveRoles(memberId: String): OrganizationMemberEffectiveRolesResource = OrganizationMemberEffectiveRolesResource(resource.at(memberId, "effective-roles"))
  def roles(memberId: String): ManagementResource = resource.at(memberId, "roles")

final case class OrganizationMemberEffectiveRolesResource(resource: ManagementResource) extends ResourceOps:
  def sources: EffectiveRoleSourcesResource = EffectiveRoleSourcesResource(resource.at("sources"))

final case class PromptsResource(resource: ManagementResource) extends ResourceOps:
  def customText(prompt: String): ManagementResource = resource.at(prompt, "custom-text")
  def partials(prompt: String): ManagementResource = resource.at(prompt, "partials")
  def rendering(prompt: String): ManagementResource = resource.at(prompt, "rendering")

final case class RiskAssessmentsResource(resource: ManagementResource) extends ResourceOps:
  def settings: RiskAssessmentSettingsResource = RiskAssessmentSettingsResource(resource.at("settings"))

final case class RiskAssessmentSettingsResource(resource: ManagementResource) extends ResourceOps:
  def newDevice: ManagementResource = resource.at("new-device")

final case class RolesResource(resource: ManagementResource) extends ResourceOps:
  def groups(roleId: String): ManagementResource = resource.at(roleId, "groups")
  def permissions(roleId: String): ManagementResource = resource.at(roleId, "permissions")
  def users(roleId: String): ManagementResource = resource.at(roleId, "users")

final case class SelfServiceProfilesResource(resource: ManagementResource) extends ResourceOps:
  def customText(profileId: String): ManagementResource = resource.at(profileId, "custom-text")
  def ssoTicket(profileId: String): ManagementResource = resource.at(profileId, "sso-ticket")

final case class TenantsResource(resource: ManagementResource) extends ResourceOps:
  def settings: ManagementResource = resource.at("settings")

final case class UsersResource(resource: ManagementResource) extends ResourceOps:
  def authenticationMethods(userId: String): ManagementResource = resource.at(userId, "authentication-methods")
  def authenticators(userId: String): ManagementResource = resource.at(userId, "authenticators")
  def connectedAccounts(userId: String): ManagementResource = resource.at(userId, "connected-accounts")
  def effectivePermissions(userId: String): UserEffectivePermissionsResource = UserEffectivePermissionsResource(resource.at(userId, "effective-permissions"))
  def effectiveRoles(userId: String): UserEffectiveRolesResource = UserEffectiveRolesResource(resource.at(userId, "effective-roles"))
  def enrollments(userId: String): ManagementResource = resource.at(userId, "enrollments")
  def federatedConnectionsTokensets(userId: String): ManagementResource = resource.at(userId, "federated-connections-tokensets")
  def groups(userId: String): ManagementResource = resource.at(userId, "groups")
  def identities(userId: String): ManagementResource = resource.at(userId, "identities")
  def logs(userId: String): ManagementResource = resource.at(userId, "logs")
  def multifactor(userId: String): ManagementResource = resource.at(userId, "multifactor")
  def organizations(userId: String): ManagementResource = resource.at(userId, "organizations")
  def permissions(userId: String): ManagementResource = resource.at(userId, "permissions")
  def refreshToken(userId: String): ManagementResource = resource.at(userId, "refresh-token")
  def riskAssessments(userId: String): ManagementResource = resource.at(userId, "risk-assessments")
  def roles(userId: String): ManagementResource = resource.at(userId, "roles")
  def sessions(userId: String): ManagementResource = resource.at(userId, "sessions")

final case class UserEffectivePermissionsResource(resource: ManagementResource) extends ResourceOps:
  def sources: EffectivePermissionSourcesResource = EffectivePermissionSourcesResource(resource.at("sources"))

final case class UserEffectiveRolesResource(resource: ManagementResource) extends ResourceOps:
  def sources: EffectiveRoleSourcesResource = EffectiveRoleSourcesResource(resource.at("sources"))

final case class EffectivePermissionSourcesResource(resource: ManagementResource) extends ResourceOps:
  def roles: ManagementResource = resource.at("roles")

final case class EffectiveRoleSourcesResource(resource: ManagementResource) extends ResourceOps:
  def groups: ManagementResource = resource.at("groups")

final case class VerifiableCredentialsResource(resource: ManagementResource) extends ResourceOps:
  def verification: VerifiableCredentialVerificationResource = VerifiableCredentialVerificationResource(resource.at("verification"))

final case class VerifiableCredentialVerificationResource(resource: ManagementResource) extends ResourceOps:
  def templates: ManagementResource = resource.at("templates")

final class CachedToken(baseUrl: String, credentials: ClientCredentials, audience: String, transport: Transport):
  private var token = Option.empty[String]
  def get: String =
    token.getOrElse {
      val response = AuthenticationApi(AuthConfig(baseUrl, credentials.clientId, Some(credentials.clientSecret), transport = transport)).requestToken(audience).send()
      val accessToken = response match
        case Json.Obj(values) => values.get("access_token").collect { case Json.Str(v) => v }.getOrElse(throw IllegalStateException("access_token missing"))
        case _ => throw IllegalStateException("access_token missing")
      token = Some(accessToken)
      accessToken
    }

final case class Auth0(auth: AuthenticationApi, management: Option[ManagementApi])

object Auth0:
  def auth(domain: String, clientId: String, clientSecret: Option[String] = None, transport: Transport = JavaNetTransport()): AuthenticationApi =
    AuthenticationApi(AuthConfig(domain, clientId, clientSecret, transport = transport))
  def withManagementToken(domain: String, clientId: String, clientSecret: Option[String], token: String, transport: Transport = JavaNetTransport()): Auth0 =
    Auth0(auth(domain, clientId, clientSecret, transport), Some(ManagementApi.withToken(domain, token, transport)))

object Urls:
  def baseDomain(domain: String): String =
    val withScheme = if domain.startsWith("http://") || domain.startsWith("https://") then domain else "https://" + domain
    withScheme.stripSuffix("/")
  def join(base: String, path: String): String =
    base.stripSuffix("/") + "/" + path.stripPrefix("/")
