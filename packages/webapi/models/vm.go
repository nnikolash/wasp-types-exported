package models

import (
	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/vm/gas"
)

type ReceiptResponse struct {
	Request       isc.RequestJSON            `json:"request" swagger:"required"`
	RawError      *isc.UnresolvedVMErrorJSON `json:"rawError,omitempty"`
	ErrorMessage  string                     `json:"errorMessage,omitempty"`
	GasBudget     string                     `json:"gasBudget" swagger:"required,desc(The gas budget (uint64 as string))"`
	GasBurned     string                     `json:"gasBurned" swagger:"required,desc(The burned gas (uint64 as string))"`
	GasFeeCharged string                     `json:"gasFeeCharged" swagger:"required,desc(The charged gas fee (uint64 as string))"`
	SDCharged     string                     `json:"storageDepositCharged" swagger:"required,desc(Storage deposit charged (uint64 as string))"`
	BlockIndex    uint32                     `json:"blockIndex" swagger:"required,min(1)"`
	RequestIndex  uint16                     `json:"requestIndex" swagger:"required,min(1)"`
	GasBurnLog    []gas.BurnRecord           `json:"gasBurnLog" swagger:"required"`
}

func MapReceiptResponse(receipt *isc.Receipt) *ReceiptResponse {
	burnRecords := make([]gas.BurnRecord, 0)

	if receipt.GasBurnLog != nil {
		burnRecords = append(burnRecords, receipt.GasBurnLog.Records...)
	}

	req, err := isc.RequestFromBytes(receipt.Request)
	if err != nil {
		panic(err)
	}

	return &ReceiptResponse{
		Request:       isc.RequestToJSONObject(req),
		RawError:      receipt.Error.ToJSONStruct(),
		ErrorMessage:  receipt.ResolvedError,
		BlockIndex:    receipt.BlockIndex,
		RequestIndex:  receipt.RequestIndex,
		GasBudget:     iotago.EncodeUint64(receipt.GasBudget),
		GasBurned:     iotago.EncodeUint64(receipt.GasBurned),
		GasFeeCharged: iotago.EncodeUint64(receipt.GasFeeCharged),
		SDCharged:     iotago.EncodeUint64(receipt.SDCharged),
		GasBurnLog:    burnRecords,
	}
}
