package solo

import (
	"math"

	"github.com/stretchr/testify/require"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/nnikolash/wasp-types-exported/packages/cryptolib"
	"github.com/nnikolash/wasp-types-exported/packages/testutil/testkey"
	"github.com/nnikolash/wasp-types-exported/packages/testutil/utxodb"
)

func (env *Solo) NewKeyPairFromIndex(index int) *cryptolib.KeyPair {
	seed := env.NewSeedFromIndex(index)
	return cryptolib.KeyPairFromSeed(*seed)
}

func (env *Solo) NewSeedFromIndex(index int) *cryptolib.Seed {
	if index < 0 {
		// SubSeed takes a "uint31"
		index += math.MaxUint32 / 2
	}
	seed := cryptolib.SubSeed(env.seed[:], uint32(index))
	return &seed
}

// NewKeyPairWithFunds generates new ed25519 signature scheme
// and requests some tokens from the UTXODB faucet.
// The amount of tokens is equal to utxodb.FundsFromFaucetAmount (=1000Mi) base tokens
// Returns signature scheme interface and public key in binary form
func (env *Solo) NewKeyPairWithFunds(seed ...*cryptolib.Seed) (*cryptolib.KeyPair, iotago.Address) {
	keyPair, addr := env.NewKeyPair(seed...)

	env.ledgerMutex.Lock()
	defer env.ledgerMutex.Unlock()

	_, err := env.utxoDB.GetFundsFromFaucet(addr)
	require.NoError(env.T, err)
	env.AssertL1BaseTokens(addr, utxodb.FundsFromFaucetAmount)

	return keyPair, addr
}

func (env *Solo) GetFundsFromFaucet(target iotago.Address, amount ...uint64) (*iotago.Transaction, error) {
	return env.utxoDB.GetFundsFromFaucet(target, amount...)
}

// NewSignatureSchemeAndPubKey generates new ed25519 signature scheme
// Returns signature scheme interface and public key in binary form
func (env *Solo) NewKeyPair(seedOpt ...*cryptolib.Seed) (*cryptolib.KeyPair, iotago.Address) {
	return testkey.GenKeyAddr(seedOpt...)
}
