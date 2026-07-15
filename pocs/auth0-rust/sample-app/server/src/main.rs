use auth0_rust::AuthApi;
use base64::Engine;
use base64::engine::general_purpose::URL_SAFE_NO_PAD;
use serde_json::{Map, Value, json};
use std::collections::HashMap;
use std::fs;
use std::io::{BufRead, BufReader, Read, Write};
use std::net::{TcpListener, TcpStream};
use std::path::{Path, PathBuf};
use std::sync::{Mutex, OnceLock};
use std::thread;

const SDK: &str = "rust";
const PORT: u16 = 3000;
const CALLBACK_PATH: &str = "/login/oauth2/code/okta";

fn sessions() -> &'static Mutex<HashMap<String, Value>> {
    static SESSIONS: OnceLock<Mutex<HashMap<String, Value>>> = OnceLock::new();
    SESSIONS.get_or_init(|| Mutex::new(HashMap::new()))
}

fn states() -> &'static Mutex<Vec<String>> {
    static STATES: OnceLock<Mutex<Vec<String>>> = OnceLock::new();
    STATES.get_or_init(|| Mutex::new(Vec::new()))
}

fn env(name: &str) -> String {
    std::env::var(name).unwrap_or_else(|_| panic!("{name} must be set, source ./auth0-env.sh"))
}

fn auth() -> &'static AuthApi {
    static AUTH: OnceLock<AuthApi> = OnceLock::new();
    AUTH.get_or_init(|| {
        AuthApi::builder(env("AUTH0_DOMAIN"), env("AUTH0_CLIENT_ID"))
            .client_secret(env("AUTH0_CLIENT_SECRET"))
            .build()
            .expect("Auth API client")
    })
}

fn redirect_uri() -> String {
    format!("http://localhost:{PORT}{CALLBACK_PATH}")
}

fn assets() -> PathBuf {
    PathBuf::from("webapp").join("dist")
}

fn token() -> String {
    let mut bytes = [0u8; 32];
    fs::File::open("/dev/urandom")
        .and_then(|mut file| file.read_exact(&mut bytes))
        .expect("random bytes");
    URL_SAFE_NO_PAD.encode(bytes)
}

fn main() {
    if !assets().is_dir() {
        eprintln!(
            "missing {}, run: cd webapp && npm install && npm run build",
            assets().display()
        );
        std::process::exit(1);
    }

    auth();
    let listener = TcpListener::bind(("127.0.0.1", PORT)).expect("bind port");
    println!("{SDK} webapp on http://localhost:{PORT}");

    for stream in listener.incoming() {
        match stream {
            Ok(stream) => {
                thread::spawn(|| handle(stream));
            }
            Err(error) => eprintln!("connection failed: {error}"),
        }
    }
}

fn handle(mut stream: TcpStream) {
    let peek = match stream.try_clone() {
        Ok(clone) => clone,
        Err(error) => return eprintln!("clone failed: {error}"),
    };
    let mut reader = BufReader::new(peek);

    let mut request_line = String::new();
    if reader.read_line(&mut request_line).is_err() {
        return;
    }

    let target = request_line.split_whitespace().nth(1).unwrap_or("/").to_string();

    let mut cookie = String::new();
    loop {
        let mut line = String::new();
        match reader.read_line(&mut line) {
            Ok(0) => break,
            Ok(_) => {
                if line.trim().is_empty() {
                    break;
                }
                if let Some((name, value)) = line.split_once(':') {
                    if name.trim().eq_ignore_ascii_case("cookie") {
                        cookie = value.trim().to_string();
                    }
                }
            }
            Err(_) => break,
        }
    }

    let (path, query) = match target.split_once('?') {
        Some((path, query)) => (path, query),
        None => (target.as_str(), ""),
    };

    match path {
        "/login" => start_login(&mut stream),
        CALLBACK_PATH => complete_login(&mut stream, query),
        "/api/me" => session(&mut stream, &cookie),
        "/logout" => logout(&mut stream, &cookie),
        _ => serve_static(&mut stream, path),
    }
}

fn start_login(stream: &mut TcpStream) {
    let state = token();
    states().lock().expect("states").push(state.clone());
    let url = auth()
        .authorize_url(redirect_uri())
        .scope("openid profile email")
        .state(state)
        .build();
    redirect(stream, &url);
}

fn complete_login(stream: &mut TcpStream, query: &str) {
    let params = parse_query(query);
    let code = params.get("code");
    let state = params.get("state");

    let known = match state {
        Some(state) => {
            let mut states = states().lock().expect("states");
            match states.iter().position(|value| value == state) {
                Some(index) => {
                    states.remove(index);
                    true
                }
                None => false,
            }
        }
        None => false,
    };

    let code = match (code, known) {
        (Some(code), true) => code,
        _ => return redirect(stream, "/"),
    };

    let response = match auth().exchange_code(code, redirect_uri()).execute() {
        Ok(response) => response,
        Err(error) => return send(stream, 500, "text/plain", format!("{error}").as_bytes()),
    };

    let id_token = response
        .body
        .as_ref()
        .and_then(|body| body.get("id_token"))
        .and_then(Value::as_str)
        .unwrap_or_default()
        .to_string();

    let claims = decode_claims(&id_token);
    let sid = token();
    sessions().lock().expect("sessions").insert(sid.clone(), claims);

    write_response(
        stream,
        302,
        &[
            ("Location", "/profile"),
            ("Set-Cookie", &format!("sid={sid}; Path=/; HttpOnly; SameSite=Lax")),
        ],
        &[],
    );
}

