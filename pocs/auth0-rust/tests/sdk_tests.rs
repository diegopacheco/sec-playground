use auth0_rust::{
    ApiErrorKind, AuthApi, Body, MANAGEMENT_ENDPOINTS, ManagementApi, Method, MultipartPart,
    PasswordlessEmailType, RequestOptions, TokenResponse, UserInfo,
};
use serde_json::json;
use std::collections::HashSet;

#[test]
fn auth_domain_defaults_to_https() {
    let auth = AuthApi::builder("tenant.auth0.com", "client")
        .build()
        .unwrap();
    assert_eq!(auth.base_url().as_str(), "https://tenant.auth0.com/");
}

#[test]
fn auth_rejects_insecure_remote_domains() {
    let result = AuthApi::builder("http://tenant.auth0.com", "client").build();

    assert!(result.is_err());
}

#[test]
fn authorize_url_contains_expected_parameters() {
    let auth = AuthApi::builder("tenant.auth0.com", "client")
        .build()
        .unwrap();
    let url = auth
        .authorize_url("https://app.local/callback")
        .scope("openid profile")
        .audience("https://api.local")
        .state("state-1")
        .build();
    assert!(url.starts_with("https://tenant.auth0.com/authorize?"));
    assert!(url.contains("client_id=client"));
    assert!(url.contains("redirect_uri=https%3A%2F%2Fapp.local%2Fcallback"));
    assert!(url.contains("scope=openid+profile"));
    assert!(url.contains("audience=https%3A%2F%2Fapi.local"));
    assert!(url.contains("state=state-1"));
}

#[test]
fn passwordless_email_builds_json_body() {
    let auth = AuthApi::builder("tenant.auth0.com", "client")
        .build()
        .unwrap();
    let request = auth
        .passwordless_email(
            "user@app.local",
            "email",
            PasswordlessEmailType::Code,
            Some(json!({"scope": "openid"})),
        )
        .prepared();
    assert_eq!(request.method, Method::Post);
    assert_eq!(request.url, "https://tenant.auth0.com/passwordless/start");
    assert_eq!(
        request.body,
        Body::Json(json!({
            "client_id": "client",
            "email": "user@app.local",
            "connection": "email",
            "send": "code",
            "authParams": {"scope": "openid"}
        }))
    );
}

#[test]
fn mfa_token_exchanges_build_expected_forms() {
    let auth = AuthApi::builder("tenant.auth0.com", "client")
        .client_secret("secret")
        .build()
        .unwrap();
    let otp = auth.exchange_mfa_otp("mfa-token", "123456").prepared();
    let recovery = auth
        .exchange_mfa_recovery_code("mfa-token", "recovery-code")
        .prepared();

    assert_eq!(otp.method, Method::Post);
    assert!(
        matches!(otp.body, Body::Form(ref values) if values.contains(&(
            "grant_type".into(),
            "http://auth0.com/oauth/grant-type/mfa-otp".into()
        )))
    );
    assert!(
        matches!(recovery.body, Body::Form(ref values) if values.contains(&(
            "recovery_code".into(),
            "recovery-code".into()
        )))
    );
}

#[test]
fn client_assertion_authentication_does_not_send_a_client_secret() {
    let auth = AuthApi::builder("tenant.auth0.com", "client")
        .client_secret("secret")
        .client_assertion("assertion")
        .build()
        .unwrap();
    let request = auth.request_token("https://api.local").prepared();

    let Body::Form(values) = request.body else {
        panic!("expected form");
    };
    assert!(
        values
            .iter()
            .any(|(key, value)| key == "client_assertion" && value == "assertion")
    );
    assert!(
        !values.iter().any(|(key, _)| key == "client_secret"),
        "{values:?}"
    );
}

#[test]
fn management_catalog_has_java_routes() {
    assert!(MANAGEMENT_ENDPOINTS.len() >= 400);
    assert!(
        MANAGEMENT_ENDPOINTS
            .iter()
            .any(|endpoint| endpoint.name == "clients.list")
    );
    assert!(
        MANAGEMENT_ENDPOINTS
            .iter()
            .any(|endpoint| endpoint.name == "users.get")
    );
    assert!(
        MANAGEMENT_ENDPOINTS
            .iter()
            .any(|endpoint| endpoint.name == "roles.permissions.list")
    );
    assert!(
        MANAGEMENT_ENDPOINTS
            .iter()
            .any(|endpoint| endpoint.name == "organizations.members.roles.assign")
    );
}

