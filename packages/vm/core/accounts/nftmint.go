package accounts

import (
	"io"
	"slices"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/kv"
	"github.com/nnikolash/wasp-types-exported/packages/kv/codec"
	"github.com/nnikolash/wasp-types-exported/packages/kv/collections"
	"github.com/nnikolash/wasp-types-exported/packages/kv/dict"
	"github.com/nnikolash/wasp-types-exported/packages/util/rwutil"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/errors/coreerrors"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/evm"
)

type MintedNFTRecord struct {
	positionInMintedList uint16
	outputIndex          uint16
	owner                isc.AgentID
	output               *iotago.NFTOutput
}

func (rec *MintedNFTRecord) Read(r io.Reader) error {
	rr := rwutil.NewReader(r)
	rec.positionInMintedList = rr.ReadUint16()
	rec.outputIndex = rr.ReadUint16()
	rec.owner = isc.AgentIDFromReader(rr)
	rec.output = new(iotago.NFTOutput)
	rr.ReadSerialized(rec.output)
	return rr.Err
}

func (rec *MintedNFTRecord) Write(w io.Writer) error {
	ww := rwutil.NewWriter(w)
	ww.WriteUint16(rec.positionInMintedList)
	ww.WriteUint16(rec.outputIndex)
	if rec.owner != nil {
		ww.Write(rec.owner)
	} else {
		ww.Write(&isc.NilAgentID{})
	}
	ww.WriteSerialized(rec.output)
	return ww.Err
}

func (rec *MintedNFTRecord) Bytes() []byte {
	return rwutil.WriteToBytes(rec)
}

func MintedNFTRecordFromBytes(data []byte) *MintedNFTRecord {
	record, err := rwutil.ReadFromBytes(data, new(MintedNFTRecord))
	if err != nil {
		panic(err)
	}
	return record
}

func NewlyMintedNFTsMap(state kv.KVStore) *collections.Map {
	return collections.NewMap(state, PrefixNewlyMintedNFTs)
}

func MintIDMap(state kv.KVStore) *collections.Map {
	return collections.NewMap(state, PrefixMintIDMap)
}

func MintIDMapR(state kv.KVStoreReader) *collections.ImmutableMap {
	return collections.NewMapReadOnly(state, PrefixMintIDMap)
}

var (
	errMintNFTWithdraw      = coreerrors.Register("can only withdraw on mint to a L1 address").Create()
	errInvalidAgentID       = coreerrors.Register("invalid agentID").Create()
	errCollectionNotAllowed = coreerrors.Register("caller doesn't own the collection").Create()
)

type mintParameters struct {
	immutableMetadata []byte
	targetAddress     iotago.Address
	issuerAddress     iotago.Address
	ownerAgentID      isc.AgentID
	withdrawOnMint    bool
	collectionID      iotago.NFTID
}

func mintParams(ctx isc.Sandbox) mintParameters {
	params := ctx.Params()

	immutableMetadata := params.MustGetBytes(ParamNFTImmutableData)
	targetAgentID := params.MustGetAgentID(ParamAgentID)
	withdrawOnMint := params.MustGetBool(ParamNFTWithdrawOnMint, false)
	emptyNFTID := iotago.NFTID{}
	collectionID := params.MustGetNFTID(ParamCollectionID, emptyNFTID)

	chainAddress := ctx.ChainID().AsAddress()
	ret := mintParameters{
		immutableMetadata: slices.Clone(immutableMetadata),
		targetAddress:     chainAddress,
		issuerAddress:     chainAddress,
		ownerAgentID:      targetAgentID,
		withdrawOnMint:    withdrawOnMint,
		collectionID:      collectionID,
	}

	if !collectionID.Empty() {
		// assert the NFT of collectionID is on-chain and owned by the caller
		if !hasNFT(ctx.State(), ctx.Caller(), collectionID) {
			panic(errCollectionNotAllowed)
		}
		ret.issuerAddress = collectionID.ToAddress()
	}

	switch targetAgentID.Kind() {
	case isc.AgentIDKindContract, isc.AgentIDKindEthereumAddress:
		if withdrawOnMint {
			panic(errMintNFTWithdraw)
		}
		return ret
	case isc.AgentIDKindAddress:
		if withdrawOnMint {
			ret.targetAddress = targetAgentID.(*isc.AddressAgentID).Address()
			return ret
		}
		return ret
	default:
		panic(errInvalidAgentID)
	}
}

func mintID(blockIndex uint32, positionInMintedList uint16) []byte {
	ret := make([]byte, 6)
	copy(ret[0:], codec.EncodeUint32(blockIndex))
	copy(ret[4:], codec.EncodeUint16(positionInMintedList))
	return ret
}

