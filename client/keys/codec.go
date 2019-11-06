package keys

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// cdc defines codec to be used with key operations
var cdc *codec.Codec

func init() {
	cdc = codec.New()
	codec.RegisterCrypto(cdc)
	// cdc.Seal()
}

// RegisterKeyTypeCodec registers an external account type defined in
// another module for the internal ModuleCdc.
func RegisterKeyTypeCodec(o interface{}, name string) {
	cdc.RegisterConcrete(o, name, nil)
}

// marshal keys
func MarshalJSON(o interface{}) ([]byte, error) {
	return cdc.MarshalJSON(o)
}

// unmarshal json
func UnmarshalJSON(bz []byte, ptr interface{}) error {
	return cdc.UnmarshalJSON(bz, ptr)
}
