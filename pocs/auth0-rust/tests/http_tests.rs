use auth0_rust::{AuthApi, ManagementApi};
use std::io::{Read, Write};
use std::net::{TcpListener, TcpStream};
use std::sync::{Arc, Mutex};
use std::thread;

struct TestServer {
    domain: String,
    requests: Arc<Mutex<Vec<String>>>,
    handle: thread::JoinHandle<()>,
}

impl TestServer {
    fn start(responses: Vec<(u16, &'static str)>) -> Self {
        let listener = TcpListener::bind("127.0.0.1:0").expect("listener");
        let address = listener.local_addr().expect("address");
        let requests = Arc::new(Mutex::new(Vec::new()));
        let captured = requests.clone();
        let handle = thread::spawn(move || {
            for (status, body) in responses {
                let (mut stream, _) = listener.accept().expect("connection");
                let request = read_request(&mut stream);
                captured.lock().expect("requests").push(request);
                write_response(&mut stream, status, body);
            }
        });
        Self {
            domain: format!("http://{}", address),
            requests,
            handle,
        }
    }

    fn finish(self) -> Vec<String> {
        self.handle.join().expect("server thread");
        Arc::try_unwrap(self.requests)
            .expect("request ownership")
            .into_inner()
            .expect("requests")
    }
}

fn read_request(stream: &mut TcpStream) -> String {
    let mut bytes = Vec::new();
    let mut buffer = [0; 2048];
    loop {
        let count = stream.read(&mut buffer).expect("request bytes");
        if count == 0 {
            break;
        }
        bytes.extend_from_slice(&buffer[..count]);
        if let Some(header_end) = find_header_end(&bytes) {
            let headers = String::from_utf8_lossy(&bytes[..header_end]);
            let content_length = headers
                .lines()
                .find_map(|line| {
                    line.split_once(':').and_then(|(name, value)| {
                        name.eq_ignore_ascii_case("content-length")
                            .then(|| value.trim().parse::<usize>().ok())
                            .flatten()
                    })
                })
                .unwrap_or(0);
            if bytes.len() >= header_end + 4 + content_length {
                break;
            }
        }
    }
    String::from_utf8(bytes).expect("request text")
}

fn find_header_end(bytes: &[u8]) -> Option<usize> {
    bytes.windows(4).position(|value| value == b"\r\n\r\n")
}

fn write_response(stream: &mut TcpStream, status: u16, body: &str) {
    let reason = if status == 200 { "OK" } else { "Unauthorized" };
    let response = format!(
        "HTTP/1.1 {} {}\r\nContent-Type: application/json\r\nContent-Length: {}\r\nConnection: close\r\n\r\n{}",
        status,
        reason,
        body.len(),
        body
    );
    stream.write_all(response.as_bytes()).expect("response");
}

#[test]
fn management_refreshes_a_rejected_automatic_token() {
    let server = TestServer::start(vec![
        (200, r#"{"access_token":"token-1","expires_in":3600}"#),
        (401, r#"{"error":"invalid_token"}"#),
        (200, r#"{"access_token":"token-2","expires_in":3600}"#),
        (200, "[]"),
    ]);
    let api = ManagementApi::builder()
        .domain(&server.domain)
        .client_credentials("client", "secret")
        .build()
        .expect("management client");

    let response = api
        .clients_list()
        .expect("request")
        .execute()
        .expect("response");

    assert_eq!(response.status, 200);
    let requests = server.finish();
    assert!(requests[1].contains("authorization: Bearer token-1"));
    assert!(requests[3].contains("authorization: Bearer token-2"));
}

#[test]
fn async_management_refreshes_expired_tokens_without_blocking_calls() {
    let server = TestServer::start(vec![
        (200, r#"{"access_token":"token-1","expires_in":0}"#),
        (200, "[]"),
        (200, r#"{"access_token":"token-2","expires_in":0}"#),
        (200, "[]"),
    ]);
    let api = ManagementApi::builder()
        .domain(&server.domain)
        .client_credentials("client", "secret")
        .build()
        .expect("management client");
    let runtime = tokio::runtime::Runtime::new().expect("runtime");

    runtime.block_on(async {
        api.clients_list()
            .expect("first request")
            .execute_async()
            .await
            .expect("first response");
        api.clients_list()
            .expect("second request")
            .execute_async()
            .await
            .expect("second response");
    });

    let requests = server.finish();
    assert!(requests[1].contains("authorization: Bearer token-1"));
    assert!(requests[3].contains("authorization: Bearer token-2"));
}

#[test]
fn async_context_can_own_and_drop_clients() {
    let runtime = tokio::runtime::Runtime::new().expect("runtime");

    runtime.block_on(async {
        let auth = AuthApi::builder("tenant.auth0.com", "client")
            .build()
            .expect("authentication client");
        let management = ManagementApi::builder()
            .domain("tenant.auth0.com")
            .token("token")
            .build()
            .expect("management client");
        drop(auth);
        drop(management);
    });
}
