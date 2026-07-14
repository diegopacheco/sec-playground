use crate::assertion::RsaClientAssertionSigner;
use crate::client::ClientOptions;
use crate::error::Auth0Error;
use crate::request::{ApiResponse, Method, Request};
use serde::de::DeserializeOwned;
use serde_json::{Value, json};
use url::Url;

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum PasswordlessEmailType {
    Code,
    Link,
}

impl PasswordlessEmailType {
    pub fn as_str(self) -> &'static str {
        match self {
            PasswordlessEmailType::Code => "code",
            PasswordlessEmailType::Link => "link",
        }
    }
}

#[derive(Clone)]
pub struct AuthApiBuilder {
    domain: String,
    client_id: String,
    client_secret: Option<String>,
    client_assertion: Option<String>,
    client_assertion_signer: Option<RsaClientAssertionSigner>,
    client_options: ClientOptions,
}

impl AuthApiBuilder {
    pub fn new(domain: impl Into<String>, client_id: impl Into<String>) -> Self {
        Self {
            domain: domain.into(),
            client_id: client_id.into(),
            client_secret: None,
            client_assertion: None,
            client_assertion_signer: None,
            client_options: ClientOptions::default(),
        }
    }

    pub fn client_secret(mut self, value: impl Into<String>) -> Self {
        self.client_secret = Some(value.into());
        self
    }

    pub fn client_assertion(mut self, value: impl Into<String>) -> Self {
        self.client_assertion = Some(value.into());
        self
    }

    pub fn client_assertion_signer(mut self, value: RsaClientAssertionSigner) -> Self {
        self.client_assertion_signer = Some(value);
        self
    }

    pub fn client_options(mut self, value: ClientOptions) -> Self {
        self.client_options = value;
        self
    }

    pub fn build(self) -> Result<AuthApi, Auth0Error> {
        let base_url = normalize_domain(&self.domain)?;
        AuthApi::new_with_signer(
            base_url,
            self.client_id,
            self.client_secret,
            self.client_assertion,
            self.client_assertion_signer,
            self.client_options,
        )
    }
}

#[derive(Clone)]
pub struct AuthApi {
    base_url: Url,
    client_id: String,
    client_secret: Option<String>,
    client_assertion: Option<String>,
    client_assertion_signer: Option<RsaClientAssertionSigner>,
    client_options: ClientOptions,
}

impl AuthApi {
    pub fn builder(domain: impl Into<String>, client_id: impl Into<String>) -> AuthApiBuilder {
        AuthApiBuilder::new(domain, client_id)
    }

    pub fn new(
        domain: impl Into<String>,
        client_id: impl Into<String>,
        client_secret: Option<String>,
        client_assertion: Option<String>,
    ) -> Result<Self, Auth0Error> {
        Self::new_with_options(
            normalize_domain(&domain.into())?,
            client_id.into(),
            client_secret,
            client_assertion,
            ClientOptions::default(),
        )
    }

    pub fn new_with_options(
        base_url: Url,
        client_id: String,
        client_secret: Option<String>,
        client_assertion: Option<String>,
        client_options: ClientOptions,
    ) -> Result<Self, Auth0Error> {
        Self::new_with_signer(
            base_url,
            client_id,
            client_secret,
            client_assertion,
            None,
            client_options,
        )
    }

    fn new_with_signer(
        base_url: Url,
        client_id: String,
        client_secret: Option<String>,
        client_assertion: Option<String>,
        client_assertion_signer: Option<RsaClientAssertionSigner>,
        client_options: ClientOptions,
    ) -> Result<Self, Auth0Error> {
        Ok(Self {
            base_url,
            client_id,
            client_secret,
            client_assertion,
            client_assertion_signer,
            client_options,
        })
    }

    pub fn base_url(&self) -> &Url {
        &self.base_url
    }

    pub fn authorize_url(&self, redirect_uri: impl Into<String>) -> AuthorizeUrlBuilder {
        AuthorizeUrlBuilder::new(
            self.base_url.clone(),
            self.client_id.clone(),
            redirect_uri.into(),
        )
    }

    pub fn authorize_url_with_par(&self, request_uri: impl AsRef<str>) -> String {
        let mut url = self.base_url.join("authorize").unwrap();
        url.query_pairs_mut()
            .append_pair("client_id", &self.client_id)
            .append_pair("request_uri", request_uri.as_ref());
        url.to_string()
    }

