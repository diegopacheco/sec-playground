# sec-playground

Hands-on security POCs: cryptography, identity, tokens, secrets, scanning and AI agent security. Every project lives under [pocs/](pocs/) and runs locally.

## 🤖 AI Agent Security

Agents call tools, hold credentials and read untrusted text.
These POCs put a guard between the agent and everything it can touch.

* [agent-credential-broker-mcp](pocs/agent-credential-broker-mcp/) - MCP broker that calls APIs for the agent, credentials stay in Vault or Keychain, single use intent tokens, hash chained audit
* [mcp-rug-pull-proxy](pocs/mcp-rug-pull-proxy/) - stdio MCP proxy that pins a SHA-256 per tool definition, blocks rug pulls and poisoned descriptions
* [agent-prompt-injection-firewall](pocs/agent-prompt-injection-firewall/) - Go reference monitor that treats retrieved instructions as data and enforces capabilities before any tool runs
* [agent-intent-bound-tokens](pocs/agent-intent-bound-tokens/) - Go service issuing signed one-time tokens bound to audience, action, resource, amount and deadline
* [agent-safehouse-poc](pocs/agent-safehouse-poc/) - Agent Safehouse: macOS-native sandboxing for LLM coding agents
* [agent-intent-eval](pocs/agent-intent-eval/) - Rust tool that scores how consistently an LLM answers 10 paraphrases of the same prompt, LLM as judge
* [agent-semantic-drift-detector](pocs/agent-semantic-drift-detector/) - Rust CLI that reruns a prompt, compares TF-IDF vectors to a baseline and plots model drift

## 🧰 Agent Skills

Skills that make a coding agent do a security job end to end.

* [agent-skill-threat-analyst](pocs/agent-skill-threat-analyst/) - Maps endpoints, inputs and auth boundaries, writes a threat analysis with attack vectors and scores
* [leak-detector-skill](pocs/leak-detector-skill/) - Scans a codebase for leaked PII, secrets and vulnerabilities before you push
* [VSDD-poc](pocs/VSDD-poc/) - Verified Spec-Driven Development: spec, test first, adversarial refinement, formal hardening

## 🔐 Identity & Authentication

Identity providers prove who the user is.
Frameworks turn that proof into sessions and access rules.

* [auth0-webapp](pocs/auth0-webapp/) - Spring Boot with Auth0 OpenID Connect login, protected profile and full logout
* [auth0-kotlin](pocs/auth0-kotlin/) - Idiomatic Kotlin wrapper over the official Auth0 Java SDK
* [auth0-scala3](pocs/auth0-scala3/) - Idiomatic Scala 3 wrapper over the official Auth0 Java SDK
* [auth0-rust](pocs/auth0-rust/) - Rust SDK for the Auth0 Authentication and Management APIs
* [okta-spring-boot-security-27x](pocs/okta-spring-boot-security-27x/) - Okta login on Spring Boot 2.7 with Spring Security
* [go-istio-keycloak](pocs/go-istio-keycloak/) - Go book service on kind where Istio enforces Keycloak JWTs
* [webauthn-passkeys-go](pocs/webauthn-passkeys-go/) - Passkey registration and login in Go with SQLite, no passwords at all
* [java-25-apache-shiro](pocs/java-25-apache-shiro/) - Apache Shiro authentication and authorization on Java 25
* [spring-security-boot-simple](pocs/spring-security-boot-simple/) - Minimal Spring Security setup on Spring Boot
* [spring-boot-2x-web-security](pocs/spring-boot-2x-web-security/) - Web security on Spring Boot 2.x

## 🎟️ Tokens & Capabilities

A token carries claims that a service must verify before it trusts them.
Capability tokens go further and limit exactly what the holder can do.

* [jwt-java-fun](pocs/jwt-java-fun/) - Create and validate JWTs in Java
* [jwt-sign-verify](pocs/jwt-sign-verify/) - Sign and verify a JWT with Bun and jsonwebtoken
* [jwt-open](pocs/jwt-open/) - Decode a JWT without verifying it, showing that claims are readable by anyone
* [macaroon-auth-go](pocs/macaroon-auth-go/) - Macaroon tokens in Go, attenuated by resource, operation, location and expiry
* [token-blast-radius](pocs/token-blast-radius/) - Go + React app that calculates how far a credential really reaches

## 🔒 Symmetric Encryption

One secret key encrypts and decrypts.
Mode, IV and padding choices decide if it is actually safe.

* [aes-cipher-fun](pocs/aes-cipher-fun/) - AES with Go crypto/aes and crypto/cipher
* [aes-encryption-fun-python](pocs/aes-encryption-fun-python/) - AES with PyCryptodome in Python
* [aes-encryption-fun-rust](pocs/aes-encryption-fun-rust/) - AES-128 CBC with PKCS7 padding in Rust
* [security-aes256-fun](pocs/security-aes256-fun/) - AES-256 in Java
* [security-file-enc-dec](pocs/security-file-enc-dec/) - Encrypt and decrypt files in Java
* [commons-crypto-fun](pocs/commons-crypto-fun/) - Apache Commons Crypto ciphers
* [security-google-tink-fun](pocs/security-google-tink-fun/) - Google Tink AEAD with keysets
* [dumb-encryption-poc](pocs/dumb-encryption-poc/) - Naive encryption in Java, what not to do
* [encryption-deep-dive-poc](pocs/encryption-deep-dive-poc/) - Keys, PRNG, encryptor and decryptor built step by step
* [compression-encryption-encoding-fun](pocs/compression-encryption-encoding-fun/) - Compression, encryption and encoding chained together
* [native-creds-sc](pocs/native-creds-sc/) - Spring Boot credentials encryption service compiled to a GraalVM native image

