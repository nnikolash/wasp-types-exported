package cryptolib_test

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nnikolash/wasp-types-exported/packages/cryptolib"
	"github.com/nnikolash/wasp-types-exported/packages/util/rwutil"
)

func TestPrivateKeySerialization(t *testing.T) {
	seedBytes := make([]byte, cryptolib.SeedSize)
	rand.Read(seedBytes)
	pivkey1 := cryptolib.PrivateKeyFromSeed((cryptolib.SeedFromBytes(seedBytes)))
	pivkey2, err := cryptolib.PrivateKeyFromBytes(pivkey1.AsBytes())
	require.NoError(t, err)
	require.Equal(t, pivkey1, pivkey2)

	rwutil.ReadWriteTest(t, pivkey1, cryptolib.NewPrivateKey())
}