    pub fn authorize_url_with_jar(&self, request: impl AsRef<str>) -> String {
        let mut url = self.base_url.join("authorize").unwrap();
        url.query_pairs_mut()
            .append_pair("client_id", &self.client_id)
            .append_pair("request", request.as_ref());
        url.to_string()
    }

    pub fn logout_url(&self) -> LogoutUrlBuilder {
        LogoutUrlBuilder::new(self.base_url.clone(), self.client_id.clone())
    }

    pub fn logout_url_with_return_to(
        &self,
        return_to: impl Into<String>,
        set_client_id: bool,
    ) -> LogoutUrlBuilder {
        LogoutUrlBuilder::new_with_return_to(
            self.base_url.clone(),
            self.client_id.clone(),
            return_to.into(),
            set_client_id,
        )
    }

    pub fn request_token(&self, audience: impl Into<String>) -> Request {
        let mut form = self.client_auth_form();
        form.push(("grant_type".into(), "client_credentials".into()));
        form.push(("audience".into(), audience.into()));
        self.form_request(Method::Post, "oauth/token", form)
    }

    pub fn request_token_for_organization(
        &self,
        audience: impl Into<String>,
        organization: impl Into<String>,
    ) -> Request {
        self.request_token(audience)
            .form_append("organization", organization)
    }

    pub fn exchange_code(
        &self,
        code: impl Into<String>,
        redirect_uri: impl Into<String>,
    ) -> Request {
        let mut form = self.client_auth_form();
        form.push(("grant_type".into(), "authorization_code".into()));
        form.push(("code".into(), code.into()));
        form.push(("redirect_uri".into(), redirect_uri.into()));
        self.form_request(Method::Post, "oauth/token", form)
    }

    pub fn exchange_code_with_verifier(
        &self,
        code: impl Into<String>,
        verifier: impl Into<String>,
        redirect_uri: impl Into<String>,
    ) -> Request {
        self.exchange_code(code, redirect_uri)
            .form_append("code_verifier", verifier)
    }

    pub fn refresh_token(&self, refresh_token: impl Into<String>) -> Request {
        let mut form = self.client_auth_form();
        form.push(("grant_type".into(), "refresh_token".into()));
        form.push(("refresh_token".into(), refresh_token.into()));
        self.form_request(Method::Post, "oauth/token", form)
    }

    pub fn renew_auth(&self, refresh_token: impl Into<String>) -> Request {
        self.refresh_token(refresh_token)
    }

    pub fn login(
        &self,
        username: impl Into<String>,
        password: impl Into<String>,
        realm: Option<String>,
        audience: Option<String>,
        scope: Option<String>,
    ) -> Request {
        let mut form = self.client_auth_form();
        if realm.is_some() {
            form.push((
                "grant_type".into(),
                "http://auth0.com/oauth/grant-type/password-realm".into(),
            ));
        } else {
            form.push(("grant_type".into(), "password".into()));
        }
        form.push(("username".into(), username.into()));
        form.push(("password".into(), password.into()));
        if let Some(value) = realm {
            form.push(("realm".into(), value));
        }
        if let Some(value) = audience {
            form.push(("audience".into(), value));
        }
        if let Some(value) = scope {
            form.push(("scope".into(), value));
        }
        self.form_request(Method::Post, "oauth/token", form)
    }

    pub fn login_password(
        &self,
        username: impl Into<String>,
        password: impl Into<String>,
    ) -> Request {
        self.login(username, password, None, None, None)
    }

    pub fn login_realm(
        &self,
        username: impl Into<String>,
        password: impl Into<String>,
        realm: impl Into<String>,
    ) -> Request {
        self.login(username, password, Some(realm.into()), None, None)
    }

    pub fn exchange_passwordless_otp(
        &self,
        email_or_phone: impl Into<String>,
        realm: impl Into<String>,
        otp: impl Into<String>,
    ) -> Request {
        let mut form = self.client_auth_form();
        form.push((
            "grant_type".into(),
            "http://auth0.com/oauth/grant-type/passwordless/otp".into(),
        ));
        form.push(("username".into(), email_or_phone.into()));
        form.push(("realm".into(), realm.into()));
        form.push(("otp".into(), otp.into()));
        self.form_request(Method::Post, "oauth/token", form)
    }

