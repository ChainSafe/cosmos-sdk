package state

import (
	"encoding/binary"
	"fmt"
	"strconv"
)

// IntEncoding is an enum type defining the integer serialization scheme.
// All encoding schemes preserves order.
type IntEncoding byte

const (
	// Dec is human readable decimal encoding scheme.
	// Has fixed length of 20 bytes.
	Dec IntEncoding = iota
	// Hex is human readable hexadecimal encoding scheme
	// Has fixed length of 16 bytes.
	Hex
	// Bin is machine readable big endian encoding scheme
	// Has fixed length of 8 bytes
	Bin
)

// Indexer is a integer typed key wrapper for Mapping. Except for the type
// checking, it does not alter the behaviour. All keys are encoded depending on
// the IntEncoding.
type Indexer struct {
	mapping Mapping

	encoding IntEncoding
}

// Indexer() constructs the Indexer with an IntEncoding
func (m Mapping) Indexer(enc IntEncoding) Indexer {
	return Indexer{
		m:   m,
		enc: enc,
	}
}

// EncodeInt provides order preserving integer encoding function.
func EncodeInt(index uint64, enc IntEncoding) (res []byte) {
	switch enc {
	case Dec:
		// return decimal number index, 20-length 0 padded
		return []byte(fmt.Sprintf("%020d", index))

	case Hex:
		// return hexadecimal number index, 20-length 0 padded
		return []byte(fmt.Sprintf("%016x", index))

	case Bin:
		// return bigendian encoded number index, 8-length
		res = make([]byte, 8)
		binary.BigEndian.PutUint64(res, index)
		return

	default:
		panic("invalid IntEncoding")
	}
}

// DecodeInt provides integer decoding function, inversion of EncodeInt.
func DecodeInt(bz []byte, enc IntEncoding) (res uint64, err error) {
	switch enc {
	case Dec:
		return strconv.ParseUint(string(bz), 10, 64)

	case Hex:
		return strconv.ParseUint(string(bz), 16, 64)

	case Bin:
		return binary.BigEndian.Uint64(bz), nil

	default:
		panic("invalid IntEncoding")
	}
}

// Value returns the Value corresponding to the provided index.
func (ix Indexer) Value(index uint64) Value {
	return ix.m.Value(EncodeInt(index, ix.enc))
}

// Get decodes and sets the stored value to the pointer if it exists. It will
// panic if the value exists but not unmarshalable.
func (ix Indexer) Get(ctx Context, index uint64, ptr interface{}) {
	ix.Value(index).Get(ctx, ptr)
}

// GetSafe decodes and sets the stored value to the pointer. It will return an
// error if the value does not exist or unmarshalable.
func (ix Indexer) GetSafe(ctx Context, index uint64, ptr interface{}) error {
	return ix.Value(index).GetSafe(ctx, ptr)
}

// Set encodes and sets the argument to the state.
func (ix Indexer) Set(ctx Context, index uint64, o interface{}) {
	ix.Value(index).Set(ctx, o)
}

// SetRaw sets the raw bytestring argument to the state
func (ix Indexer) SetRaw(ctx Context, index uint64, value []byte) {
	ix.Value(index).SetRaw(ctx, value)
}

// Has returns true if the stored value is not nil.
func (ix Indexer) Has(ctx Context, index uint64) bool {
	return ix.Value(index).Exists(ctx)
}

// Delete removes the stored value.
func (ix Indexer) Delete(ctx Context, index uint64) {
	ix.Value(index).Delete(ctx)
}

// Prefix returns a new Indexer with the updated prefix
func (ix Indexer) Prefix(prefix []byte) Indexer {
	return Indexer{
		m: ix.m.Prefix(prefix),

		enc: ix.enc,
	}
}
