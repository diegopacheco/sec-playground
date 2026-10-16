package webauthn

import (
	"encoding/binary"
	"errors"
	"fmt"

	"passkey-server/internal/cbor"
	"passkey-server/internal/cose"
)

const (
	FlagUP byte = 0x01
	FlagUV byte = 0x04
	FlagBE byte = 0x08
	FlagBS byte = 0x10
	FlagAT byte = 0x40
	FlagED byte = 0x80

	minAuthDataLen  = 37
	maxCredentialID = 1023
)

type AttestedCredential struct {
	AAGUID       []byte
	CredentialID []byte
	PublicKeyRaw []byte
	PublicKey    *cose.Key
}

type AuthenticatorData struct {
	RPIDHash  []byte
	Flags     byte
	SignCount uint32
	Attested  *AttestedCredential
}

func (a *AuthenticatorData) Has(flag byte) bool {
	return a.Flags&flag == flag
}

func ParseAuthenticatorData(raw []byte) (*AuthenticatorData, error) {
	if len(raw) < minAuthDataLen {
		return nil, errors.New("authenticatorData: too short")
	}
	ad := &AuthenticatorData{
		RPIDHash:  raw[:32],
		Flags:     raw[32],
		SignCount: binary.BigEndian.Uint32(raw[33:37]),
	}
	rest := raw[37:]
	if ad.Has(FlagAT) {
		att, after, err := parseAttested(rest)
		if err != nil {
			return nil, err
		}
		ad.Attested = att
		rest = after
	}
	if ad.Has(FlagED) {
		_, after, err := cbor.Decode(rest)
		if err != nil {
			return nil, fmt.Errorf("authenticatorData: extensions: %w", err)
		}
		rest = after
	}
	if len(rest) != 0 {
		return nil, errors.New("authenticatorData: trailing bytes")
	}
	return ad, nil
}

func parseAttested(b []byte) (*AttestedCredential, []byte, error) {
	if len(b) < 18 {
		return nil, nil, errors.New("authenticatorData: attested credential data too short")
	}
	aaguid := b[:16]
	idLen := int(binary.BigEndian.Uint16(b[16:18]))
	if idLen == 0 || idLen > maxCredentialID {
		return nil, nil, errors.New("authenticatorData: invalid credential id length")
	}
	b = b[18:]
	if len(b) < idLen {
		return nil, nil, errors.New("authenticatorData: credential id truncated")
	}
	id := b[:idLen]
	b = b[idLen:]
	keyVal, rest, err := cbor.Decode(b)
	if err != nil {
		return nil, nil, fmt.Errorf("authenticatorData: public key: %w", err)
	}
	key, err := cose.FromMap(keyVal)
	if err != nil {
		return nil, nil, err
	}
	return &AttestedCredential{
		AAGUID:       aaguid,
		CredentialID: id,
		PublicKeyRaw: b[:len(b)-len(rest)],
		PublicKey:    key,
	}, rest, nil
}