    pub fn exchange_token(
        &self,
        subject_token: impl Into<String>,
        subject_token_type: impl Into<String>,
    ) -> Request {
        let mut form = self.client_auth_form();
        form.push((
            "grant_type".into(),
            "urn:ietf:params:oauth:grant-type:token-exchange".into(),
        ));
        form.push(("subject_token".into(), subject_token.into()));
        form.push(("subject_token_type".into(), subject_token_type.into()));
        self.form_request(Method::Post, "oauth/token", form)
    }

    pub fn revoke_token(&self, refresh_token: impl Into<String>) -> Request {
        let mut form = self.client_auth_form();
        form.push(("token".into(), refresh_token.into()));
        self.form_request(Method::Post, "oauth/revoke", form)
    }

    pub fn user_info(&self, access_token: impl AsRef<str>) -> Request {
        self.request(Method::Get, "userinfo")
            .bearer(access_token.as_ref())
    }

    pub fn jwks(&self) -> Request {
        self.request(Method::Get, ".well-known/jwks.json")
    }

    pub fn reset_password(
        &self,
        email: impl Into<String>,
        connection: impl Into<String>,
        organization: Option<String>,
    ) -> Request {
        let mut body = json!({
            "client_id": self.client_id,
            "email": email.into(),
            "connection": connection.into()
        });
        if let Some(value) = organization {
            body["organization"] = Value::String(value);
        }
        self.request(Method::Post, "dbconnections/change_password")
            .json(body)
    }

    pub fn reset_password_with_client_id(
        &self,
        client_id: impl Into<String>,
        email: impl Into<String>,
        connection: impl Into<String>,
    ) -> Request {
        self.request(Method::Post, "dbconnections/change_password")
            .json(json!({
                "client_id": client_id.into(),
                "email": email.into(),
                "connection": connection.into()
            }))
    }

    pub fn sign_up(
        &self,
        email: impl Into<String>,
        password: impl Into<String>,
        connection: impl Into<String>,
        username: Option<String>,
        user_metadata: Option<Value>,
    ) -> Request {
        let mut body = json!({
            "client_id": self.client_id,
            "email": email.into(),
            "password": password.into(),
            "connection": connection.into()
        });
        if let Some(value) = username {
            body["username"] = Value::String(value);
        }
        if let Some(value) = user_metadata {
            body["user_metadata"] = value;
        }
        self.request(Method::Post, "dbconnections/signup")
            .json(body)
    }

    pub fn sign_up_with_username(
        &self,
        email: impl Into<String>,
        username: impl Into<String>,
        password: impl Into<String>,
        connection: impl Into<String>,
    ) -> Request {
        self.sign_up(email, password, connection, Some(username.into()), None)
    }

    pub fn sign_up_with_phone_number(
        &self,
        email: impl Into<String>,
        username: impl Into<String>,
        password: impl Into<String>,
        connection: impl Into<String>,
        phone_number: impl Into<String>,
    ) -> Request {
        self.sign_up_with_username(email, username, password, connection)
            .json_insert("phone_number", Value::String(phone_number.into()))
    }

    pub fn passwordless_email(
        &self,
        email: impl Into<String>,
        connection: impl Into<String>,
        send: PasswordlessEmailType,
        auth_params: Option<Value>,
    ) -> Request {
        let mut body = json!({
            "client_id": self.client_id,
            "email": email.into(),
            "connection": connection.into(),
            "send": send.as_str()
        });
        if let Some(value) = auth_params {
            body["authParams"] = value;
        }
        self.request(Method::Post, "passwordless/start").json(body)
    }

    pub fn start_passwordless_email_flow(
        &self,
        email: impl Into<String>,
        send: PasswordlessEmailType,
    ) -> Request {
        self.passwordless_email(email, "email", send, None)
    }

    pub fn passwordless_sms(
        &self,
        phone_number: impl Into<String>,
        connection: impl Into<String>,
        auth_params: Option<Value>,
    ) -> Request {
        let mut body = json!({
            "client_id": self.client_id,
            "phone_number": phone_number.into(),
            "connection": connection.into()
        });
        if let Some(value) = auth_params {
            body["authParams"] = value;
        }
        self.request(Method::Post, "passwordless/start").json(body)
    }

    pub fn start_passwordless_sms_flow(&self, phone_number: impl Into<String>) -> Request {
        self.passwordless_sms(phone_number, "sms", None)
    }

