package accounts

import (
	"fmt"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/kv"
	"github.com/nnikolash/wasp-types-exported/packages/kv/codec"
	"github.com/nnikolash/wasp-types-exported/packages/kv/collections"
	"github.com/nnikolash/wasp-types-exported/packages/util"
)

// observation: this uses the entire agentID as key, unlike acccounts.accountKey, which skips the chainID if it is the current chain. This means some bytes are wasted when saving NFTs

func NftsMapKey(agentID isc.AgentID) string {
	return PrefixNFTs + string(agentID.Bytes())
}

func NftsByCollectionMapKey(agentID isc.AgentID, collectionKey kv.Key) string {
	return PrefixNFTsByCollection + string(agentID.Bytes()) + string(collectionKey)
}

func FoundriesMapKey(agentID isc.AgentID) string {
	return PrefixFoundries + string(agentID.Bytes())
}

func AccountToNFTsMapR(state kv.KVStoreReader, agentID isc.AgentID) *collections.ImmutableMap {
	return collections.NewMapReadOnly(state, NftsMapKey(agentID))
}

func AccountToNFTsMap(state kv.KVStore, agentID isc.AgentID) *collections.Map {
	return collections.NewMap(state, NftsMapKey(agentID))
}

func NftToOwnerMap(state kv.KVStore) *collections.Map {
	return collections.NewMap(state, KeyNFTOwner)
}

func NftToOwnerMapR(state kv.KVStoreReader) *collections.ImmutableMap {
	return collections.NewMapReadOnly(state, KeyNFTOwner)
}

func NftCollectionKey(issuer iotago.Address) kv.Key {
	if issuer == nil {
		return NoCollection
	}
	nftAddr, ok := issuer.(*iotago.NFTAddress)
	if !ok {
		return NoCollection
	}
	id := nftAddr.NFTID()
	return kv.Key(id[:])
}

func NftsByCollectionMapR(state kv.KVStoreReader, agentID isc.AgentID, collectionKey kv.Key) *collections.ImmutableMap {
	return collections.NewMapReadOnly(state, NftsByCollectionMapKey(agentID, collectionKey))
}

func NftsByCollectionMap(state kv.KVStore, agentID isc.AgentID, collectionKey kv.Key) *collections.Map {
	return collections.NewMap(state, NftsByCollectionMapKey(agentID, collectionKey))
}

func hasNFT(state kv.KVStoreReader, agentID isc.AgentID, nftID iotago.NFTID) bool {
	return AccountToNFTsMapR(state, agentID).HasAt(nftID[:])
}

func removeNFTOwner(state kv.KVStore, nftID iotago.NFTID, agentID isc.AgentID) bool {
	// remove the mapping of NFTID => owner
	nftMap := NftToOwnerMap(state)
	if !nftMap.HasAt(nftID[:]) {
		return false
	}
	nftMap.DelAt(nftID[:])

	// add to the mapping of agentID => []NFTIDs
	nfts := AccountToNFTsMap(state, agentID)
	if !nfts.HasAt(nftID[:]) {
		return false
	}
	nfts.DelAt(nftID[:])
	return true
}

func setNFTOwner(state kv.KVStore, nftID iotago.NFTID, agentID isc.AgentID) {
	// add to the mapping of NFTID => owner
	nftMap := NftToOwnerMap(state)
	nftMap.SetAt(nftID[:], agentID.Bytes())

	// add to the mapping of agentID => []NFTIDs
	nfts := AccountToNFTsMap(state, agentID)
	nfts.SetAt(nftID[:], codec.EncodeBool(true))
}

func GetNFTData(state kv.KVStoreReader, nftID iotago.NFTID) *isc.NFT {
	o, oID := GetNFTOutput(state, nftID)
	if o == nil {
		return nil
	}
	owner, err := isc.AgentIDFromBytes(NftToOwnerMapR(state).GetAt(nftID[:]))
	if err != nil {
		panic("error parsing AgentID in NFTToOwnerMap")
	}
	return &isc.NFT{
		ID:       util.NFTIDFromNFTOutput(o, oID),
		Issuer:   o.ImmutableFeatureSet().IssuerFeature().Address,
		Metadata: o.ImmutableFeatureSet().MetadataFeature().Data,
		Owner:    owner,
	}
}

