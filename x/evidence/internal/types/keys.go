package types

import (
	cmn "github.com/tendermint/tendermint/libs/common"
)

const (
	// ModuleName defines the module name
	ModuleName = "evidence"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// DefaultParamspace defines the module's default paramspace name
	DefaultParamspace = ModuleName
)

// KVStore key prefixes
var (
	KeyPrefixEvidence = []byte{0x00}
)

// EvidenceKey returns the KVStore key for persisting Evidence.
func EvidenceKey(hash cmn.HexBytes) []byte {
	return append(KeyPrefixEvidence, hash...)
}