    pub fn pushed_authorization_request(
        &self,
        redirect_uri: impl Into<String>,
        response_type: impl Into<String>,
        parameters: Vec<(String, String)>,
    ) -> Request {
        let mut form = self.client_auth_form();
        form.push(("redirect_uri".into(), redirect_uri.into()));
        form.push(("response_type".into(), response_type.into()));
        for value in parameters {
            form.push(value);
        }
        self.form_request(Method::Post, "oauth/par", form)
    }

    pub fn pushed_authorization_request_with_details(
        &self,
        redirect_uri: impl Into<String>,
        response_type: impl Into<String>,
        parameters: Vec<(String, String)>,
        authorization_details: Value,
    ) -> Request {
        self.pushed_authorization_request(redirect_uri, response_type, parameters)
            .form_append("authorization_details", authorization_details.to_string())
    }

    pub fn pushed_authorization_request_with_jar(&self, request: impl Into<String>) -> Request {
        let mut form = self.client_auth_form();
        form.push(("request".into(), request.into()));
        self.form_request(Method::Post, "oauth/par", form)
    }

    pub fn pushed_authorization_request_with_jar_and_details(
        &self,
        request: impl Into<String>,
        authorization_details: Value,
    ) -> Request {
        self.pushed_authorization_request_with_jar(request)
            .form_append("authorization_details", authorization_details.to_string())
    }

    pub fn authorize_back_channel(
        &self,
        scope: impl Into<String>,
        binding_message: impl Into<String>,
        login_hint: Value,
        audience: Option<String>,
        requested_expiry: Option<u32>,
    ) -> Request {
        let mut form = self.client_auth_form();
        form.push(("scope".into(), scope.into()));
        form.push(("binding_message".into(), binding_message.into()));
        form.push(("login_hint".into(), login_hint.to_string()));
        if let Some(value) = audience {
            form.push(("audience".into(), value));
        }
        if let Some(value) = requested_expiry {
            form.push(("requested_expiry".into(), value.to_string()));
        }
        self.form_request(Method::Post, "bc-authorize", form)
    }

    pub fn back_channel_login_status(
        &self,
        auth_req_id: impl Into<String>,
        grant_type: impl Into<String>,
    ) -> Request {
        let mut form = self.client_auth_form();
        form.push(("auth_req_id".into(), auth_req_id.into()));
        form.push(("grant_type".into(), grant_type.into()));
        self.form_request(Method::Post, "oauth/token", form)
    }

    pub fn mfa_challenge(
        &self,
        mfa_token: impl Into<String>,
        challenge_type: impl Into<String>,
    ) -> Request {
        self.request(Method::Post, "mfa/challenge").json(json!({
            "client_id": self.client_id,
            "challenge_type": challenge_type.into(),
            "mfa_token": mfa_token.into()
        }))
    }

    pub fn mfa_challenge_request(
        &self,
        mfa_token: impl Into<String>,
        challenge_type: impl Into<String>,
    ) -> Request {
        self.mfa_challenge(mfa_token, challenge_type)
    }

    pub fn mfa_challenge_with_authenticator(
        &self,
        mfa_token: impl Into<String>,
        challenge_type: impl Into<String>,
        authenticator_id: impl Into<String>,
    ) -> Request {
        self.mfa_challenge(mfa_token, challenge_type)
            .json_insert("authenticator_id", Value::String(authenticator_id.into()))
    }

    pub fn exchange_mfa_otp(
        &self,
        mfa_token: impl Into<String>,
        otp: impl Into<String>,
    ) -> Request {
        let mut form = self.client_auth_form();
        form.push((
            "grant_type".into(),
            "http://auth0.com/oauth/grant-type/mfa-otp".into(),
        ));
        form.push(("mfa_token".into(), mfa_token.into()));
        form.push(("otp".into(), otp.into()));
        self.form_request(Method::Post, "oauth/token", form)
    }

    pub fn exchange_mfa_oob(
        &self,
        mfa_token: impl Into<String>,
        oob_code: impl Into<String>,
        binding_code: impl Into<String>,
    ) -> Request {
        let mut form = self.client_auth_form();
        form.push((
            "grant_type".into(),
            "http://auth0.com/oauth/grant-type/mfa-oob".into(),
        ));
        form.push(("mfa_token".into(), mfa_token.into()));
        form.push(("oob_code".into(), oob_code.into()));
        form.push(("binding_code".into(), binding_code.into()));
        self.form_request(Method::Post, "oauth/token", form)
    }

