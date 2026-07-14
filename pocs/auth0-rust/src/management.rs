use crate::auth::AuthApi;
use crate::client::ClientOptions;
use crate::error::Auth0Error;
use crate::generated::{Endpoint, MANAGEMENT_ENDPOINTS};
use crate::request::{
    ApiResponse, Body, Method, MultipartPart, PreparedRequest, Request, RequestOptions,
};
use serde::de::DeserializeOwned;
use serde_json::Value;
use std::collections::HashMap;
use std::marker::PhantomData;
use std::sync::Mutex;
use url::Url;

pub struct ManagementApiBuilder {
    domain: Option<String>,
    token: Option<String>,
    client_id: Option<String>,
    client_secret: Option<String>,
    client_options: ClientOptions,
}

impl ManagementApiBuilder {
    pub fn new() -> Self {
        Self {
            domain: None,
            token: None,
            client_id: None,
            client_secret: None,
            client_options: ClientOptions::default(),
        }
    }

    pub fn domain(mut self, value: impl Into<String>) -> Self {
        self.domain = Some(value.into());
        self
    }

    pub fn token(mut self, value: impl Into<String>) -> Self {
        self.token = Some(value.into());
        self
    }

    pub fn client_credentials(
        mut self,
        client_id: impl Into<String>,
        client_secret: impl Into<String>,
    ) -> Self {
        self.client_id = Some(client_id.into());
        self.client_secret = Some(client_secret.into());
        self
    }

    pub fn client_options(mut self, value: ClientOptions) -> Self {
        self.client_options = value;
        self
    }

    pub fn build(self) -> Result<ManagementApi, Auth0Error> {
        let domain = self
            .domain
            .ok_or_else(|| Auth0Error::InvalidInput("domain".into()))?;
        ManagementApi::new(
            domain,
            self.token,
            self.client_id,
            self.client_secret,
            self.client_options,
        )
    }
}

impl Default for ManagementApiBuilder {
    fn default() -> Self {
        Self::new()
    }
}

pub struct ManagementApi {
    domain: String,
    base_url: Url,
    token: Option<String>,
    client_id: Option<String>,
    client_secret: Option<String>,
    client_options: ClientOptions,
    token_cache: Mutex<Option<String>>,
}

impl ManagementApi {
    pub fn builder() -> ManagementApiBuilder {
        ManagementApiBuilder::new()
    }

    pub fn new(
        domain: impl Into<String>,
        token: Option<String>,
        client_id: Option<String>,
        client_secret: Option<String>,
        client_options: ClientOptions,
    ) -> Result<Self, Auth0Error> {
        let domain = domain.into();
        let auth_base = crate::auth::normalize_domain(&domain)?;
        let base_url = auth_base.join("api/v2/")?;
        Ok(Self {
            domain,
            base_url,
            token,
            client_id,
            client_secret,
            client_options,
            token_cache: Mutex::new(None),
        })
    }

    pub fn endpoint(&self, name: &str) -> Option<&'static Endpoint> {
        MANAGEMENT_ENDPOINTS
            .iter()
            .find(|endpoint| endpoint.name == name)
    }

    pub fn endpoints(&self) -> &'static [Endpoint] {
        MANAGEMENT_ENDPOINTS
    }

    pub fn request(&self, name: &str) -> Result<ManagementRequest<'_>, Auth0Error> {
        let endpoint = self
            .endpoint(name)
            .ok_or_else(|| Auth0Error::MissingEndpoint(name.to_string()))?;
        Ok(ManagementRequest::new(self, endpoint))
    }

    pub fn resource(&self, group: impl Into<String>) -> ResourceClient<'_> {
        ResourceClient {
            api: self,
            group: group.into(),
        }
    }

    pub fn raw(&self) -> RawManagementApi<'_> {
        RawManagementApi { api: self }
    }

    pub fn actions(&self) -> ResourceClient<'_> {
        self.resource("actions")
    }

    pub fn clients(&self) -> ResourceClient<'_> {
        self.resource("clients")
    }

    pub fn connections(&self) -> ResourceClient<'_> {
        self.resource("connections")
    }

    pub fn organizations(&self) -> ResourceClient<'_> {
        self.resource("organizations")
    }

    pub fn roles(&self) -> ResourceClient<'_> {
        self.resource("roles")
    }

    pub fn users(&self) -> ResourceClient<'_> {
        self.resource("users")
    }

    fn access_token(&self) -> Result<String, Auth0Error> {
        if let Some(value) = &self.token {
            return Ok(value.clone());
        }
        {
            let guard = self.token_cache.lock().unwrap();
            if let Some(value) = guard.as_ref() {
                return Ok(value.clone());
            }
        }
        let client_id = self
            .client_id
            .clone()
            .ok_or(Auth0Error::MissingAccessToken)?;
        let client_secret = self
            .client_secret
            .clone()
            .ok_or(Auth0Error::MissingAccessToken)?;
        let auth = AuthApi::builder(&self.domain, client_id)
            .client_secret(client_secret)
            .client_options(self.client_options.clone())
            .build()?;
        let audience = format!(
            "{}/api/v2/",
            crate::auth::normalize_domain(&self.domain)?
                .as_str()
                .trim_end_matches('/')
        );
        let response = auth
            .request_token(audience)
            .execute_with_options(&self.client_options)?;
        let token = response
            .body
            .and_then(|value| {
                value
                    .get("access_token")
                    .and_then(Value::as_str)
                    .map(str::to_string)
            })
            .ok_or(Auth0Error::MissingAccessToken)?;
        let mut guard = self.token_cache.lock().unwrap();
        *guard = Some(token.clone());
        Ok(token)
    }
}

