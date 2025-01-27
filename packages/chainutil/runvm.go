package chainutil

import (
	"errors"
	"time"

	"go.uber.org/zap"

	"github.com/samber/lo"

	"github.com/nnikolash/wasp-types-exported/packages/chain"
	"github.com/nnikolash/wasp-types-exported/packages/hashing"
	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/state"
	"github.com/nnikolash/wasp-types-exported/packages/state/indexedstore"
	"github.com/nnikolash/wasp-types-exported/packages/transaction"
	"github.com/nnikolash/wasp-types-exported/packages/vm"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/accounts"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/migrations"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/migrations/allmigrations"
	"github.com/nnikolash/wasp-types-exported/packages/vm/vmimpl"
)

func runISCTask(
	ch chain.ChainCore,
	aliasOutput *isc.AliasOutputWithID,
	blockTime time.Time,
	reqs []isc.Request,
	estimateGasMode bool,
	evmTracer *isc.EVMTracer,
) ([]*vm.RequestResult, error) {
	store := ch.Store()
	migs, err := getMigrationsForBlock(store, aliasOutput)
	if err != nil {
		return nil, err
	}
	task := &vm.VMTask{
		Processors:           ch.Processors(),
		AnchorOutput:         aliasOutput.GetAliasOutput(),
		AnchorOutputID:       aliasOutput.OutputID(),
		Store:                store,
		Requests:             reqs,
		TimeAssumption:       blockTime,
		Entropy:              hashing.PseudoRandomHash(nil),
		ValidatorFeeTarget:   accounts.CommonAccount(),
		EnableGasBurnLogging: estimateGasMode,
		EstimateGasMode:      estimateGasMode,
		EVMTracer:            evmTracer,
		Log:                  ch.Log().Desugar().WithOptions(zap.AddCallerSkip(1)).Sugar(),
		Migrations:           migs,
	}
	res, err := vmimpl.Run(task)
	if err != nil {
		return nil, err
	}
	return res.RequestResults, nil
}

func getMigrationsForBlock(store indexedstore.IndexedStore, aliasOutput *isc.AliasOutputWithID) (*migrations.MigrationScheme, error) {
	prevL1Commitment, err := transaction.L1CommitmentFromAliasOutput(aliasOutput.GetAliasOutput())
	if err != nil {
		panic(err)
	}
	prevState, err := store.StateByTrieRoot(prevL1Commitment.TrieRoot())
	if err != nil {
		if errors.Is(err, state.ErrTrieRootNotFound) {
			return allmigrations.DefaultScheme, nil
		}
		panic(err)
	}
	if lo.Must(store.LatestBlockIndex()) == prevState.BlockIndex() {
		return allmigrations.DefaultScheme, nil
	}
	newState := lo.Must(store.StateByIndex(prevState.BlockIndex() + 1))
	targetSchemaVersion := newState.SchemaVersion()
	return allmigrations.DefaultScheme.WithTargetSchemaVersion(targetSchemaVersion)
}

func runISCRequest(
	ch chain.ChainCore,
	aliasOutput *isc.AliasOutputWithID,
	blockTime time.Time,
	req isc.Request,
	estimateGasMode bool,
) (*vm.RequestResult, error) {
	results, err := runISCTask(
		ch,
		aliasOutput,
		blockTime,
		[]isc.Request{req},
		estimateGasMode,
		nil,
	)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, errors.New("request was skipped")
	}
	return results[0], nil
}