// CreditNFTToAccount credits an NFT to the on chain ledger
func CreditNFTToAccount(state kv.KVStore, agentID isc.AgentID, nftOutput *iotago.NFTOutput, chainID isc.ChainID) {
	if nftOutput.NFTID.Empty() {
		panic("empty NFTID")
	}

	issuerFeature := nftOutput.ImmutableFeatureSet().IssuerFeature()
	var issuer iotago.Address
	if issuerFeature != nil {
		issuer = issuerFeature.Address
	}
	creditNFTToAccount(state, agentID, nftOutput.NFTID, issuer)
	touchAccount(state, agentID, chainID)

	// save the NFTOutput with a temporary outputIndex so the NFTData is readily available (it will be updated upon block closing)
	SaveNFTOutput(state, nftOutput, 0)
}

func creditNFTToAccount(state kv.KVStore, agentID isc.AgentID, nftID iotago.NFTID, issuer iotago.Address) {
	setNFTOwner(state, nftID, agentID)

	collectionKey := NftCollectionKey(issuer)
	nftsByCollection := NftsByCollectionMap(state, agentID, collectionKey)
	nftsByCollection.SetAt(nftID[:], codec.EncodeBool(true))
}

// DebitNFTFromAccount removes an NFT from an account.
// If the account does not own the nft, it panics.
func DebitNFTFromAccount(state kv.KVStore, agentID isc.AgentID, nftID iotago.NFTID, chainID isc.ChainID) {
	nft := GetNFTData(state, nftID)
	if nft == nil {
		panic(fmt.Errorf("cannot debit unknown NFT %s", nftID.String()))
	}
	if !debitNFTFromAccount(state, agentID, nft) {
		panic(fmt.Errorf("cannot debit NFT %s from %s: %w", nftID.String(), agentID, ErrNotEnoughFunds))
	}
	touchAccount(state, agentID, chainID)
}

// DebitNFTFromAccount removes an NFT from the internal map of an account
func debitNFTFromAccount(state kv.KVStore, agentID isc.AgentID, nft *isc.NFT) bool {
	if !removeNFTOwner(state, nft.ID, agentID) {
		return false
	}

	collectionKey := NftCollectionKey(nft.Issuer)
	nftsByCollection := NftsByCollectionMap(state, agentID, collectionKey)
	if !nftsByCollection.HasAt(nft.ID[:]) {
		panic("inconsistency: NFT not found in collection")
	}
	nftsByCollection.DelAt(nft.ID[:])

	return true
}

func collectNFTIDs(m *collections.ImmutableMap) []iotago.NFTID {
	var ret []iotago.NFTID
	m.Iterate(func(idBytes []byte, val []byte) bool {
		id := iotago.NFTID{}
		copy(id[:], idBytes)
		ret = append(ret, id)
		return true
	})
	return ret
}

func getAccountNFTs(state kv.KVStoreReader, agentID isc.AgentID) []iotago.NFTID {
	return collectNFTIDs(AccountToNFTsMapR(state, agentID))
}

func getAccountNFTsInCollection(state kv.KVStoreReader, agentID isc.AgentID, collectionID iotago.NFTID) []iotago.NFTID {
	return collectNFTIDs(NftsByCollectionMapR(state, agentID, kv.Key(collectionID[:])))
}

func getL2TotalNFTs(state kv.KVStoreReader) []iotago.NFTID {
	return collectNFTIDs(NftToOwnerMapR(state))
}

// GetAccountNFTs returns all NFTs belonging to the agentID on the state
func GetAccountNFTs(state kv.KVStoreReader, agentID isc.AgentID) []iotago.NFTID {
	return getAccountNFTs(state, agentID)
}

func GetTotalL2NFTs(state kv.KVStoreReader) []iotago.NFTID {
	return getL2TotalNFTs(state)
}
