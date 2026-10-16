package cbor

import (
	"encoding/binary"
	"fmt"
)

type Map [][2]any

func Encode(v any) ([]byte, error) {
	var out []byte
	return appendValue(out, v)
}

func MustEncode(v any) []byte {
	b, err := Encode(v)
	if err != nil {
		panic(err)
	}
	return b
}

func appendHead(out []byte, major byte, n uint64) []byte {
	m := major << 5
	switch {
	case n < 24:
		return append(out, m|byte(n))
	case n <= 0xff:
		return append(out, m|24, byte(n))
	case n <= 0xffff:
		return binary.BigEndian.AppendUint16(append(out, m|25), uint16(n))
	case n <= 0xffffffff:
		return binary.BigEndian.AppendUint32(append(out, m|26), uint32(n))
	}
	return binary.BigEndian.AppendUint64(append(out, m|27), n)
}

func appendInt(out []byte, n int64) []byte {
	if n >= 0 {
		return appendHead(out, 0, uint64(n))
	}
	return appendHead(out, 1, uint64(-1-n))
}

func appendValue(out []byte, v any) ([]byte, error) {
	var err error
	switch x := v.(type) {
	case nil:
		return append(out, 0xf6), nil
	case bool:
		if x {
			return append(out, 0xf5), nil
		}
		return append(out, 0xf4), nil
	case int:
		return appendInt(out, int64(x)), nil
	case int64:
		return appendInt(out, x), nil
	case []byte:
		return append(appendHead(out, 2, uint64(len(x))), x...), nil
	case string:
		return append(appendHead(out, 3, uint64(len(x))), x...), nil
	case []any:
		out = appendHead(out, 4, uint64(len(x)))
		for _, e := range x {
			if out, err = appendValue(out, e); err != nil {
				return nil, err
			}
		}
		return out, nil
	case Map:
		out = appendHead(out, 5, uint64(len(x)))
		for _, p := range x {
			if out, err = appendValue(out, p[0]); err != nil {
				return nil, err
			}
			if out, err = appendValue(out, p[1]); err != nil {
				return nil, err
			}
		}
		return out, nil
	}
	return nil, fmt.Errorf("cbor: cannot encode %T", v)
}