pub struct RawManagementApi<'a> {
    api: &'a ManagementApi,
}

impl<'a> RawManagementApi<'a> {
    pub fn request(&self, name: &str) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.api.request(name)
    }

    pub fn endpoints(&self) -> &'static [Endpoint] {
        self.api.endpoints()
    }
}

#[derive(Clone)]
pub struct ResourceClient<'a> {
    api: &'a ManagementApi,
    group: String,
}

impl<'a> ResourceClient<'a> {
    pub fn request(&self, operation: &str) -> Result<ManagementRequest<'a>, Auth0Error> {
        self.api.request(&format!("{}.{}", self.group, operation))
    }

    pub fn endpoint(&self, operation: &str) -> Option<&'static Endpoint> {
        self.api.endpoint(&format!("{}.{}", self.group, operation))
    }

    pub fn endpoints(&self) -> Vec<&'static Endpoint> {
        MANAGEMENT_ENDPOINTS
            .iter()
            .filter(|endpoint| endpoint.group == self.group)
            .collect()
    }
}

#[derive(Clone)]
pub struct ManagementRequest<'a> {
    api: &'a ManagementApi,
    endpoint: &'static Endpoint,
    path_params: HashMap<String, String>,
    query: Vec<(String, String)>,
    body: Option<Body>,
    headers: Vec<(String, String)>,
    options: RequestOptions,
}

impl<'a> ManagementRequest<'a> {
    pub fn new(api: &'a ManagementApi, endpoint: &'static Endpoint) -> Self {
        Self {
            api,
            endpoint,
            path_params: HashMap::new(),
            query: Vec::new(),
            body: None,
            headers: Vec::new(),
            options: RequestOptions::new(),
        }
    }

    pub fn path_param(mut self, key: impl Into<String>, value: impl Into<String>) -> Self {
        self.path_params.insert(key.into(), value.into());
        self
    }

    pub fn query(mut self, key: impl Into<String>, value: impl Into<String>) -> Self {
        self.query.push((key.into(), value.into()));
        self
    }

    pub fn body(mut self, value: Value) -> Self {
        self.body = Some(Body::Json(value));
        self
    }

    pub fn multipart(mut self, value: Vec<MultipartPart>) -> Self {
        self.body = Some(Body::Multipart(value));
        self
    }

    pub fn header(mut self, key: impl Into<String>, value: impl Into<String>) -> Self {
        self.headers.push((key.into(), value.into()));
        self
    }

    pub fn options(mut self, value: RequestOptions) -> Self {
        self.options = value;
        self
    }

