package state

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Value is a capability for reading and writing on a specific key-value point
// in the state. Value consists of a mapping and a key. An actor holding a Value
// has a full access right on that state point.
type Value struct {
	mapping   Mapping
	key []byte
}

// NewValue constructs a Value.
func NewValue(m Mapping, key []byte) Value {
	return Value{
		m:   m,
		key: key,
	}
}

func (v Value) store(ctx Context) KVStore {
	return ctx.KVStore(v.m.storeKey)
}

// Cdc() returns the codec that the value is using to marshal/unmarshal
func (v Value) Cdc() *Codec {
	return v.m.Cdc()
}

func (v Value) Marshal(value interface{}) []byte {
	return v.m.cdc.MustMarshalBinaryBare(value)
}

func (v Value) Unmarshal(bz []byte, ptr interface{}) error {
	return v.m.cdc.UnmarshalBinaryBare(bz, ptr)
}

func (v Value) mustUnmarshal(bz []byte, ptr interface{}) {
	v.m.cdc.MustUnmarshalBinaryBare(bz, ptr)
}

// Get decodes and sets the stored value to the pointer if it exists. It will
// panic if the value exists but not unmarshalable.
func (v Value) Get(ctx Context, ptr interface{}) {
	bz := v.store(ctx).Get(v.KeyBytes())
	if bz != nil {
		v.mustUnmarshal(bz, ptr)
	}
}

// GetSafe decodes and sets the stored value to the pointer. It will return an
// error if the value does not exist or unmarshalable.
func (v Value) GetSafe(ctx Context, ptr interface{}) error {
	bz := v.store(ctx).Get(v.KeyBytes())
	if bz == nil {
		return ErrEmptyValue()
	}
	err := v.Unmarshal(bz, ptr)
	if err != nil {
		return ErrUnmarshal(err)
	}
	return nil
}

// GetRaw returns the raw bytes that is stored in the state.
func (v Value) GetRaw(ctx Context) []byte {
	return v.store(ctx).Get(v.KeyBytes())
}

// Set encodes and sets the argument to the state.
func (v Value) Set(ctx Context, o interface{}) {
	v.store(ctx).Set(v.KeyBytes(), v.Marshal(o))
}

// SetRaw sets the raw bytes to the state.
func (v Value) SetRaw(ctx Context, bz []byte) {
	v.store(ctx).Set(v.KeyBytes(), bz)
}

// Exists returns true if the stored value is not nil. It calls KVStore.Has()
// internally.
func (v Value) Exists(ctx Context) bool {
	return v.store(ctx).Has(v.KeyBytes())
}

// Delete removes the stored value. It calls KVStore.Delete() internally.
func (v Value) Delete(ctx Context) {
	v.store(ctx).Delete(v.KeyBytes())
}

func (v Value) StoreName() string {
	return v.m.StoreName()
}

func (v Value) PrefixBytes() []byte {
	return v.m.PrefixBytes()
}

// KeyBytes returns the prefixed key that the Value is providing to the KVStore.
func (v Value) KeyBytes() []byte {
	return v.m.KeyBytes(v.key)
}

func (v Value) QueryRaw(q ABCIQuerier) ([]byte, *Proof, error) {
	resp, err := q.Query(v.m.StoreName(), v.KeyBytes())
	if err != nil {
		return nil, nil, err
	}

	if !resp.IsOK() {
		return nil, nil, sdk.NewError(sdk.CodespaceRoot, sdk.CodeType(resp.Code), resp.Log)
	}

	return resp.Value, resp.Proof, nil
}

func (v Value) Query(q ABCIQuerier, ptr interface{}) (*Proof, error) {
	value, proof, err := v.QueryRaw(q)
	if err != nil {
		return nil, err
	}
	err = v.Cdc().UnmarshalBinaryBare(value, ptr)
	return proof, err
}
