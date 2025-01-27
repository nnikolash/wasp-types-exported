package testutil

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/nnikolash/wasp-types-exported/packages/cryptolib"
	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/kv/dict"
	"github.com/nnikolash/wasp-types-exported/packages/testutil/testkey"
	"github.com/nnikolash/wasp-types-exported/packages/vm/gas"
)

func DummyOffledgerRequest(chainID isc.ChainID) isc.OffLedgerRequest {
	contract := isc.Hn("somecontract")
	entrypoint := isc.Hn("someentrypoint")
	args := dict.Dict{}
	req := isc.NewOffLedgerRequest(chainID, contract, entrypoint, args, 0, gas.LimitsDefault.MaxGasPerRequest)
	keys, _ := testkey.GenKeyAddr()
	return req.Sign(keys)
}

func DummyOffledgerRequestForAccount(chainID isc.ChainID, nonce uint64, kp *cryptolib.KeyPair) isc.OffLedgerRequest {
	contract := isc.Hn("somecontract")
	entrypoint := isc.Hn("someentrypoint")
	args := dict.Dict{}
	req := isc.NewOffLedgerRequest(chainID, contract, entrypoint, args, nonce, gas.LimitsDefault.MaxGasPerRequest)
	return req.Sign(kp)
}

func DummyEVMRequest(chainID isc.ChainID, gasPrice *big.Int) isc.OffLedgerRequest {
	key, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}

	tx := types.MustSignNewTx(key, types.NewEIP155Signer(big.NewInt(0)),
		&types.LegacyTx{
			Nonce:    0,
			To:       &common.MaxAddress,
			Value:    big.NewInt(123),
			Gas:      10000,
			GasPrice: gasPrice,
			Data:     []byte{},
		})

	req, err := isc.NewEVMOffLedgerTxRequest(chainID, tx)
	if err != nil {
		panic(err)
	}
	return req
}
