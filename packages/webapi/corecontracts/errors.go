package corecontracts

import (
	"github.com/nnikolash/wasp-types-exported/packages/chain"
	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/kv/codec"
	"github.com/nnikolash/wasp-types-exported/packages/kv/kvdecoder"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/errors"
	"github.com/nnikolash/wasp-types-exported/packages/webapi/common"
)

func ErrorMessageFormat(ch chain.Chain, contractID isc.Hname, errorID uint16, blockIndexOrTrieRoot string) (string, error) {
	errorCode := isc.NewVMErrorCode(contractID, errorID)

	ret, err := common.CallView(
		ch,
		errors.Contract.Hname(),
		errors.ViewGetErrorMessageFormat.Hname(),
		codec.MakeDict(map[string]interface{}{errors.ParamErrorCode: errorCode.Bytes()}),
		blockIndexOrTrieRoot,
	)
	if err != nil {
		return "", err
	}

	resultDecoder := kvdecoder.New(ret)
	messageFormat, err := resultDecoder.GetString(errors.ParamErrorMessageFormat)
	if err != nil {
		return "", err
	}

	return messageFormat, nil
}
