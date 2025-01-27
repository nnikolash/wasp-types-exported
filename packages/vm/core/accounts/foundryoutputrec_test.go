package accounts

import (
	"math/big"
	"testing"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/nnikolash/wasp-types-exported/packages/util/rwutil"
)

func TestFoundryOutputRecSerialization(t *testing.T) {
	o := FoundryOutputRec{
		OutputID: iotago.OutputID{1, 2, 3},
		Amount:   300,
		TokenScheme: &iotago.SimpleTokenScheme{
			MaximumSupply: big.NewInt(1000),
			MintedTokens:  big.NewInt(20),
			MeltedTokens:  big.NewInt(1),
		},
		Metadata: []byte("Tralala"),
	}
	rwutil.ReadWriteTest(t, &o, new(FoundryOutputRec))
	rwutil.BytesTest(t, &o, FoundryOutputRecFromBytes)
}