fn session(stream: &mut TcpStream, cookie: &str) {
    let claims = sid(cookie).and_then(|sid| sessions().lock().expect("sessions").get(&sid).cloned());
    let body = json!({
        "authenticated": claims.is_some(),
        "sdk": SDK,
        "claims": claims.unwrap_or_else(|| Value::Object(Map::new())),
    });
    send(stream, 200, "application/json", body.to_string().as_bytes());
}

fn logout(stream: &mut TcpStream, cookie: &str) {
    if let Some(sid) = sid(cookie) {
        sessions().lock().expect("sessions").remove(&sid);
    }
    let url = auth()
        .logout_url()
        .return_to(format!("http://localhost:{PORT}/"))
        .build();
    write_response(
        stream,
        302,
        &[("Location", &url), ("Set-Cookie", "sid=; Path=/; HttpOnly; Max-Age=0")],
        &[],
    );
}

fn serve_static(stream: &mut TcpStream, path: &str) {
    let requested = if path == "/" || path == "/profile" { "/index.html" } else { path };
    let root = assets();
    let file = root.join(requested.trim_start_matches('/'));

    let allowed = file
        .canonicalize()
        .ok()
        .zip(root.canonicalize().ok())
        .map(|(file, root)| file.starts_with(root))
        .unwrap_or(false);

    if !allowed || !file.is_file() {
        return send(stream, 404, "text/plain", b"not found");
    }

    match fs::read(&file) {
        Ok(body) => send(stream, 200, content_type(&file), &body),
        Err(_) => send(stream, 404, "text/plain", b"not found"),
    }
}

fn decode_claims(id_token: &str) -> Value {
    id_token
        .split('.')
        .nth(1)
        .and_then(|payload| URL_SAFE_NO_PAD.decode(payload).ok())
        .and_then(|bytes| serde_json::from_slice::<Value>(&bytes).ok())
        .unwrap_or_else(|| Value::Object(Map::new()))
}

fn content_type(file: &Path) -> &'static str {
    match file.extension().and_then(|value| value.to_str()) {
        Some("html") => "text/html",
        Some("js") => "text/javascript",
        Some("css") => "text/css",
        Some("svg") => "image/svg+xml",
        _ => "application/octet-stream",
    }
}

fn parse_query(query: &str) -> HashMap<String, String> {
    query
        .split('&')
        .filter_map(|pair| pair.split_once('='))
        .map(|(key, value)| (key.to_string(), decode_component(value)))
        .collect()
}

fn decode_component(value: &str) -> String {
    let bytes = value.as_bytes();
    let mut out = Vec::with_capacity(bytes.len());
    let mut index = 0;

    while index < bytes.len() {
        match bytes[index] {
            b'%' if index + 2 < bytes.len() => {
                match u8::from_str_radix(&value[index + 1..index + 3], 16) {
                    Ok(byte) => {
                        out.push(byte);
                        index += 3;
                    }
                    Err(_) => {
                        out.push(bytes[index]);
                        index += 1;
                    }
                }
            }
            b'+' => {
                out.push(b' ');
                index += 1;
            }
            byte => {
                out.push(byte);
                index += 1;
            }
        }
    }

    String::from_utf8_lossy(&out).into_owned()
}

fn sid(cookie: &str) -> Option<String> {
    cookie
        .split(';')
        .map(str::trim)
        .find_map(|pair| pair.strip_prefix("sid="))
        .map(str::to_string)
}

fn redirect(stream: &mut TcpStream, location: &str) {
    write_response(stream, 302, &[("Location", location)], &[]);
}

fn send(stream: &mut TcpStream, status: u16, content_type: &str, body: &[u8]) {
    write_response(stream, status, &[("Content-Type", content_type)], body);
}

fn write_response(stream: &mut TcpStream, status: u16, headers: &[(&str, &str)], body: &[u8]) {
    let reason = match status {
        200 => "OK",
        302 => "Found",
        404 => "Not Found",
        _ => "Internal Server Error",
    };

    let mut head = format!("HTTP/1.1 {status} {reason}\r\n");
    for (name, value) in headers {
        head.push_str(&format!("{name}: {value}\r\n"));
    }
    head.push_str(&format!("Content-Length: {}\r\n", body.len()));
    head.push_str("Connection: close\r\n\r\n");

    let _ = stream.write_all(head.as_bytes());
    let _ = stream.write_all(body);
    let _ = stream.flush();
}