    pub fn exchange_mfa_recovery_code(
        &self,
        mfa_token: impl Into<String>,
        recovery_code: impl Into<String>,
    ) -> Request {
        let mut form = self.client_auth_form();
        form.push((
            "grant_type".into(),
            "http://auth0.com/oauth/grant-type/mfa-recovery-code".into(),
        ));
        form.push(("mfa_token".into(), mfa_token.into()));
        form.push(("recovery_code".into(), recovery_code.into()));
        self.form_request(Method::Post, "oauth/token", form)
    }

    pub fn add_otp_authenticator(&self, mfa_token: impl Into<String>) -> Request {
        self.request(Method::Post, "mfa/associate").json(json!({
            "client_id": self.client_id,
            "authenticator_types": ["otp"],
            "mfa_token": mfa_token.into()
        }))
    }

    pub fn add_oob_authenticator(
        &self,
        mfa_token: impl Into<String>,
        authenticator_type: impl Into<String>,
    ) -> Request {
        self.request(Method::Post, "mfa/associate").json(json!({
            "client_id": self.client_id,
            "authenticator_types": [authenticator_type.into()],
            "mfa_token": mfa_token.into()
        }))
    }

    pub fn list_authenticators(&self, access_token: impl AsRef<str>) -> Request {
        self.request(Method::Get, "mfa/authenticators")
            .bearer(access_token.as_ref())
    }

    pub fn delete_authenticator(
        &self,
        access_token: impl AsRef<str>,
        authenticator_id: impl AsRef<str>,
    ) -> Request {
        let path = format!("mfa/authenticators/{}", encode(authenticator_id.as_ref()));
        self.request(Method::Delete, path)
            .bearer(access_token.as_ref())
    }

    pub fn execute(&self, request: &Request) -> Result<ApiResponse, Auth0Error> {
        request.execute_with_options(&self.client_options)
    }

    pub fn execute_json<T: DeserializeOwned>(&self, request: &Request) -> Result<T, Auth0Error> {
        request.execute_json_with_options(&self.client_options)
    }

    pub async fn execute_async(&self, request: &Request) -> Result<ApiResponse, Auth0Error> {
        request
            .execute_async_with_options(&self.client_options)
            .await
    }

    pub async fn execute_json_async<T: DeserializeOwned>(
        &self,
        request: &Request,
    ) -> Result<T, Auth0Error> {
        request
            .execute_json_async_with_options(&self.client_options)
            .await
    }

    fn request(&self, method: Method, path: impl AsRef<str>) -> Request {
        Request::new(method, self.base_url.join(path.as_ref()).unwrap())
    }

    fn form_request(
        &self,
        method: Method,
        path: impl AsRef<str>,
        form: Vec<(String, String)>,
    ) -> Request {
        let request = self.request(method, path).form(form);
        if self.client_assertion.is_none()
            && let Some(signer) = &self.client_assertion_signer
        {
            return request.client_assertion(
                signer.clone(),
                self.client_id.clone(),
                self.client_id.clone(),
            );
        }
        request
    }

    fn client_auth_form(&self) -> Vec<(String, String)> {
        let mut form = vec![("client_id".into(), self.client_id.clone())];
        let client_secret = self
            .client_secret
            .as_ref()
            .filter(|_| self.client_assertion.is_none() && self.client_assertion_signer.is_none());
        if let Some(value) = client_secret {
            form.push(("client_secret".into(), value.clone()));
        }
        if let Some(value) = &self.client_assertion {
            form.push(("client_assertion".into(), value.clone()));
            form.push((
                "client_assertion_type".into(),
                "urn:ietf:params:oauth:client-assertion-type:jwt-bearer".into(),
            ));
        }
        form
    }
}

#[derive(Debug, Clone)]
pub struct AuthorizeUrlBuilder {
    base_url: Url,
    params: Vec<(String, String)>,
}

impl AuthorizeUrlBuilder {
    pub fn new(base_url: Url, client_id: String, redirect_uri: String) -> Self {
        Self {
            base_url,
            params: vec![
                ("client_id".into(), client_id),
                ("redirect_uri".into(), redirect_uri),
                ("response_type".into(), "code".into()),
            ],
        }
    }

