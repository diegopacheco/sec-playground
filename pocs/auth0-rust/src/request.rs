use crate::assertion::RsaClientAssertionSigner;
use crate::client::{ClientOptions, LogLevel, RetryOptions};
use crate::error::{ApiError, Auth0Error};
use serde::de::DeserializeOwned;
use serde_json::Value;
use std::collections::HashMap;
use std::thread;
use std::time::Duration;
use url::Url;

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum Method {
    Get,
    Post,
    Put,
    Patch,
    Delete,
}

impl Method {
    pub fn as_str(self) -> &'static str {
        match self {
            Method::Get => "GET",
            Method::Post => "POST",
            Method::Put => "PUT",
            Method::Patch => "PATCH",
            Method::Delete => "DELETE",
        }
    }

    fn as_reqwest(self) -> reqwest::Method {
        match self {
            Method::Get => reqwest::Method::GET,
            Method::Post => reqwest::Method::POST,
            Method::Put => reqwest::Method::PUT,
            Method::Patch => reqwest::Method::PATCH,
            Method::Delete => reqwest::Method::DELETE,
        }
    }

    fn is_idempotent(self) -> bool {
        matches!(self, Method::Get | Method::Put | Method::Delete)
    }
}

#[derive(Debug, Clone)]
struct DynamicClientAssertion {
    signer: RsaClientAssertionSigner,
    issuer: String,
    subject: String,
}

#[derive(Debug, Clone, PartialEq)]
pub enum Body {
    Empty,
    Json(Value),
    Form(Vec<(String, String)>),
    Multipart(Vec<MultipartPart>),
    Text(String),
}

#[derive(Debug, Clone, PartialEq, Eq)]
pub struct MultipartPart {
    pub name: String,
    pub file_name: Option<String>,
    pub mime: Option<String>,
    pub bytes: Vec<u8>,
}

impl MultipartPart {
    pub fn text(name: impl Into<String>, value: impl Into<String>) -> Self {
        Self {
            name: name.into(),
            file_name: None,
            mime: None,
            bytes: value.into().into_bytes(),
        }
    }

    pub fn file(name: impl Into<String>, file_name: impl Into<String>, bytes: Vec<u8>) -> Self {
        Self {
            name: name.into(),
            file_name: Some(file_name.into()),
            mime: None,
            bytes,
        }
    }

    pub fn mime(mut self, value: impl Into<String>) -> Self {
        self.mime = Some(value.into());
        self
    }
}

#[derive(Debug, Clone, PartialEq)]
pub struct PreparedRequest {
    pub method: Method,
    pub url: String,
    pub headers: Vec<(String, String)>,
    pub body: Body,
    pub timeout: Option<Duration>,
}

#[derive(Debug, Clone, PartialEq)]
pub struct ApiResponse {
    pub status: u16,
    pub headers: HashMap<String, String>,
    pub body: Option<Value>,
    pub text: String,
}

#[derive(Debug, Clone, Default, PartialEq, Eq)]
pub struct RequestOptions {
    pub headers: Vec<(String, String)>,
    pub query: Vec<(String, String)>,
    pub timeout: Option<Duration>,
}

impl RequestOptions {
    pub fn new() -> Self {
        Self::default()
    }

    pub fn header(mut self, key: impl Into<String>, value: impl Into<String>) -> Self {
        self.headers.push((key.into(), value.into()));
        self
    }

    pub fn query(mut self, key: impl Into<String>, value: impl Into<String>) -> Self {
        self.query.push((key.into(), value.into()));
        self
    }

    pub fn timeout(mut self, value: Duration) -> Self {
        self.timeout = Some(value);
        self
    }
}

#[derive(Debug, Clone)]
pub struct Request {
    method: Method,
    url: Url,
    headers: Vec<(String, String)>,
    query: Vec<(String, String)>,
    body: Body,
    timeout: Option<Duration>,
    client_assertion: Option<DynamicClientAssertion>,
}

impl Request {
    pub fn new(method: Method, url: Url) -> Self {
        Self {
            method,
            url,
            headers: Vec::new(),
            query: Vec::new(),
            body: Body::Empty,
            timeout: None,
            client_assertion: None,
        }
    }

    pub fn header(mut self, key: impl Into<String>, value: impl Into<String>) -> Self {
        self.headers.push((key.into(), value.into()));
        self
    }

    pub fn bearer(self, token: impl AsRef<str>) -> Self {
        self.header("Authorization", format!("Bearer {}", token.as_ref()))
    }

    pub fn query(mut self, key: impl Into<String>, value: impl Into<String>) -> Self {
        self.query.push((key.into(), value.into()));
        self
    }

    pub fn options(mut self, options: RequestOptions) -> Self {
        for (key, value) in options.headers {
            self.headers.push((key, value));
        }
        for (key, value) in options.query {
            self.query.push((key, value));
        }
        if options.timeout.is_some() {
            self.timeout = options.timeout;
        }
        self
    }

