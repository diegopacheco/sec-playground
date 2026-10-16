package cbor

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
)

const maxDepth = 16

var ErrTruncated = errors.New("cbor: truncated input")

func Decode(data []byte) (any, []byte, error) {
	d := &decoder{data: data}
	v, err := d.value(0)
	if err != nil {
		return nil, nil, err
	}
	return v, d.data[d.pos:], nil
}

func DecodeAll(data []byte) (any, error) {
	v, rest, err := Decode(data)
	if err != nil {
		return nil, err
	}
	if len(rest) != 0 {
		return nil, fmt.Errorf("cbor: %d trailing bytes", len(rest))
	}
	return v, nil
}

type decoder struct {
	data []byte
	pos  int
}

func (d *decoder) remaining() int {
	return len(d.data) - d.pos
}

func (d *decoder) take(n uint64) ([]byte, error) {
	if n > uint64(d.remaining()) {
		return nil, ErrTruncated
	}
	b := d.data[d.pos : d.pos+int(n)]
	d.pos += int(n)
	return b, nil
}

func (d *decoder) head() (byte, uint64, error) {
	b, err := d.take(1)
	if err != nil {
		return 0, 0, err
	}
	major := b[0] >> 5
	info := b[0] & 0x1f
	if major == 7 && info >= 24 {
		return 0, 0, fmt.Errorf("cbor: floats and extended simple values are not supported")
	}
	switch {
	case info < 24:
		return major, uint64(info), nil
	case info == 24:
		x, err := d.take(1)
		if err != nil {
			return 0, 0, err
		}
		return major, uint64(x[0]), nil
	case info == 25:
		x, err := d.take(2)
		if err != nil {
			return 0, 0, err
		}
		return major, uint64(binary.BigEndian.Uint16(x)), nil
	case info == 26:
		x, err := d.take(4)
		if err != nil {
			return 0, 0, err
		}
		return major, uint64(binary.BigEndian.Uint32(x)), nil
	case info == 27:
		x, err := d.take(8)
		if err != nil {
			return 0, 0, err
		}
		return major, binary.BigEndian.Uint64(x), nil
	}
	return 0, 0, fmt.Errorf("cbor: unsupported additional info %d", info)
}

func (d *decoder) value(depth int) (any, error) {
	if depth > maxDepth {
		return nil, errors.New("cbor: nesting too deep")
	}
	major, arg, err := d.head()
	if err != nil {
		return nil, err
	}
	switch major {
	case 0:
		if arg > math.MaxInt64 {
			return nil, errors.New("cbor: integer overflow")
		}
		return int64(arg), nil
	case 1:
		if arg > math.MaxInt64 {
			return nil, errors.New("cbor: integer overflow")
		}
		return -1 - int64(arg), nil
	case 2:
		b, err := d.take(arg)
		if err != nil {
			return nil, err
		}
		out := make([]byte, len(b))
		copy(out, b)
		return out, nil
	case 3:
		b, err := d.take(arg)
		if err != nil {
			return nil, err
		}
		return string(b), nil
	case 4:
		return d.array(arg, depth)
	case 5:
		return d.mapValue(arg, depth)
	case 6:
		return d.value(depth + 1)
	case 7:
		return d.simple(arg)
	}
	return nil, fmt.Errorf("cbor: unknown major type %d", major)
}

func (d *decoder) array(n uint64, depth int) ([]any, error) {
	if n > uint64(d.remaining()) {
		return nil, ErrTruncated
	}
	out := make([]any, 0, n)
	for i := uint64(0); i < n; i++ {
		v, err := d.value(depth + 1)
		if err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, nil
}

func (d *decoder) mapValue(n uint64, depth int) (map[any]any, error) {
	if n > uint64(d.remaining()) {
		return nil, ErrTruncated
	}
	out := make(map[any]any, n)
	for i := uint64(0); i < n; i++ {
		k, err := d.value(depth + 1)
		if err != nil {
			return nil, err
		}
		switch k.(type) {
		case int64, string:
		default:
			return nil, fmt.Errorf("cbor: unsupported map key type %T", k)
		}
		if _, dup := out[k]; dup {
			return nil, fmt.Errorf("cbor: duplicate map key %v", k)
		}
		v, err := d.value(depth + 1)
		if err != nil {
			return nil, err
		}
		out[k] = v
	}
	return out, nil
}

func (d *decoder) simple(arg uint64) (any, error) {
	switch arg {
	case 20:
		return false, nil
	case 21:
		return true, nil
	case 22, 23:
		return nil, nil
	}
	return nil, fmt.Errorf("cbor: unsupported simple value %d", arg)
}