#[test]
fn management_catalog_has_unique_names() {
    let mut names = HashSet::new();
    for endpoint in MANAGEMENT_ENDPOINTS {
        assert!(names.insert(endpoint.name));
    }
}

#[test]
fn management_request_renders_path_query_and_token() {
    let api = ManagementApi::builder()
        .domain("tenant.auth0.com")
        .token("token")
        .build()
        .unwrap();
    let request = api
        .users()
        .request("get")
        .unwrap()
        .path_param("id", "auth0|abc")
        .query("fields", "email,name")
        .prepared()
        .unwrap();
    assert_eq!(request.method, Method::Get);
    assert_eq!(
        request.url,
        "https://tenant.auth0.com/api/v2/users/auth0%7Cabc?fields=email%2Cname"
    );
    assert!(
        request
            .headers
            .contains(&("Authorization".into(), "Bearer token".into()))
    );
}

#[test]
fn generated_management_methods_build_requests() {
    let api = ManagementApi::builder()
        .domain("tenant.auth0.com")
        .token("token")
        .build()
        .unwrap();
    let request = api.clients_list().unwrap().prepared().unwrap();
    assert_eq!(request.method, Method::Get);
    assert_eq!(request.url, "https://tenant.auth0.com/api/v2/clients");
    let raw = api.raw().clients_list().unwrap().prepared().unwrap();
    assert_eq!(raw.method, Method::Get);
    assert_eq!(raw.url, request.url);
}

#[test]
fn request_options_are_applied() {
    let api = ManagementApi::builder()
        .domain("tenant.auth0.com")
        .token("token")
        .build()
        .unwrap();
    let request = api
        .clients_get()
        .unwrap()
        .path_param("id", "abc")
        .options(
            RequestOptions::new()
                .query("fields", "name")
                .header("Auth0-Client", "sdk"),
        )
        .prepared()
        .unwrap();
    assert_eq!(
        request.url,
        "https://tenant.auth0.com/api/v2/clients/abc?fields=name"
    );
    assert!(
        request
            .headers
            .contains(&("Auth0-Client".into(), "sdk".into()))
    );
}

#[test]
fn multipart_request_is_preserved() {
    let api = ManagementApi::builder()
        .domain("tenant.auth0.com")
        .token("token")
        .build()
        .unwrap();
    let request = api
        .jobs_users_imports_create()
        .unwrap()
        .multipart(vec![MultipartPart::file(
            "users",
            "users.json",
            b"[]".to_vec(),
        )])
        .prepared()
        .unwrap();
    assert_eq!(request.method, Method::Post);
    match request.body {
        Body::Multipart(parts) => {
            assert_eq!(parts.len(), 1);
            assert_eq!(parts[0].name, "users");
            assert_eq!(parts[0].file_name.as_deref(), Some("users.json"));
        }
        _ => panic!("expected multipart"),
    }
}

#[test]
fn java_auth_fixtures_decode_into_typed_structs() {
    let token: TokenResponse = serde_json::from_str(
        r#"{
          "id_token": "eyJ0eXAiOiJKV1Qi...",
          "access_token": "A9CvPwFojaBI...",
          "refresh_token": "GEbRxBN...edjnXbL",
          "token_type": "bearer",
          "expires_in": 86000
        }"#,
    )
    .unwrap();
    assert_eq!(token.access_token.as_deref(), Some("A9CvPwFojaBI..."));
    assert_eq!(token.expires_in, Some(86000));
    let user: UserInfo = serde_json::from_str(
        r#"{
          "email_verified": false,
          "email": "test.account@userinfo.com",
          "clientID": "q2hnj2iu...",
          "updated_at": "2016-12-05T15:15:40.545Z",
          "name": "test.account@userinfo.com",
          "picture": "https://s.gravatar.com/avatar/dummy.png",
          "user_id": "auth0|58454...",
          "nickname": "test.account",
          "identities": [
            {
              "user_id": "58454...",
              "provider": "auth0",
              "connection": "Username-Password-Authentication",
              "isSocial": false
            }
          ],
          "created_at": "2016-12-05T11:16:59.640Z",
          "sub": "auth0|58454..."
        }"#,
    )
    .unwrap();
    assert_eq!(user.email.as_deref(), Some("test.account@userinfo.com"));
    assert!(user.extra.contains_key("identities"));
}

#[test]
fn api_error_maps_status_structs() {
    let error = auth0_rust::ApiError {
        status: 429,
        code: Some("too_many_requests".into()),
        message: Some("slow down".into()),
        body: "{}".into(),
    };
    assert!(matches!(error.kind(), ApiErrorKind::TooManyRequests(_)));
}