    pub fn timeout(mut self, value: Duration) -> Self {
        self.timeout = Some(value);
        self
    }

    pub fn json(mut self, value: Value) -> Self {
        self.body = Body::Json(value);
        self
    }

    pub fn form(mut self, values: Vec<(String, String)>) -> Self {
        self.body = Body::Form(values);
        self
    }

    pub fn form_append(mut self, key: impl Into<String>, value: impl Into<String>) -> Self {
        match &mut self.body {
            Body::Form(values) => values.push((key.into(), value.into())),
            _ => self.body = Body::Form(vec![(key.into(), value.into())]),
        }
        self
    }

    pub fn json_insert(mut self, key: impl Into<String>, value: Value) -> Self {
        match &mut self.body {
            Body::Json(Value::Object(map)) => {
                map.insert(key.into(), value);
            }
            _ => {
                let mut map = serde_json::Map::new();
                map.insert(key.into(), value);
                self.body = Body::Json(Value::Object(map));
            }
        }
        self
    }

    pub fn multipart(mut self, values: Vec<MultipartPart>) -> Self {
        self.body = Body::Multipart(values);
        self
    }

    pub fn text(mut self, value: impl Into<String>) -> Self {
        self.body = Body::Text(value.into());
        self
    }

    pub(crate) fn client_assertion(
        mut self,
        signer: RsaClientAssertionSigner,
        issuer: impl Into<String>,
        subject: impl Into<String>,
    ) -> Self {
        self.client_assertion = Some(DynamicClientAssertion {
            signer,
            issuer: issuer.into(),
            subject: subject.into(),
        });
        self
    }

    pub fn prepared(&self) -> PreparedRequest {
        let mut url = self.url.clone();
        if !self.query.is_empty() {
            let mut pairs = url.query_pairs_mut();
            for (key, value) in &self.query {
                pairs.append_pair(key, value);
            }
        }
        PreparedRequest {
            method: self.method,
            url: url.to_string(),
            headers: self.headers.clone(),
            body: self.body.clone(),
            timeout: self.timeout,
        }
    }

    pub fn execute(&self) -> Result<ApiResponse, Auth0Error> {
        self.execute_with_options(&ClientOptions::default())
    }

    pub fn execute_json<T: DeserializeOwned>(&self) -> Result<T, Auth0Error> {
        let response = self.execute()?;
        decode_json(response)
    }

    pub fn execute_with_options(&self, options: &ClientOptions) -> Result<ApiResponse, Auth0Error> {
        let mut attempt = 0;
        loop {
            let prepared = self.prepare_for_execution()?;
            log_request(options, &prepared);
            let response = match send_blocking(options.blocking_client()?, &prepared) {
                Ok(value) => value,
                Err(_) if should_retry_transport(&options.retry, attempt, prepared.method) => {
                    let delay = retry_delay(&options.retry, &HashMap::new(), attempt);
                    attempt += 1;
                    thread::sleep(delay);
                    continue;
                }
                Err(value) => return Err(value),
            };
            log_response(options, &response);
            if should_retry_status(&options.retry, attempt, prepared.method, response.status) {
                let delay = retry_delay(&options.retry, &response.headers, attempt);
                attempt += 1;
                thread::sleep(delay);
                continue;
            }
            return response.into_result();
        }
    }

    pub fn execute_json_with_options<T: DeserializeOwned>(
        &self,
        options: &ClientOptions,
    ) -> Result<T, Auth0Error> {
        let response = self.execute_with_options(options)?;
        decode_json(response)
    }

    pub async fn execute_async(&self) -> Result<ApiResponse, Auth0Error> {
        self.execute_async_with_options(&ClientOptions::default())
            .await
    }

    pub async fn execute_json_async<T: DeserializeOwned>(&self) -> Result<T, Auth0Error> {
        let response = self.execute_async().await?;
        decode_json(response)
    }

    pub async fn execute_async_with_options(
        &self,
        options: &ClientOptions,
    ) -> Result<ApiResponse, Auth0Error> {
        let mut attempt = 0;
        loop {
            let prepared = self.prepare_for_execution()?;
            log_request(options, &prepared);
            let response = match send_async(&options.async_client, &prepared).await {
                Ok(value) => value,
                Err(_) if should_retry_transport(&options.retry, attempt, prepared.method) => {
                    let delay = retry_delay(&options.retry, &HashMap::new(), attempt);
                    attempt += 1;
                    tokio::time::sleep(delay).await;
                    continue;
                }
                Err(value) => return Err(value),
            };
            log_response(options, &response);
            if should_retry_status(&options.retry, attempt, prepared.method, response.status) {
                let delay = retry_delay(&options.retry, &response.headers, attempt);
                attempt += 1;
                tokio::time::sleep(delay).await;
                continue;
            }
            return response.into_result();
        }
    }

