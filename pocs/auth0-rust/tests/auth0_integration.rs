use auth0_rust::AuthApi;

fn env(name: &str) -> String {
    std::env::var(name).unwrap_or_else(|_| panic!("{name} must be set"))
}

#[test]
#[ignore = "requires auth0-env.sh and network access"]
fn live_tenant_accepts_authorization_request_built_by_sdk() {
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

    assert!(url.starts_with(&format!("https://{domain}/authorize")));
    assert!(url.contains(&format!("client_id={client_id}")));

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
}
