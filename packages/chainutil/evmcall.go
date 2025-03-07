package chainutil

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum"

	"github.com/nnikolash/wasp-types-exported/packages/chain"
	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/parameters"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/evm"
	"github.com/nnikolash/wasp-types-exported/packages/vm/gas"
)

// EVMCall executes an EVM contract call and returns its output, discarding any state changes
func EVMCall(ch chain.ChainCore, aliasOutput *isc.AliasOutputWithID, call ethereum.CallMsg) ([]byte, error) {
	info := getChainInfo(ch)

	// 0 means view call
	gasLimit := gas.EVMCallGasLimit(info.GasLimits, &info.GasFeePolicy.EVMGasRatio)
	if call.Gas != 0 && call.Gas > gasLimit {
		call.Gas = gasLimit
	}

	if call.GasPrice == nil {
		call.GasPrice = info.GasFeePolicy.DefaultGasPriceFullDecimals(parameters.L1().BaseToken.Decimals)
	}

	iscReq := isc.NewEVMOffLedgerCallRequest(ch.ID(), call)
	// TODO: setting EstimateGasMode = true feels wrong here
	res, err := runISCRequest(ch, aliasOutput, time.Now(), iscReq, true)
	if err != nil {
		return nil, err
	}
	if res.Receipt.Error != nil {
		vmerr, resolvingErr := ResolveError(ch, res.Receipt.Error)
		if resolvingErr != nil {
			panic(fmt.Errorf("error resolving vmerror %w", resolvingErr))
		}
		return nil, vmerr
	}
	return res.Return[evm.FieldResult], nil
}