    pub async fn execute_json_async_with_options<T: DeserializeOwned>(
        &self,
        options: &ClientOptions,
    ) -> Result<T, Auth0Error> {
        let response = self.execute_async_with_options(options).await?;
        decode_json(response)
    }

    fn prepare_for_execution(&self) -> Result<PreparedRequest, Auth0Error> {
        let mut prepared = self.prepared();
        if let Some(assertion) = &self.client_assertion {
            let value = assertion.signer.sign(
                assertion.issuer.clone(),
                prepared.url.clone(),
                assertion.subject.clone(),
            )?;
            match &mut prepared.body {
                Body::Form(values) => {
                    values.retain(|(key, _)| {
                        key != "client_assertion" && key != "client_assertion_type"
                    });
                    values.push(("client_assertion".into(), value));
                    values.push((
                        "client_assertion_type".into(),
                        "urn:ietf:params:oauth:client-assertion-type:jwt-bearer".into(),
                    ));
                }
                _ => {
                    return Err(Auth0Error::InvalidInput(
                        "client assertion requires a form request".into(),
                    ));
                }
            }
        }
        Ok(prepared)
    }
}

fn send_blocking(
    client: &reqwest::blocking::Client,
    prepared: &PreparedRequest,
) -> Result<ApiResponse, Auth0Error> {
    let mut builder = client.request(prepared.method.as_reqwest(), &prepared.url);
    for (key, value) in &prepared.headers {
        builder = builder.header(key, value);
    }
    if let Some(timeout) = prepared.timeout {
        builder = builder.timeout(timeout);
    }
    builder = match &prepared.body {
        Body::Empty => builder,
        Body::Json(value) => builder.json(value),
        Body::Form(values) => builder.form(values),
        Body::Multipart(values) => builder.multipart(blocking_multipart(values)?),
        Body::Text(value) => builder.body(value.clone()),
    };
    let response = builder.send()?;
    read_blocking_response(response)
}

async fn send_async(
    client: &reqwest::Client,
    prepared: &PreparedRequest,
) -> Result<ApiResponse, Auth0Error> {
    let mut builder = client.request(prepared.method.as_reqwest(), &prepared.url);
    for (key, value) in &prepared.headers {
        builder = builder.header(key, value);
    }
    if let Some(timeout) = prepared.timeout {
        builder = builder.timeout(timeout);
    }
    builder = match &prepared.body {
        Body::Empty => builder,
        Body::Json(value) => builder.json(value),
        Body::Form(values) => builder.form(values),
        Body::Multipart(values) => builder.multipart(async_multipart(values)?),
        Body::Text(value) => builder.body(value.clone()),
    };
    let response = builder.send().await?;
    read_async_response(response).await
}

fn blocking_multipart(
    values: &[MultipartPart],
) -> Result<reqwest::blocking::multipart::Form, Auth0Error> {
    let mut form = reqwest::blocking::multipart::Form::new();
    for value in values {
        let mut part = reqwest::blocking::multipart::Part::bytes(value.bytes.clone());
        if let Some(file_name) = &value.file_name {
            part = part.file_name(file_name.clone());
        }
        if let Some(mime) = &value.mime {
            part = part
                .mime_str(mime)
                .map_err(|value| Auth0Error::InvalidInput(value.to_string()))?;
        }
        form = form.part(value.name.clone(), part);
    }
    Ok(form)
}

fn async_multipart(values: &[MultipartPart]) -> Result<reqwest::multipart::Form, Auth0Error> {
    let mut form = reqwest::multipart::Form::new();
    for value in values {
        let mut part = reqwest::multipart::Part::bytes(value.bytes.clone());
        if let Some(file_name) = &value.file_name {
            part = part.file_name(file_name.clone());
        }
        if let Some(mime) = &value.mime {
            part = part
                .mime_str(mime)
                .map_err(|value| Auth0Error::InvalidInput(value.to_string()))?;
        }
        form = form.part(value.name.clone(), part);
    }
    Ok(form)
}

fn read_blocking_response(
    response: reqwest::blocking::Response,
) -> Result<ApiResponse, Auth0Error> {
    let status = response.status().as_u16();
    let headers = header_map(response.headers());
    let text = response.text()?;
    api_response(status, headers, text)
}

async fn read_async_response(response: reqwest::Response) -> Result<ApiResponse, Auth0Error> {
    let status = response.status().as_u16();
    let headers = header_map(response.headers());
    let text = response.text().await?;
    api_response(status, headers, text)
}

