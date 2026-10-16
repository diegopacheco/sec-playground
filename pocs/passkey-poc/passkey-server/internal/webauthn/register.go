package webauthn

import (
	"crypto/sha256"
	"crypto/subtle"
	"crypto/x509"
	"errors"
	"fmt"

	"passkey-server/internal/cbor"
	"passkey-server/internal/cose"
)

type Credential struct {
	ID             []byte   `json:"id"`
	PublicKey      []byte   `json:"publicKey"`
	Alg            int64    `json:"alg"`
	SignCount      uint32   `json:"signCount"`
	AAGUID         []byte   `json:"aaguid"`
	AttestationFmt string   `json:"attestationFmt"`
	BackupEligible bool     `json:"backupEligible"`
	BackedUp       bool     `json:"backedUp"`
	Transports     []string `json:"transports"`
}

type RegistrationInput struct {
	ClientDataJSON    []byte
	AttestationObject []byte
	Transports        []string
}

func VerifyRegistration(rp RelyingParty, challenge []byte, in RegistrationInput) (*Credential, error) {
	if _, err := VerifyClientData(in.ClientDataJSON, TypeCreate, challenge, rp.Origins); err != nil {
		return nil, err
	}
	obj, err := cbor.DecodeAll(in.AttestationObject)
	if err != nil {
		return nil, fmt.Errorf("attestationObject: %w", err)
	}
	m, ok := obj.(map[any]any)
	if !ok {
		return nil, errors.New("attestationObject: not a map")
	}
	format, _ := m["fmt"].(string)
	stmt, _ := m["attStmt"].(map[any]any)
	rawAuth, _ := m["authData"].([]byte)
	if rawAuth == nil || stmt == nil {
		return nil, errors.New("attestationObject: missing authData or attStmt")
	}
	ad, err := ParseAuthenticatorData(rawAuth)
	if err != nil {
		return nil, err
	}
	if err := checkAuthData(ad, rp); err != nil {
		return nil, err
	}
	if ad.Attested == nil {
		return nil, errors.New("authenticatorData: no attested credential data")
	}
	if ad.Has(FlagBS) && !ad.Has(FlagBE) {
		return nil, errors.New("authenticatorData: backed up without backup eligible")
	}
	clientHash := sha256.Sum256(in.ClientDataJSON)
	if err := verifyAttestation(format, stmt, rawAuth, clientHash[:], ad.Attested.PublicKey); err != nil {
		return nil, err
	}
	return &Credential{
		ID:             ad.Attested.CredentialID,
		PublicKey:      ad.Attested.PublicKeyRaw,
		Alg:            ad.Attested.PublicKey.Alg,
		SignCount:      ad.SignCount,
		AAGUID:         ad.Attested.AAGUID,
		AttestationFmt: format,
		BackupEligible: ad.Has(FlagBE),
		BackedUp:       ad.Has(FlagBS),
		Transports:     in.Transports,
	}, nil
}

func checkAuthData(ad *AuthenticatorData, rp RelyingParty) error {
	want := sha256.Sum256([]byte(rp.ID))
	if subtle.ConstantTimeCompare(ad.RPIDHash, want[:]) != 1 {
		return errors.New("authenticatorData: rpIdHash mismatch")
	}
	if !ad.Has(FlagUP) {
		return errors.New("authenticatorData: user not present")
	}
	if !ad.Has(FlagUV) {
		return errors.New("authenticatorData: user not verified")
	}
	return nil
}

func verifyAttestation(format string, stmt map[any]any, rawAuth, clientHash []byte, credKey *cose.Key) error {
	switch format {
	case "none":
		if len(stmt) != 0 {
			return errors.New("attestation none: attStmt must be empty")
		}
		return nil
	case "packed":
		return verifyPacked(stmt, append(append([]byte{}, rawAuth...), clientHash...), credKey)
	}
	return fmt.Errorf("attestation format %q not supported", format)
}

func verifyPacked(stmt map[any]any, signed []byte, credKey *cose.Key) error {
	alg, ok := stmt["alg"].(int64)
	if !ok {
		return errors.New("attestation packed: missing alg")
	}
	sig, ok := stmt["sig"].([]byte)
	if !ok {
		return errors.New("attestation packed: missing sig")
	}
	x5c, hasCert := stmt["x5c"].([]any)
	if !hasCert {
		if alg != credKey.Alg {
			return errors.New("attestation packed: self attestation alg mismatch")
		}
		return credKey.Verify(signed, sig)
	}
	if len(x5c) == 0 {
		return errors.New("attestation packed: empty x5c")
	}
	der, ok := x5c[0].([]byte)
	if !ok {
		return errors.New("attestation packed: x5c entry is not bytes")
	}
	cert, err := x509.ParseCertificate(der)
	if err != nil {
		return fmt.Errorf("attestation packed: %w", err)
	}
	attKey := &cose.Key{Alg: alg, Public: cert.PublicKey}
	if alg != cose.AlgES256 && alg != cose.AlgRS256 && alg != cose.AlgEdDSA {
		return fmt.Errorf("attestation packed: unsupported alg %d", alg)
	}
	return attKey.Verify(signed, sig)
}