    pub fn audience(self, value: impl Into<String>) -> Self {
        self.param("audience", value)
    }

    pub fn scope(self, value: impl Into<String>) -> Self {
        self.param("scope", value)
    }

    pub fn state(self, value: impl Into<String>) -> Self {
        self.param("state", value)
    }

    pub fn connection(self, value: impl Into<String>) -> Self {
        self.param("connection", value)
    }

    pub fn organization(self, value: impl Into<String>) -> Self {
        self.param("organization", value)
    }

    pub fn invitation(self, value: impl Into<String>) -> Self {
        self.param("invitation", value)
    }

    pub fn prompt(self, value: impl Into<String>) -> Self {
        self.param("prompt", value)
    }

    pub fn nonce(self, value: impl Into<String>) -> Self {
        self.param("nonce", value)
    }

    pub fn code_challenge(self, value: impl Into<String>) -> Self {
        self.param("code_challenge", value)
            .param("code_challenge_method", "S256")
    }

    pub fn response_type(mut self, value: impl Into<String>) -> Self {
        self.params.retain(|(key, _)| key != "response_type");
        self.params.push(("response_type".into(), value.into()));
        self
    }

    pub fn param(mut self, key: impl Into<String>, value: impl Into<String>) -> Self {
        self.params.push((key.into(), value.into()));
        self
    }

    pub fn build(self) -> String {
        let mut url = self.base_url.join("authorize").unwrap();
        {
            let mut query = url.query_pairs_mut();
            for (key, value) in self.params {
                query.append_pair(&key, &value);
            }
        }
        url.to_string()
    }
}

#[derive(Debug, Clone)]
pub struct LogoutUrlBuilder {
    base_url: Url,
    params: Vec<(String, String)>,
}

impl LogoutUrlBuilder {
    pub fn new(base_url: Url, client_id: String) -> Self {
        Self {
            base_url,
            params: vec![("client_id".into(), client_id)],
        }
    }

    pub fn new_with_return_to(
        base_url: Url,
        client_id: String,
        return_to: String,
        set_client_id: bool,
    ) -> Self {
        let mut params = vec![("returnTo".into(), return_to)];
        if set_client_id {
            params.push(("client_id".into(), client_id));
        }
        Self { base_url, params }
    }

    pub fn return_to(self, value: impl Into<String>) -> Self {
        self.param("returnTo", value)
    }

    pub fn federated(self) -> Self {
        self.use_federated(true)
    }

    pub fn use_federated(mut self, value: bool) -> Self {
        self.params.retain(|(key, _)| key != "federated");
        if value {
            self.params.push(("federated".into(), "".into()));
        }
        self
    }

    pub fn param(mut self, key: impl Into<String>, value: impl Into<String>) -> Self {
        self.params.push((key.into(), value.into()));
        self
    }

    pub fn build(self) -> String {
        let mut url = self.base_url.join("v2/logout").unwrap();
        {
            let mut query = url.query_pairs_mut();
            for (key, value) in self.params {
                query.append_pair(&key, &value);
            }
        }
        url.to_string()
    }
}

pub fn normalize_domain(domain: &str) -> Result<Url, Auth0Error> {
    let value = if domain.starts_with("https://") || domain.starts_with("http://") {
        domain.to_string()
    } else {
        format!("https://{}", domain)
    };
    let mut url = Url::parse(&value)?;
    let host = url
        .host_str()
        .ok_or_else(|| Auth0Error::InvalidInput("domain host".into()))?;
    let loopback =
        host == "localhost" || host == "127.0.0.1" || host == "::1" || host.ends_with(".localhost");
    if url.scheme() != "https" && !(url.scheme() == "http" && loopback) {
        return Err(Auth0Error::InvalidInput("domain must use HTTPS".into()));
    }
    if !url.username().is_empty()
        || url.password().is_some()
        || url.query().is_some()
        || url.fragment().is_some()
    {
        return Err(Auth0Error::InvalidInput(
            "domain contains unsupported URL parts".into(),
        ));
    }
    if url.path() != "/" && !url.path().is_empty() {
        return Err(Auth0Error::InvalidInput(
            "domain must not contain a path".into(),
        ));
    }
    url.set_path("/");
    Ok(url)
}

fn encode(value: &str) -> String {
    url::form_urlencoded::byte_serialize(value.as_bytes()).collect()
}
