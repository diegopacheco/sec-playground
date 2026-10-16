package cose

import (
	"bytes"
	"testing"

	"passkey-server/internal/cbor"
)

func TestPointNotOnCurveIsRejected(t *testing.T) {
	raw := cbor.MustEncode(cbor.Map{{int64(1), int64(2)}, {int64(3), AlgES256}, {int64(-1), int64(1)}, {int64(-2), bytes.Repeat([]byte{1}, 32)}, {int64(-3), bytes.Repeat([]byte{2}, 32)}})
	if _, err := Parse(raw); err == nil {
		t.Fatalf("an invalid curve point enables invalid-curve attacks and must be rejected")
	}
}

func TestWeakRSAKeyIsRejected(t *testing.T) {
	raw := cbor.MustEncode(cbor.Map{{int64(1), int64(3)}, {int64(3), AlgRS256}, {int64(-1), bytes.Repeat([]byte{0xff}, 128)}, {int64(-2), []byte{1, 0, 1}}})
	if _, err := Parse(raw); err == nil {
		t.Fatalf("1024 bit RSA must be rejected")
	}
}

func TestAlgorithmAndKeyTypeMustMatch(t *testing.T) {
	raw := cbor.MustEncode(cbor.Map{{int64(1), int64(1)}, {int64(3), AlgES256}, {int64(-1), int64(6)}, {int64(-2), bytes.Repeat([]byte{1}, 32)}})
	if _, err := Parse(raw); err == nil {
		t.Fatalf("an OKP key claiming ES256 is algorithm confusion and must be rejected")
	}
}
