use auth0_rust::{Auth0Error, AuthApi};

fn env(name: &str) -> String {
    std::env::var(name).unwrap_or_else(|_| panic!("{name} must be set"))
}

#[test]
#[ignore = "requires auth0-env.sh and network access"]
fn live_tenant_accepts_the_credentials_and_the_authorization_request() {
    let domain = env("AUTH0_DOMAIN");
    let client_id = env("AUTH0_CLIENT_ID");
    let redirect_uri = env("AUTH0_REDIRECT_URI");

    let auth = AuthApi::builder(&domain, &client_id)
        .client_secret(env("AUTH0_CLIENT_SECRET"))
        .build()
        .expect("Auth API client");

    let url = auth
        .authorize_url(&redirect_uri)
        .scope("openid profile email")
        .state("integration-state")
        .build();

    let client = reqwest::blocking::Client::builder()
        .redirect(reqwest::redirect::Policy::none())
        .build()
        .expect("http client");
    let response = client.get(&url).send().expect("authorize response");

    let status = response.status().as_u16();
    let location = response
        .headers()
        .get("location")
        .and_then(|value| value.to_str().ok())
        .unwrap_or_default()
        .to_string();

    assert_eq!(302, status, "tenant rejected the authorization request: {location}");
    assert!(
        location.starts_with("/u/login"),
        "expected redirect to Universal Login, got {location}"
    );

    let error = auth
        .exchange_code("integration-invalid-code", &redirect_uri)
        .execute()
        .expect_err("a bogus authorization code must be rejected");

    let error = match error {
        Auth0Error::Http(error) => error,
        other => panic!("expected an Auth0 HTTP error, got {other}"),
    };

    assert_ne!(
        401, error.status,
        "Auth0 refused the client secret, so the credentials are not valid: {}",
        error.body
    );
    assert!(
        error.body.contains("invalid_grant"),
        "expected Auth0 to authenticate the client and reject only the code, got {}",
        error.body
    );
}