fn header_map(headers: &reqwest::header::HeaderMap) -> HashMap<String, String> {
    let mut map = HashMap::new();
    for (key, value) in headers {
        map.insert(
            key.as_str().to_string(),
            value.to_str().unwrap_or("").to_string(),
        );
    }
    map
}

fn api_response(
    status: u16,
    headers: HashMap<String, String>,
    text: String,
) -> Result<ApiResponse, Auth0Error> {
    let body = if text.trim().is_empty() {
        None
    } else {
        serde_json::from_str(&text).ok()
    };
    Ok(ApiResponse {
        status,
        headers,
        body,
        text,
    })
}

fn decode_json<T: DeserializeOwned>(response: ApiResponse) -> Result<T, Auth0Error> {
    match response.body {
        Some(value) => Ok(serde_json::from_value(value)?),
        None => Ok(serde_json::from_str("null")?),
    }
}

fn should_retry_status(retry: &RetryOptions, attempt: usize, method: Method, status: u16) -> bool {
    attempt < retry.max_retries
        && retry.should_retry(status)
        && (method.is_idempotent() || retry.retry_non_idempotent)
}

fn should_retry_transport(retry: &RetryOptions, attempt: usize, method: Method) -> bool {
    attempt < retry.max_retries && (method.is_idempotent() || retry.retry_non_idempotent)
}

fn retry_delay(
    retry: &RetryOptions,
    headers: &HashMap<String, String>,
    attempt: usize,
) -> Duration {
    if retry.respect_retry_after
        && let Some(value) = headers
            .get("retry-after")
            .or_else(|| headers.get("Retry-After"))
            .and_then(|value| value.parse::<u64>().ok())
    {
        return Duration::from_secs(value).min(retry.max_delay);
    }
    let multiplier = 1u32.checked_shl(attempt.min(31) as u32).unwrap_or(u32::MAX);
    retry
        .base_delay
        .saturating_mul(multiplier)
        .min(retry.max_delay)
}

fn log_request(options: &ClientOptions, prepared: &PreparedRequest) {
    match options.logging.level {
        LogLevel::None => {}
        LogLevel::Basic => eprintln!("{} {}", prepared.method.as_str(), prepared.url),
        LogLevel::Headers | LogLevel::Body => {
            eprintln!("{} {}", prepared.method.as_str(), prepared.url);
            for (key, value) in &prepared.headers {
                eprintln!("{}: {}", key, options.logging.redact_value(key, value));
            }
            if matches!(options.logging.level, LogLevel::Body) {
                eprintln!(
                    "{}",
                    options.logging.body_value(format!("{:?}", prepared.body))
                );
            }
        }
    }
}

fn log_response(options: &ClientOptions, response: &ApiResponse) {
    match options.logging.level {
        LogLevel::None => {}
        LogLevel::Basic => eprintln!("status {}", response.status),
        LogLevel::Headers | LogLevel::Body => {
            eprintln!("status {}", response.status);
            for (key, value) in &response.headers {
                eprintln!("{}: {}", key, options.logging.redact_value(key, value));
            }
            if matches!(options.logging.level, LogLevel::Body) {
                eprintln!("{}", options.logging.body_value(&response.text));
            }
        }
    }
}

impl ApiResponse {
    fn into_result(self) -> Result<Self, Auth0Error> {
        if self.status >= 400 {
            return Err(Auth0Error::Http(parse_api_error(self.status, &self.text)));
        }
        Ok(self)
    }

    pub fn json<T: DeserializeOwned>(self) -> Result<T, Auth0Error> {
        decode_json(self)
    }
}

pub fn parse_api_error(status: u16, body: &str) -> ApiError {
    let parsed = serde_json::from_str::<Value>(body).ok();
    let code = parsed
        .as_ref()
        .and_then(|value| value.get("error").or_else(|| value.get("code")))
        .and_then(Value::as_str)
        .map(str::to_string);
    let message = parsed
        .as_ref()
        .and_then(|value| {
            value
                .get("error_description")
                .or_else(|| value.get("message"))
                .or_else(|| value.get("description"))
        })
        .and_then(Value::as_str)
        .map(str::to_string);
    ApiError {
        status,
        code,
        message,
        body: body.to_string(),
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn retries_idempotent_methods_only_by_default() {
        let retry = RetryOptions::new(1);

        assert!(should_retry_status(&retry, 0, Method::Get, 503));
        assert!(should_retry_transport(&retry, 0, Method::Delete));
        assert!(!should_retry_status(&retry, 0, Method::Post, 503));
        assert!(!should_retry_transport(&retry, 0, Method::Patch));
    }

    #[test]
    fn retry_delay_is_bounded() {
        let retry = RetryOptions::new(1).max_delay(Duration::from_secs(2));
        let headers = HashMap::from([("retry-after".into(), "60".into())]);

        assert_eq!(retry_delay(&retry, &headers, 0), Duration::from_secs(2));
    }
}
