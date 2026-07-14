use crate::error::Auth0Error;
use reqwest::Proxy;
use std::collections::HashSet;
use std::time::Duration;

#[derive(Debug, Clone, PartialEq, Eq)]
pub struct ProxyOptions {
    pub url: String,
    pub basic_auth: Option<(String, String)>,
}

impl ProxyOptions {
    pub fn new(url: impl Into<String>) -> Self {
        Self {
            url: url.into(),
            basic_auth: None,
        }
    }

    pub fn basic_auth(mut self, username: impl Into<String>, password: impl Into<String>) -> Self {
        self.basic_auth = Some((username.into(), password.into()));
        self
    }

    pub fn to_proxy(&self) -> Result<Proxy, Auth0Error> {
        let mut proxy =
            Proxy::all(&self.url).map_err(|value| Auth0Error::Transport(value.to_string()))?;
        if let Some((username, password)) = &self.basic_auth {
            proxy = proxy.basic_auth(username, password);
        }
        Ok(proxy)
    }
}

#[derive(Debug, Clone, PartialEq, Eq)]
pub enum LogLevel {
    None,
    Basic,
    Headers,
    Body,
}

#[derive(Debug, Clone, PartialEq, Eq)]
pub struct LoggingOptions {
    pub level: LogLevel,
    pub headers_to_redact: HashSet<String>,
}

impl LoggingOptions {
    pub fn new(level: LogLevel) -> Self {
        Self {
            level,
            headers_to_redact: HashSet::new(),
        }
    }

    pub fn redact_header(mut self, header: impl Into<String>) -> Self {
        self.headers_to_redact
            .insert(header.into().to_ascii_lowercase());
        self
    }

    pub fn redact_headers(mut self, headers: impl IntoIterator<Item = impl Into<String>>) -> Self {
        for header in headers {
            self.headers_to_redact
                .insert(header.into().to_ascii_lowercase());
        }
        self
    }

    pub fn redact_value(&self, header: &str, value: &str) -> String {
        if self
            .headers_to_redact
            .contains(&header.to_ascii_lowercase())
        {
            "[redacted]".into()
        } else {
            value.into()
        }
    }
}

#[derive(Debug, Clone, PartialEq, Eq)]
pub struct RetryOptions {
    pub max_retries: usize,
    pub retry_statuses: Vec<u16>,
    pub respect_retry_after: bool,
}

impl RetryOptions {
    pub fn new(max_retries: usize) -> Self {
        Self {
            max_retries,
            retry_statuses: vec![429, 500, 502, 503, 504],
            respect_retry_after: true,
        }
    }

    pub fn retry_statuses(mut self, statuses: impl IntoIterator<Item = u16>) -> Self {
        self.retry_statuses = statuses.into_iter().collect();
        self
    }

    pub fn respect_retry_after(mut self, value: bool) -> Self {
        self.respect_retry_after = value;
        self
    }

    pub fn should_retry(&self, status: u16) -> bool {
        self.retry_statuses.contains(&status)
    }
}

#[derive(Clone)]
pub struct ClientOptions {
    pub blocking_client: reqwest::blocking::Client,
    pub async_client: reqwest::Client,
    pub timeout: Option<Duration>,
    pub retry: RetryOptions,
    pub logging: LoggingOptions,
}

impl ClientOptions {
    pub fn builder() -> ClientOptionsBuilder {
        ClientOptionsBuilder::new()
    }
}

impl Default for ClientOptions {
    fn default() -> Self {
        ClientOptionsBuilder::new().build().unwrap()
    }
}

pub struct ClientOptionsBuilder {
    blocking_client: Option<reqwest::blocking::Client>,
    async_client: Option<reqwest::Client>,
    timeout: Option<Duration>,
    proxy: Option<ProxyOptions>,
    retry: RetryOptions,
    logging: LoggingOptions,
}

impl ClientOptionsBuilder {
    pub fn new() -> Self {
        Self {
            blocking_client: None,
            async_client: None,
            timeout: None,
            proxy: None,
            retry: RetryOptions::new(0),
            logging: LoggingOptions::new(LogLevel::None),
        }
    }

    pub fn blocking_client(mut self, client: reqwest::blocking::Client) -> Self {
        self.blocking_client = Some(client);
        self
    }

    pub fn async_client(mut self, client: reqwest::Client) -> Self {
        self.async_client = Some(client);
        self
    }

    pub fn timeout(mut self, value: Duration) -> Self {
        self.timeout = Some(value);
        self
    }

    pub fn proxy(mut self, value: ProxyOptions) -> Self {
        self.proxy = Some(value);
        self
    }

    pub fn retry(mut self, value: RetryOptions) -> Self {
        self.retry = value;
        self
    }

    pub fn logging(mut self, value: LoggingOptions) -> Self {
        self.logging = value;
        self
    }

    pub fn build(self) -> Result<ClientOptions, Auth0Error> {
        let blocking_client = match self.blocking_client {
            Some(value) => value,
            None => {
                let mut builder = reqwest::blocking::Client::builder();
                if let Some(timeout) = self.timeout {
                    builder = builder.timeout(timeout);
                }
                if let Some(proxy) = &self.proxy {
                    builder = builder.proxy(proxy.to_proxy()?);
                }
                builder.build()?
            }
        };
        let async_client = match self.async_client {
            Some(value) => value,
            None => {
                let mut builder = reqwest::Client::builder();
                if let Some(timeout) = self.timeout {
                    builder = builder.timeout(timeout);
                }
                if let Some(proxy) = &self.proxy {
                    builder = builder.proxy(proxy.to_proxy()?);
                }
                builder.build()?
            }
        };
        Ok(ClientOptions {
            blocking_client,
            async_client,
            timeout: self.timeout,
            retry: self.retry,
            logging: self.logging,
        })
    }
}

impl Default for ClientOptionsBuilder {
    fn default() -> Self {
        Self::new()
    }
}
