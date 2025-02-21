package isc

import (
	"fmt"
	"io"
	"time"

	"github.com/ethereum/go-ethereum"

	"github.com/iotaledger/hive.go/serializer/v2"
	iotago "github.com/iotaledger/iota.go/v3"

	"github.com/nnikolash/wasp-types-exported/packages/kv/dict"
	"github.com/nnikolash/wasp-types-exported/packages/util"
	"github.com/nnikolash/wasp-types-exported/packages/util/rwutil"
)

type OnLedgerRequestData struct {
	outputID iotago.OutputID
	output   iotago.Output

	// the following originate from UTXOMetaData and output, and are created in `NewExtendedOutputData`

	featureBlocks    iotago.FeatureSet
	unlockConditions iotago.UnlockConditionSet
	requestMetadata  *RequestMetadata
}

var (
	_ Request         = new(OnLedgerRequestData)
	_ OnLedgerRequest = new(OnLedgerRequestData)
	_ Calldata        = new(OnLedgerRequestData)
	_ Features        = new(OnLedgerRequestData)
)

func OnLedgerFromUTXO(output iotago.Output, outputID iotago.OutputID) (OnLedgerRequest, error) {
	r := &OnLedgerRequestData{}
	if err := r.readFromUTXO(output, outputID); err != nil {
		return nil, err
	}
	return r, nil
}

func (req *OnLedgerRequestData) RequestMetadataRaw() *RequestMetadata {
	return req.requestMetadata
}

func (req *OnLedgerRequestData) readFromUTXO(output iotago.Output, outputID iotago.OutputID) error {
	var reqMetadata *RequestMetadata
	var err error

	fbSet := output.FeatureSet()

	reqMetadata, err = requestMetadataFromFeatureSet(fbSet)
	if err != nil {
		reqMetadata = nil // bad metadata. // we must handle these request, so that those funds are not lost forever
	}

	if reqMetadata != nil {
		reqMetadata.Allowance.fillEmptyNFTIDs(output, outputID)
	}

	req.output = output
	req.outputID = outputID
	req.featureBlocks = fbSet
	req.unlockConditions = output.UnlockConditionSet()
	req.requestMetadata = reqMetadata
	return nil
}

func (req *OnLedgerRequestData) Read(r io.Reader) error {
	rr := rwutil.NewReader(r)
	rr.ReadKindAndVerify(rwutil.Kind(requestKindOnLedger))
	rr.ReadN(req.outputID[:])
	outputData := rr.ReadBytes()
	if rr.Err != nil {
		return rr.Err
	}
	req.output, rr.Err = util.OutputFromBytes(outputData)
	if rr.Err != nil {
		return rr.Err
	}
	return req.readFromUTXO(req.output, req.outputID)
}

func (req *OnLedgerRequestData) Write(w io.Writer) error {
	ww := rwutil.NewWriter(w)
	ww.WriteKind(rwutil.Kind(requestKindOnLedger))
	ww.WriteN(req.outputID[:])
	if ww.Err != nil {
		return ww.Err
	}
	outputData, err := req.output.Serialize(serializer.DeSeriModePerformLexicalOrdering, nil)
	ww.Err = err
	ww.WriteBytes(outputData)
	return ww.Err
}

func (req *OnLedgerRequestData) Allowance() *Assets {
	if req.requestMetadata == nil {
		return NewEmptyAssets()
	}
	return req.requestMetadata.Allowance
}

func (req *OnLedgerRequestData) Assets() *Assets {
	amount := req.output.Deposit()
	tokens := req.output.NativeTokenList()
	ret := NewAssets(amount, tokens)
	NFT := req.NFT()
	if NFT != nil {
		ret.AddNFTs(NFT.ID)
	}
	return ret
}

func (req *OnLedgerRequestData) Bytes() []byte {
	return rwutil.WriteToBytes(req)
}

func (req *OnLedgerRequestData) CallTarget() CallTarget {
	if req.requestMetadata == nil {
		return CallTarget{}
	}
	return CallTarget{
		Contract:   req.requestMetadata.TargetContract,
		EntryPoint: req.requestMetadata.EntryPoint,
	}
}

func (req *OnLedgerRequestData) Clone() OnLedgerRequest {
	outputID := iotago.OutputID{}
	copy(outputID[:], req.outputID[:])

	ret := &OnLedgerRequestData{
		outputID:         outputID,
		output:           req.output.Clone(),
		featureBlocks:    req.featureBlocks.Clone(),
		unlockConditions: util.CloneMap(req.unlockConditions),
	}
	if req.requestMetadata != nil {
		ret.requestMetadata = req.requestMetadata.Clone()
	}
	return ret
}

func (req *OnLedgerRequestData) Expiry() (time.Time, iotago.Address) {
	expiration := req.unlockConditions.Expiration()
	if expiration == nil {
		return time.Time{}, nil
	}

	return time.Unix(int64(expiration.UnixTime), 0), expiration.ReturnAddress
}

