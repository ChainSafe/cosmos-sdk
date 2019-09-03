package simulation

import (
	"math/rand"

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Account contains a privkey, pubkey, address tuple
// eventually more useful data can be placed in here.
// (e.g. number of coins)
type Account struct {
	PrivKey crypto.PrivKey
	PubKey  crypto.PubKey
	Address sdk.AccAddress
}

// Equals returns true if two accounts are equal
func (acc Account) Equals(acc2 Account) bool {
	return acc.Address.Equals(acc2.Address)
}

// RandomAcc picks and returns a random account from an array and returs its
// position on the array
func RandomAcc(r *rand.Rand, accs []Account) (Account, int) {
	idx := r.Intn(len(accs))
	return accs[idx], idx
}

// RandomAccounts generates a given number of random accounts
func RandomAccounts(r *rand.Rand, nAccounts int) []Account {
	accs := make([]Account, nAccounts)
	for i := 0; i < nAccounts; i++ {
		// don't need that much entropy for simulation
		privkeySeed := make([]byte, 15)
		r.Read(privkeySeed)
		accs[i].PrivKey = secp256k1.GenPrivKeySecp256k1(privkeySeed)
		accs[i].PubKey = accs[i].PrivKey.PubKey()
		accs[i].Address = sdk.AccAddress(accs[i].PubKey.Address())
	}

	return accs
}

// FindAccount iterates over all the simulation accounts to find the one that matches
// the given address
func FindAccount(accs []Account, address sdk.AccAddress) (Account, bool) {
	for _, acc := range accs {
		if acc.Address.Equals(address) {
			return acc, true
		}
	}

	return Account{}, false
}