## 🔑 Keys, Signatures & Certificates

Asymmetric keys sign, verify and exchange secrets.
Keystores, PKCS8 and PGP are how those keys get stored and shared.

* [security-generate-key-java](pocs/security-generate-key-java/) - Generate keys with the Java crypto APIs
* [security-n-keys-java](pocs/security-n-keys-java/) - Envelope encryption with a root key, a KEK and a DEK
* [security-signature-fun](pocs/security-signature-fun/) - SHA256withDSA digital signatures
* [security-java-keystore-jks-fun](pocs/security-java-keystore-jks-fun/) - Store and load keys in a Java KeyStore
* [bouncy-castle-generate-certificate-jks-fun](pocs/bouncy-castle-generate-certificate-jks-fun/) - Generate a certificate with Bouncy Castle into a JKS
* [pkcs8-simple-pocs](pocs/pkcs8-simple-pocs/) - Write PKCS8 private keys
* [pkcs8-read-simple-poc](pocs/pkcs8-read-simple-poc/) - Read PKCS8 private keys
* [sec-bouncycastle-openpgp-fun](pocs/sec-bouncycastle-openpgp-fun/) - OpenPGP with Bouncy Castle
* [pgpainless-fun](pocs/pgpainless-fun/) - OpenPGP with PGPainless and SOP
* [certificate-requests](pocs/certificate-requests/) - Create a certificate and call HTTPS with it from Python
* [secure-random-fun](pocs/secure-random-fun/) - SHA1PRNG vs NativePRNG with SecureRandom

## #️⃣ Hashing & Passwords

Hashes turn data into a fixed fingerprint that cannot be reversed.
Passwords need slow hashes and strong policies, not fast ones.

* [bcrypt-fun](pocs/bcrypt-fun/) - bcrypt password hashing in Go
* [sha-hash-fun](pocs/sha-hash-fun/) - SHA-256 in Go
* [security-hash-sha256](pocs/security-hash-sha256/) - SHA-256 and SHA-512 in Java
* [security-hmac256](pocs/security-hmac256/) - HMAC-SHA256 message authentication
* [security-md5-fun](pocs/security-md5-fun/) - MD5 and why it is broken
* [passay-fun](pocs/passay-fun/) - Password policy rules with Passay

## ☁️ Secrets & Key Management

Secrets belong in a vault or a KMS, never in the code.
The app asks for them at runtime and the master key never leaves the service.

* [spring-vault](pocs/spring-vault/) - Spring Boot reading secrets from HashiCorp Vault
* [sec-asm-fun](pocs/sec-asm-fun/) - AWS Secrets Manager create, describe and read from Java
* [security-kms-fun](pocs/security-kms-fun/) - AWS KMS envelope encryption with data keys, key created by Terraform

## 🌐 Transport Security

TLS encrypts the wire.
mTLS makes both sides prove who they are.

* [spring-boot-2-tls](pocs/spring-boot-2-tls/) - HTTPS on Spring Boot 2
* [spring-boot-2-mTLS-fun](pocs/spring-boot-2-mTLS-fun/) - Mutual TLS between a gateway and a microservice

## 🔍 Scanning & Threat Modeling

Scanners find bugs, misconfigurations and vulnerable packages before attackers do.
Threat models find the risks before there is code to scan.

* [sec-find-sec-bugs-maven-simple](pocs/sec-find-sec-bugs-maven-simple/) - Find Security Bugs SpotBugs plugin on a Maven project
* [sonarq-cli-fun](pocs/sonarq-cli-fun/) - SonarQube CLI
* [tfsec-fun](pocs/tfsec-fun/) - tfsec scanning Terraform
* [trivy-vulnerabilities-checker](pocs/trivy-vulnerabilities-checker/) - Trivy vulnerability scanning
* [docker-registry-trivy-k8s--registry-fun](pocs/docker-registry-trivy-k8s--registry-fun/) - Local registry with podman, images scanned by Trivy
* [falco-runtime-security](pocs/falco-runtime-security/) - Falco runtime detection on Kubernetes
* [sec-threagile](pocs/sec-threagile/) - Threagile threat model as YAML
* [detect-axios-rat.sh](pocs/detect-axios-rat.sh) - Scans lockfiles for the compromised axios 1.14.1 and 0.30.4 releases

## 🕵️ Fraud, PII & Anomaly Detection

Detection finds the bad record inside a lot of good ones.
Rules catch the known patterns, models catch the unusual ones.

* [pii-simple-detector-fun](pocs/pii-simple-detector-fun/) - Regex PII detection in Java
* [credit-card-validation](pocs/credit-card-validation/) - Credit card validation with the Luhn algorithm
* [sklearn-anomaly-detections-z-score](pocs/sklearn-anomaly-detections-z-score/) - Anomaly detection with Z-score
* [sklearn-anomaly-detections-IsolationForest](pocs/sklearn-anomaly-detections-IsolationForest/) - Anomaly detection with Isolation Forest
* [sklearn-anomaly-detections-LOF](pocs/sklearn-anomaly-detections-LOF/) - Anomaly detection with Local Outlier Factor
* [transformers-dumb-event-to-fraud](pocs/transformers-dumb-event-to-fraud/) - Fine-tune a Hugging Face transformer to classify events as fraud
