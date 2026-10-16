package cbor

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestDecodeRFC8949Vectors(t *testing.T) {
	cases := []struct {
		hex  string
		want any
	}{
		{"00", int64(0)},
		{"17", int64(23)},
		{"1818", int64(24)},
		{"1903e8", int64(1000)},
		{"1a000f4240", int64(1000000)},
		{"20", int64(-1)},
		{"3863", int64(-100)},
		{"f4", false},
		{"f5", true},
		{"f6", nil},
		{"6449455446", "IETF"},
	}
	for _, c := range cases {
		raw, _ := hex.DecodeString(c.hex)
		got, err := DecodeAll(raw)
		if err != nil {
			t.Fatalf("%s: %v", c.hex, err)
		}
		if got != c.want {
			t.Fatalf("%s: got %v want %v", c.hex, got, c.want)
		}
	}
}

func TestRoundTripCOSEShapedMap(t *testing.T) {
	in := Map{{int64(1), int64(2)}, {int64(3), int64(-7)}, {int64(-1), int64(1)}, {int64(-2), bytes.Repeat([]byte{7}, 32)}}
	got, err := DecodeAll(MustEncode(in))
	if err != nil {
		t.Fatal(err)
	}
	m := got.(map[any]any)
	if m[int64(3)] != int64(-7) || !bytes.Equal(m[int64(-2)].([]byte), bytes.Repeat([]byte{7}, 32)) {
		t.Fatalf("round trip lost COSE fields: %v", m)
	}
}

func TestDecodeReturnsRestSoAuthDataExtensionsCanFollowTheKey(t *testing.T) {
	key := MustEncode(Map{{int64(1), int64(1)}})
	ext := MustEncode(Map{{"credProtect", int64(2)}})
	_, rest, err := Decode(append(key, ext...))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(rest, ext) {
		t.Fatalf("rest must be exactly the extensions bytes")
	}
}

func TestRejectsHostileInput(t *testing.T) {
	cases := map[string]string{
		"length beyond input":      "5affffffff",
		"array count beyond input": "9affffffff",
		"truncated text":           "6449",
		"duplicate key":            "a2010201 03",
		"float":                    "f93c00",
		"trailing bytes":           "0000",
	}
	for name, h := range cases {
		raw, _ := hex.DecodeString(stripSpaces(h))
		if _, err := DecodeAll(raw); err == nil {
			t.Fatalf("%s: expected error", name)
		}
	}
	deep := bytes.Repeat([]byte{0x81}, 100)
	deep = append(deep, 0x00)
	if _, err := DecodeAll(deep); err == nil {
		t.Fatalf("deep nesting must be rejected to bound stack use")
	}
}

func stripSpaces(s string) string {
	return string(bytes.ReplaceAll([]byte(s), []byte(" "), nil))
}
