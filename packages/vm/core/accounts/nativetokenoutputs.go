package accounts

import (
	iotago "github.com/iotaledger/iota.go/v3"

	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/kv"
	"github.com/nnikolash/wasp-types-exported/packages/kv/collections"
)

func newNativeTokensArray(state kv.KVStore) *collections.Array {
	return collections.NewArray(state, keyNewNativeTokens)
}

func NativeTokenOutputMap(state kv.KVStore) *collections.Map {
	return collections.NewMap(state, keyNativeTokenOutputMap)
}

func nativeTokenOutputMapR(state kv.KVStoreReader) *collections.ImmutableMap {
	return collections.NewMapReadOnly(state, keyNativeTokenOutputMap)
}

// SaveNativeTokenOutput map nativeTokenID -> foundryRec
func SaveNativeTokenOutput(state kv.KVStore, out *iotago.BasicOutput, outputIndex uint16) {
	tokenRec := NativeTokenOutputRec{
		// TransactionID is unknown yet, will be filled next block
		OutputID:          iotago.OutputIDFromTransactionIDAndIndex(iotago.TransactionID{}, outputIndex),
		StorageBaseTokens: out.Amount,
		Amount:            out.NativeTokens[0].Amount,
	}
	NativeTokenOutputMap(state).SetAt(out.NativeTokens[0].ID[:], tokenRec.Bytes())
	newNativeTokensArray(state).Push(out.NativeTokens[0].ID[:])
}

func updateNativeTokenOutputIDs(state kv.KVStore, anchorTxID iotago.TransactionID) {
	newNativeTokens := newNativeTokensArray(state)
	allNativeTokens := NativeTokenOutputMap(state)
	n := newNativeTokens.Len()
	for i := uint32(0); i < n; i++ {
		k := newNativeTokens.GetAt(i)
		rec := MustNativeTokenOutputRecFromBytes(allNativeTokens.GetAt(k))
		rec.OutputID = iotago.OutputIDFromTransactionIDAndIndex(anchorTxID, rec.OutputID.Index())
		allNativeTokens.SetAt(k, rec.Bytes())
	}
	newNativeTokens.Erase()
}

func DeleteNativeTokenOutput(state kv.KVStore, nativeTokenID iotago.NativeTokenID) {
	NativeTokenOutputMap(state).DelAt(nativeTokenID[:])
}

func GetNativeTokenOutput(state kv.KVStoreReader, nativeTokenID iotago.NativeTokenID, chainID isc.ChainID) (*iotago.BasicOutput, iotago.OutputID) {
	data := nativeTokenOutputMapR(state).GetAt(nativeTokenID[:])
	if data == nil {
		return nil, iotago.OutputID{}
	}
	tokenRec := MustNativeTokenOutputRecFromBytes(data)
	ret := &iotago.BasicOutput{
		Amount: tokenRec.StorageBaseTokens,
		NativeTokens: iotago.NativeTokens{{
			ID:     nativeTokenID,
			Amount: tokenRec.Amount,
		}},
		Conditions: iotago.UnlockConditions{
			&iotago.AddressUnlockCondition{Address: chainID.AsAddress()},
		},
		Features: iotago.Features{
			&iotago.SenderFeature{
				Address: chainID.AsAddress(),
			},
		},
	}
	return ret, tokenRec.OutputID
}