    pub fn endpoint(&self) -> &'static Endpoint {
        self.endpoint
    }

    pub fn prepared(&self) -> Result<PreparedRequest, Auth0Error> {
        self.to_request()?
            .bearer(self.api.access_token()?)
            .prepared()
            .into_ok()
    }

    pub fn execute(&self) -> Result<ApiResponse, Auth0Error> {
        self.to_request()?
            .bearer(self.api.access_token()?)
            .execute_with_options(&self.api.client_options)
    }

    pub fn execute_json<T: DeserializeOwned>(&self) -> Result<T, Auth0Error> {
        self.execute()?.json()
    }

    pub async fn execute_async(&self) -> Result<ApiResponse, Auth0Error> {
        self.to_request()?
            .bearer(self.api.access_token()?)
            .execute_async_with_options(&self.api.client_options)
            .await
    }

    pub async fn execute_json_async<T: DeserializeOwned>(&self) -> Result<T, Auth0Error> {
        self.execute_async().await?.json()
    }

    pub fn paginate<T: DeserializeOwned>(&self, per_page: usize) -> SyncPagingIterable<'a, T> {
        SyncPagingIterable {
            request: self.clone(),
            page: 0,
            per_page,
            buffer: Vec::new(),
            done: false,
            marker: PhantomData,
        }
    }

    fn to_request(&self) -> Result<Request, Auth0Error> {
        let path = render_path(self.endpoint.path, &self.path_params)?;
        let url = self.api.base_url.join(path.trim_start_matches('/'))?;
        let mut request =
            Request::new(self.endpoint.method, url).header("Accept", "application/json");
        for (key, value) in &self.headers {
            request = request.header(key, value);
        }
        for (key, value) in &self.query {
            request = request.query(key, value);
        }
        if let Some(value) = &self.body {
            request = match value {
                Body::Empty => request,
                Body::Json(value) => request.json(value.clone()),
                Body::Form(value) => request.form(value.clone()),
                Body::Multipart(value) => request.multipart(value.clone()),
                Body::Text(value) => request.text(value.clone()),
            };
        }
        request = request.options(self.options.clone());
        Ok(request)
    }
}

pub struct SyncPagingIterable<'a, T> {
    request: ManagementRequest<'a>,
    page: usize,
    per_page: usize,
    buffer: Vec<T>,
    done: bool,
    marker: PhantomData<T>,
}

impl<'a, T: DeserializeOwned> Iterator for SyncPagingIterable<'a, T> {
    type Item = Result<T, Auth0Error>;

    fn next(&mut self) -> Option<Self::Item> {
        if let Some(value) = self.buffer.pop() {
            return Some(Ok(value));
        }
        if self.done {
            return None;
        }
        let request = self
            .request
            .clone()
            .query("page", self.page.to_string())
            .query("per_page", self.per_page.to_string())
            .query("include_totals", "true");
        self.page += 1;
        match request.execute() {
            Ok(response) => match page_items::<T>(response.body) {
                Ok(mut values) => {
                    if values.is_empty() {
                        self.done = true;
                        None
                    } else {
                        values.reverse();
                        self.buffer = values;
                        self.next()
                    }
                }
                Err(value) => {
                    self.done = true;
                    Some(Err(value))
                }
            },
            Err(value) => {
                self.done = true;
                Some(Err(value))
            }
        }
    }
}

fn page_items<T: DeserializeOwned>(body: Option<Value>) -> Result<Vec<T>, Auth0Error> {
    let Some(value) = body else {
        return Ok(Vec::new());
    };
    if value.is_array() {
        return Ok(serde_json::from_value(value)?);
    }
    if let Some(items) = value.get("items").or_else(|| value.get("data")) {
        if items.is_array() {
            return Ok(serde_json::from_value(items.clone())?);
        }
    }
    if let Some(object) = value.as_object() {
        for value in object.values() {
            if value.is_array() {
                return Ok(serde_json::from_value(value.clone())?);
            }
        }
    }
    Ok(Vec::new())
}

trait IntoOk<T> {
    fn into_ok(self) -> Result<T, Auth0Error>;
}

impl<T> IntoOk<T> for T {
    fn into_ok(self) -> Result<T, Auth0Error> {
        Ok(self)
    }
}

fn render_path(template: &str, values: &HashMap<String, String>) -> Result<String, Auth0Error> {
    let mut result = String::new();
    let mut rest = template;
    while let Some(start) = rest.find('{') {
        let before = &rest[..start];
        result.push_str(before);
        let after = &rest[start + 1..];
        let end = after
            .find('}')
            .ok_or_else(|| Auth0Error::InvalidInput(template.to_string()))?;
        let key = &after[..end];
        let value = values
            .get(key)
            .ok_or_else(|| Auth0Error::MissingPathParameter(key.to_string()))?;
        result
            .push_str(&url::form_urlencoded::byte_serialize(value.as_bytes()).collect::<String>());
        rest = &after[end + 1..];
    }
    result.push_str(rest);
    Ok(result)
}

impl Endpoint {
    pub fn request<'a>(&'static self, api: &'a ManagementApi) -> ManagementRequest<'a> {
        ManagementRequest::new(api, self)
    }

    pub fn allows_query(&self, key: &str) -> bool {
        self.query.iter().any(|value| *value == key)
    }

    pub fn is_write(&self) -> bool {
        matches!(
            self.method,
            Method::Post | Method::Put | Method::Patch | Method::Delete
        )
    }
}
