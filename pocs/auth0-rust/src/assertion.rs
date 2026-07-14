use crate::error::Auth0Error;
use jsonwebtoken::{Algorithm, EncodingKey, Header};
use serde::Serialize;
use std::time::{SystemTime, UNIX_EPOCH};

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum RsaSigningAlgorithm {
    Rsa256,
    Rsa384,
}

#[derive(Clone)]
pub struct RsaClientAssertionSigner {
    key: EncodingKey,
    algorithm: RsaSigningAlgorithm,
}

impl std::fmt::Debug for RsaClientAssertionSigner {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        f.debug_struct("RsaClientAssertionSigner")
            .field("algorithm", &self.algorithm)
            .finish()
    }
}

impl RsaClientAssertionSigner {
    pub fn from_pem(pem: impl AsRef<[u8]>) -> Result<Self, Auth0Error> {
        Self::from_pem_with_algorithm(pem, RsaSigningAlgorithm::Rsa256)
    }

    pub fn from_pem_with_algorithm(
        pem: impl AsRef<[u8]>,
        algorithm: RsaSigningAlgorithm,
    ) -> Result<Self, Auth0Error> {
        Ok(Self {
            key: EncodingKey::from_rsa_pem(pem.as_ref())
                .map_err(|value| Auth0Error::InvalidInput(value.to_string()))?,
            algorithm,
        })
    }

    pub fn sign(
        &self,
        issuer: impl Into<String>,
        audience: impl Into<String>,
        subject: impl Into<String>,
    ) -> Result<String, Auth0Error> {
        let now = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .map_err(|value| Auth0Error::InvalidInput(value.to_string()))?
            .as_secs();
        let claims = ClientAssertionClaims {
            iss: issuer.into(),
            aud: audience.into(),
            sub: subject.into(),
            iat: now,
            exp: now + 180,
            jti: format!(
                "{}-{}",
                now,
                SystemTime::now().elapsed().unwrap_or_default().as_nanos()
            ),
        };
        let algorithm = match self.algorithm {
            RsaSigningAlgorithm::Rsa256 => Algorithm::RS256,
            RsaSigningAlgorithm::Rsa384 => Algorithm::RS384,
        };
        jsonwebtoken::encode(&Header::new(algorithm), &claims, &self.key)
            .map_err(|value| Auth0Error::InvalidInput(value.to_string()))
    }
}

#[derive(Debug, Serialize)]
struct ClientAssertionClaims {
    iss: String,
    aud: String,
    sub: String,
    iat: u64,
    exp: u64,
    jti: String,
}
