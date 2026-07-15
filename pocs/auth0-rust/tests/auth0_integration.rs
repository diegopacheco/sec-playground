use auth0_rust::{Auth0Error, AuthApi, ManagementApi};

fn domain() -> String {
    std::env::var("AUTH0_DOMAIN").expect("AUTH0_DOMAIN must be set")
}

fn status(error: Auth0Error) -> u16 {
    match error {
        Auth0Error::Http(error) => error.status,
        other => panic!("expected Auth0 HTTP error, got {other}"),
    }
}

#[test]
#[ignore = "requires AUTH0_DOMAIN and network access"]
fn retrieves_tenant_signing_keys() {
    let auth = AuthApi::builder(domain(), "integration-client")
        .build()
        .expect("Auth API client");

    let response = auth.jwks().execute().expect("JWKS response");
    let keys = response
        .body
        .and_then(|body| body.get("keys").cloned())
        .and_then(|keys| keys.as_array().cloned())
        .expect("JWKS keys");

    assert_eq!(200, response.status);
    assert!(!keys.is_empty());
}

#[test]
#[ignore = "requires AUTH0_DOMAIN and network access"]
fn rejects_invalid_access_token() {
    let auth = AuthApi::builder(domain(), "integration-client")
        .build()
        .expect("Auth API client");

    let error = auth
        .user_info("integration-invalid-token")
        .execute()
        .expect_err("invalid access token must fail");

    assert_eq!(401, status(error));
}

#[test]
#[ignore = "requires AUTH0_DOMAIN and network access"]
fn rejects_invalid_management_client_credentials() {
    let management = ManagementApi::builder()
        .domain(domain())
        .client_credentials("integration-invalid-client", "integration-invalid-secret")
        .build()
        .expect("Management API client");

    let error = management
        .users_get()
        .expect("users request")
        .path_param("id", "auth0|integration-invalid-user")
        .execute()
        .expect_err("invalid client credentials must fail");

    assert_eq!(401, status(error));
}
