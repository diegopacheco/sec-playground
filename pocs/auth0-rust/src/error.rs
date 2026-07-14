use std::error::Error;
use std::fmt;

#[derive(Debug, Clone, PartialEq, Eq)]
pub struct ApiError {
    pub status: u16,
    pub code: Option<String>,
    pub message: Option<String>,
    pub body: String,
}

macro_rules! status_error {
    ($name:ident, $status:expr) => {
        #[derive(Debug, Clone, PartialEq, Eq)]
        pub struct $name {
            pub error: ApiError,
        }

        impl $name {
            pub const STATUS: u16 = $status;

            pub fn new(error: ApiError) -> Self {
                Self { error }
            }
        }
    };
}

status_error!(BadRequestError, 400);
status_error!(UnauthorizedError, 401);
status_error!(PaymentRequiredError, 402);
status_error!(ForbiddenError, 403);
status_error!(NotFoundError, 404);
status_error!(ConflictError, 409);
status_error!(GoneError, 410);
status_error!(PreconditionFailedError, 412);
status_error!(ContentTooLargeError, 413);
status_error!(TooManyRequestsError, 429);
status_error!(InternalServerError, 500);
status_error!(ServiceUnavailableError, 503);

#[derive(Debug, Clone, PartialEq, Eq)]
pub enum ApiErrorKind {
    BadRequest(BadRequestError),
    Unauthorized(UnauthorizedError),
    PaymentRequired(PaymentRequiredError),
    Forbidden(ForbiddenError),
    NotFound(NotFoundError),
    Conflict(ConflictError),
    Gone(GoneError),
    PreconditionFailed(PreconditionFailedError),
    ContentTooLarge(ContentTooLargeError),
    TooManyRequests(TooManyRequestsError),
    InternalServer(InternalServerError),
    ServiceUnavailable(ServiceUnavailableError),
    Other(ApiError),
}

impl ApiError {
    pub fn kind(&self) -> ApiErrorKind {
        match self.status {
            BadRequestError::STATUS => ApiErrorKind::BadRequest(BadRequestError::new(self.clone())),
            UnauthorizedError::STATUS => {
                ApiErrorKind::Unauthorized(UnauthorizedError::new(self.clone()))
            }
            PaymentRequiredError::STATUS => {
                ApiErrorKind::PaymentRequired(PaymentRequiredError::new(self.clone()))
            }
            ForbiddenError::STATUS => ApiErrorKind::Forbidden(ForbiddenError::new(self.clone())),
            NotFoundError::STATUS => ApiErrorKind::NotFound(NotFoundError::new(self.clone())),
            ConflictError::STATUS => ApiErrorKind::Conflict(ConflictError::new(self.clone())),
            GoneError::STATUS => ApiErrorKind::Gone(GoneError::new(self.clone())),
            PreconditionFailedError::STATUS => {
                ApiErrorKind::PreconditionFailed(PreconditionFailedError::new(self.clone()))
            }
            ContentTooLargeError::STATUS => {
                ApiErrorKind::ContentTooLarge(ContentTooLargeError::new(self.clone()))
            }
            TooManyRequestsError::STATUS => {
                ApiErrorKind::TooManyRequests(TooManyRequestsError::new(self.clone()))
            }
            InternalServerError::STATUS => {
                ApiErrorKind::InternalServer(InternalServerError::new(self.clone()))
            }
            ServiceUnavailableError::STATUS => {
                ApiErrorKind::ServiceUnavailable(ServiceUnavailableError::new(self.clone()))
            }
            _ => ApiErrorKind::Other(self.clone()),
        }
    }
}

#[derive(Debug)]
pub enum Auth0Error {
    InvalidInput(String),
    MissingEndpoint(String),
    MissingPathParameter(String),
    MissingAccessToken,
    MissingTokenExpiry,
    TokenCache,
    Http(ApiError),
    Transport(String),
    Json(String),
    Url(String),
}

impl fmt::Display for Auth0Error {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            Auth0Error::InvalidInput(value) => write!(f, "invalid input: {}", value),
            Auth0Error::MissingEndpoint(value) => write!(f, "missing endpoint: {}", value),
            Auth0Error::MissingPathParameter(value) => {
                write!(f, "missing path parameter: {}", value)
            }
            Auth0Error::MissingAccessToken => write!(f, "missing access token"),
            Auth0Error::MissingTokenExpiry => write!(f, "missing token expiry"),
            Auth0Error::TokenCache => write!(f, "token cache unavailable"),
            Auth0Error::Http(value) => write!(f, "http error {}: {}", value.status, value.body),
            Auth0Error::Transport(value) => write!(f, "transport error: {}", value),
            Auth0Error::Json(value) => write!(f, "json error: {}", value),
            Auth0Error::Url(value) => write!(f, "url error: {}", value),
        }
    }
}

impl Error for Auth0Error {}

impl From<reqwest::Error> for Auth0Error {
    fn from(value: reqwest::Error) -> Self {
        Auth0Error::Transport(value.to_string())
    }
}

impl From<serde_json::Error> for Auth0Error {
    fn from(value: serde_json::Error) -> Self {
        Auth0Error::Json(value.to_string())
    }
}

impl From<url::ParseError> for Auth0Error {
    fn from(value: url::ParseError) -> Self {
        Auth0Error::Url(value.to_string())
    }
}
