package cose

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rsa"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"

	"passkey-server/internal/cbor"
)

const (
	AlgES256 int64 = -7
	AlgEdDSA int64 = -8
	AlgRS256 int64 = -257

	ktyOKP int64 = 1
	ktyEC2 int64 = 2
	ktyRSA int64 = 3

	crvP256    int64 = 1
	crvEd25519 int64 = 6

	minRSABits = 2048
)

type Key struct {
	Alg    int64
	Public crypto.PublicKey
}

func Parse(raw []byte) (*Key, error) {
	v, err := cbor.DecodeAll(raw)
	if err != nil {
		return nil, err
	}
	return FromMap(v)
}

func FromMap(v any) (*Key, error) {
	m, ok := v.(map[any]any)
	if !ok {
		return nil, errors.New("cose: key is not a map")
	}
	kty, err := intField(m, 1)
	if err != nil {
		return nil, err
	}
	alg, err := intField(m, 3)
	if err != nil {
		return nil, err
	}
	switch {
	case kty == ktyEC2 && alg == AlgES256:
		return parseEC2(m)
	case kty == ktyOKP && alg == AlgEdDSA:
		return parseOKP(m)
	case kty == ktyRSA && alg == AlgRS256:
		return parseRSA(m)
	}
	return nil, fmt.Errorf("cose: unsupported kty %d alg %d", kty, alg)
}

func parseEC2(m map[any]any) (*Key, error) {
	crv, err := intField(m, -1)
	if err != nil {
		return nil, err
	}
	if crv != crvP256 {
		return nil, fmt.Errorf("cose: unsupported EC2 curve %d", crv)
	}
	x, err := bytesField(m, -2)
	if err != nil {
		return nil, err
	}
	y, err := bytesField(m, -3)
	if err != nil {
		return nil, err
	}
	if len(x) != 32 || len(y) != 32 {
		return nil, errors.New("cose: P-256 coordinates must be 32 bytes")
	}
	point := append(append([]byte{0x04}, x...), y...)
	pub, err := ecdsa.ParseUncompressedPublicKey(elliptic.P256(), point)
	if err != nil {
		return nil, fmt.Errorf("cose: invalid P-256 point: %w", err)
	}
	return &Key{Alg: AlgES256, Public: pub}, nil
}

func parseOKP(m map[any]any) (*Key, error) {
	crv, err := intField(m, -1)
	if err != nil {
		return nil, err
	}
	if crv != crvEd25519 {
		return nil, fmt.Errorf("cose: unsupported OKP curve %d", crv)
	}
	x, err := bytesField(m, -2)
	if err != nil {
		return nil, err
	}
	if len(x) != ed25519.PublicKeySize {
		return nil, errors.New("cose: Ed25519 key must be 32 bytes")
	}
	return &Key{Alg: AlgEdDSA, Public: ed25519.PublicKey(x)}, nil
}

func parseRSA(m map[any]any) (*Key, error) {
	n, err := bytesField(m, -1)
	if err != nil {
		return nil, err
	}
	e, err := bytesField(m, -2)
	if err != nil {
		return nil, err
	}
	if len(e) == 0 || len(e) > 4 {
		return nil, errors.New("cose: invalid RSA exponent")
	}
	exp := 0
	for _, b := range e {
		exp = exp<<8 | int(b)
	}
	pub := &rsa.PublicKey{N: new(big.Int).SetBytes(n), E: exp}
	if pub.N.BitLen() < minRSABits {
		return nil, fmt.Errorf("cose: RSA key must be at least %d bits", minRSABits)
	}
	if exp < 3 || exp%2 == 0 {
		return nil, errors.New("cose: invalid RSA exponent")
	}
	return &Key{Alg: AlgRS256, Public: pub}, nil
}

func (k *Key) Verify(data, sig []byte) error {
	switch pub := k.Public.(type) {
	case *ecdsa.PublicKey:
		digest := sha256.Sum256(data)
		if !ecdsa.VerifyASN1(pub, digest[:], sig) {
			return errors.New("cose: invalid ES256 signature")
		}
		return nil
	case ed25519.PublicKey:
		if !ed25519.Verify(pub, data, sig) {
			return errors.New("cose: invalid EdDSA signature")
		}
		return nil
	case *rsa.PublicKey:
		digest := sha256.Sum256(data)
		if err := rsa.VerifyPKCS1v15(pub, crypto.SHA256, digest[:], sig); err != nil {
			return errors.New("cose: invalid RS256 signature")
		}
		return nil
	}
	return errors.New("cose: unsupported key type")
}

func intField(m map[any]any, label int64) (int64, error) {
	v, ok := m[label].(int64)
	if !ok {
		return 0, fmt.Errorf("cose: missing integer label %d", label)
	}
	return v, nil
}

func bytesField(m map[any]any, label int64) ([]byte, error) {
	v, ok := m[label].([]byte)
	if !ok {
		return nil, fmt.Errorf("cose: missing bytes label %d", label)
	}
	return v, nil
}