func (req *OnLedgerRequestData) Features() Features {
	return req
}

func (req *OnLedgerRequestData) GasBudget() (gasBudget uint64, isEVM bool) {
	if req.requestMetadata == nil {
		return 0, false
	}
	return req.requestMetadata.GasBudget, false
}

func (req *OnLedgerRequestData) ID() RequestID {
	return RequestID(req.outputID)
}

// IsInternalUTXO if true the output cannot be interpreted as a request
func (req *OnLedgerRequestData) IsInternalUTXO(chainID ChainID) bool {
	if req.output.Type() == iotago.OutputFoundry {
		return true
	}
	if req.senderAddress() == nil {
		return false
	}
	if !req.senderAddress().Equal(chainID.AsAddress()) {
		return false
	}
	if req.requestMetadata != nil {
		return false
	}
	return true
}

func (req *OnLedgerRequestData) IsOffLedger() bool {
	return false
}

func (req *OnLedgerRequestData) NFT() *NFT {
	nftOutput, ok := req.output.(*iotago.NFTOutput)
	if !ok {
		return nil
	}

	ret := &NFT{}

	ret.ID = util.NFTIDFromNFTOutput(nftOutput, req.OutputID())

	for _, featureBlock := range nftOutput.ImmutableFeatures {
		if block, ok := featureBlock.(*iotago.IssuerFeature); ok {
			ret.Issuer = block.Address
		}
		if block, ok := featureBlock.(*iotago.MetadataFeature); ok {
			ret.Metadata = block.Data
		}
	}

	return ret
}

func (req *OnLedgerRequestData) Output() iotago.Output {
	return req.output
}

func (req *OnLedgerRequestData) OutputID() iotago.OutputID {
	return req.outputID
}

func (req *OnLedgerRequestData) Params() dict.Dict {
	if req.requestMetadata == nil {
		return dict.Dict{}
	}
	return req.requestMetadata.Params
}

func (req *OnLedgerRequestData) ReturnAmount() (uint64, bool) {
	storageDepositReturn := req.unlockConditions.StorageDepositReturn()
	if storageDepositReturn == nil {
		return 0, false
	}
	return storageDepositReturn.Amount, true
}

func (req *OnLedgerRequestData) SenderAccount() AgentID {
	sender := req.senderAddress()
	if sender == nil {
		return nil
	}
	if req.requestMetadata != nil && !req.requestMetadata.SenderContract.Empty() {
		if sender.Type() == iotago.AddressAlias {
			chainID := ChainIDFromAddress(sender.(*iotago.AliasAddress))
			return req.requestMetadata.SenderContract.AgentID(chainID)
		}
	}
	return NewAgentID(sender)
}

func (req *OnLedgerRequestData) senderAddress() iotago.Address {
	senderBlock := req.featureBlocks.SenderFeature()
	if senderBlock == nil {
		return nil
	}
	return senderBlock.Address
}

func (req *OnLedgerRequestData) String() string {
	metadata := req.requestMetadata
	if metadata == nil {
		return "onledger request without metadata"
	}
	return fmt.Sprintf("onLedgerRequestData::{ ID: %s, sender: %s, target: %s, entrypoint: %s, Params: %s, GasBudget: %d }",
		req.ID().String(),
		metadata.SenderContract.String(),
		metadata.TargetContract.String(),
		metadata.EntryPoint.String(),
		metadata.Params.String(),
		metadata.GasBudget,
	)
}

func (req *OnLedgerRequestData) TargetAddress() iotago.Address {
	switch out := req.output.(type) {
	case *iotago.BasicOutput:
		return out.Ident()
	case *iotago.FoundryOutput:
		return out.Ident()
	case *iotago.NFTOutput:
		return out.Ident()
	case *iotago.AliasOutput:
		return out.AliasID.ToAddress()
	default:
		panic("onLedgerRequestData:TargetAddress implement me")
	}
}

func (req *OnLedgerRequestData) TimeLock() time.Time {
	timelock := req.unlockConditions.Timelock()
	if timelock == nil {
		return time.Time{}
	}
	return time.Unix(int64(timelock.UnixTime), 0)
}

func (req *OnLedgerRequestData) EVMCallMsg() *ethereum.CallMsg {
	return nil
}

// region RetryOnLedgerRequest //////////////////////////////////////////////////////////////////

type RetryOnLedgerRequest struct {
	OnLedgerRequest
	retryOutputID iotago.OutputID
}

func NewRetryOnLedgerRequest(req OnLedgerRequest, retryOutput iotago.OutputID) *RetryOnLedgerRequest {
	return &RetryOnLedgerRequest{
		OnLedgerRequest: req,
		retryOutputID:   retryOutput,
	}
}

func (r *RetryOnLedgerRequest) RetryOutputID() iotago.OutputID {
	return r.retryOutputID
}

func (r *RetryOnLedgerRequest) SetRetryOutputID(oid iotago.OutputID) {
	r.retryOutputID = oid
}
