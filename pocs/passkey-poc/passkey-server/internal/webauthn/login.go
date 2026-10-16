package webauthn

import (
	"crypto/sha256"
	"errors"

	"passkey-server/internal/cose"
)

type AssertionInput struct {
	ClientDataJSON    []byte
	AuthenticatorData []byte
	Signature         []byte
}

type AssertionResult struct {
	SignCount uint32
	BackedUp  bool
}

var ErrCloned = errors.New("signature counter did not increase, possible cloned authenticator")

func VerifyAssertion(rp RelyingParty, challenge []byte, cred *Credential, in AssertionInput) (*AssertionResult, error) {
	if _, err := VerifyClientData(in.ClientDataJSON, TypeGet, challenge, rp.Origins); err != nil {
		return nil, err
	}
	ad, err := ParseAuthenticatorData(in.AuthenticatorData)
	if err != nil {
		return nil, err
	}
	if err := checkAuthData(ad, rp); err != nil {
		return nil, err
	}
	if ad.Has(FlagBE) != cred.BackupEligible {
		return nil, errors.New("authenticatorData: backup eligibility changed")
	}
	key, err := cose.Parse(cred.PublicKey)
	if err != nil {
		return nil, err
	}
	clientHash := sha256.Sum256(in.ClientDataJSON)
	signed := append(append([]byte{}, in.AuthenticatorData...), clientHash[:]...)
	if err := key.Verify(signed, in.Signature); err != nil {
		return nil, err
	}
	if (ad.SignCount != 0 || cred.SignCount != 0) && ad.SignCount <= cred.SignCount {
		return nil, ErrCloned
	}
	return &AssertionResult{SignCount: ad.SignCount, BackedUp: ad.Has(FlagBS)}, nil
}