// NFTs are always minted with the minimumSD and that must be provided via allowance
func mintNFT(ctx isc.Sandbox) dict.Dict {
	params := mintParams(ctx)

	// NFTs are now automatically registered inside the EVM.
	// The EVM requires IRC27 metadata to be present. Therefore, any invalid metadata will panic here
	// This will not check the metadata according to the schema, only syntactic validation applies. "{}" would be correct.
	_, err := isc.IRC27NFTMetadataFromBytes(params.immutableMetadata)
	if err != nil {
		panic(ErrImmutableMetadataInvalid.Create(err.Error()))
	}

	positionInMintedList, nftOutput := ctx.Privileged().MintNFT(
		params.targetAddress,
		params.immutableMetadata,
		params.issuerAddress,
	)

	// debit the SD required for the NFT from the sender account
	ctx.TransferAllowedFunds(ctx.AccountID(), isc.NewAssetsBaseTokens(nftOutput.Amount))                                          // claim tokens from allowance
	DebitFromAccount(ctx.SchemaVersion(), ctx.State(), ctx.AccountID(), isc.NewAssetsBaseTokens(nftOutput.Amount), ctx.ChainID()) // debit from this SC account

	rec := MintedNFTRecord{
		positionInMintedList: positionInMintedList,
		outputIndex:          0, // to be filled on block close by `SaveMintedNFTOutput`
		owner:                params.ownerAgentID,
		output:               nftOutput,
	}
	// save the info required to credit the NFT on next block
	NewlyMintedNFTsMap(ctx.State()).SetAt(codec.Encode(positionInMintedList), rec.Bytes())

	// register the collection in the EVM
	if !params.collectionID.Empty() {
		res := ctx.CallView(
			evm.Contract.Hname(),
			evm.FuncGetERC721CollectionAddress.Hname(),
			dict.Dict{evm.FieldNFTCollectionID: codec.Encode(params.collectionID)},
		)
		exists := codec.MustDecodeBool(res.Get(evm.FieldResult))

		if !exists {
			// NOTE must not call `RegisterERC721NFTCollection` if already registered, otherwise it will panic
			ctx.Call(
				evm.Contract.Hname(),
				evm.FuncRegisterERC721NFTCollection.Hname(),
				dict.Dict{evm.FieldNFTCollectionID: codec.Encode(params.collectionID)},
				nil,
			)
		}
	}

	return dict.Dict{
		ParamMintID: mintID(ctx.StateAnchor().StateIndex+1, positionInMintedList),
	}
}

func viewNFTIDbyMintID(ctx isc.SandboxView) dict.Dict {
	internalMintID := ctx.Params().MustGetBytes(ParamMintID)
	nftID := MintIDMapR(ctx.StateR()).GetAt(internalMintID)
	return dict.Dict{
		ParamNFTID: nftID,
	}
}

// ----  output management

func SaveMintedNFTOutput(state kv.KVStore, positionInMintedList, outputIndex uint16) {
	mintMap := NewlyMintedNFTsMap(state)
	key := codec.Encode(positionInMintedList)
	recBytes := mintMap.GetAt(key)
	if recBytes == nil {
		return
	}
	rec := MintedNFTRecordFromBytes(recBytes)
	rec.outputIndex = outputIndex
	mintMap.SetAt(key, rec.Bytes())
}

func updateNewlyMintedNFTOutputIDs(state kv.KVStore, anchorTxID iotago.TransactionID, blockIndex uint32) map[iotago.NFTID]isc.AgentID {
	mintMap := NewlyMintedNFTsMap(state)
	nftMap := NFTOutputMap(state)
	newNFTIDs := make(map[iotago.NFTID]isc.AgentID, 0)

	// iterate the minted collection of NFT's that we're looking to co-relate with their NFTIDs
	mintMap.Iterate(func(_, recBytes []byte) bool {
		mintedRec := MintedNFTRecordFromBytes(recBytes)
		// calculate the NFTID from the anchor txID	+ outputIndex
		outputID := iotago.OutputIDFromTransactionIDAndIndex(anchorTxID, mintedRec.outputIndex)
		nftID := iotago.NFTIDFromOutputID(outputID)

		if mintedRec.owner.Kind() != isc.AgentIDKindNil { // when owner is nil, means the NFT was minted directly to a L1 wallet
			outputRec := NFTOutputRec{
				OutputID: outputID,
				Output:   mintedRec.output,
			}
			// save the updated data in the NFT map
			nftMap.SetAt(nftID[:], outputRec.Bytes())
			// credit the NFT to the target owner
			creditNFTToAccount(state, mintedRec.owner, nftID, mintedRec.output.ImmutableFeatureSet().IssuerFeature().Address)
		}
		// save the mapping of [mintID => NFTID]
		MintIDMap(state).SetAt(mintID(blockIndex, mintedRec.positionInMintedList), nftID[:])
		newNFTIDs[nftID] = mintedRec.owner

		return true
	})
	// clear the minted collection
	mintMap.Erase()

	return newNFTIDs
}
