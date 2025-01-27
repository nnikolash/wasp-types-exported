// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package chain

import (
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/nnikolash/wasp-types-exported/components/app"
	"github.com/nnikolash/wasp-types-exported/packages/apilib"
	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/kv/codec"
	"github.com/nnikolash/wasp-types-exported/packages/kv/dict"
	"github.com/nnikolash/wasp-types-exported/packages/origin"
	"github.com/nnikolash/wasp-types-exported/packages/parameters"
	"github.com/nnikolash/wasp-types-exported/packages/util"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/evm"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/governance"
	"github.com/nnikolash/wasp-types-exported/tools/wasp-cli/cli/cliclients"
	"github.com/nnikolash/wasp-types-exported/tools/wasp-cli/cli/config"
	"github.com/nnikolash/wasp-types-exported/tools/wasp-cli/cli/wallet"
	"github.com/nnikolash/wasp-types-exported/tools/wasp-cli/log"
	"github.com/nnikolash/wasp-types-exported/tools/wasp-cli/waspcmd"
)

func GetAllWaspNodes() []int {
	ret := []int{}
	for index := range viper.GetStringMap("wasp") {
		i, err := strconv.Atoi(index)
		log.Check(err)
		ret = append(ret, i)
	}
	return ret
}

func controllerAddrDefaultFallback(addr string) iotago.Address {
	if addr == "" {
		return wallet.Load().Address()
	}
	prefix, govControllerAddr, err := iotago.ParseBech32(addr)
	log.Check(err)
	if parameters.L1().Protocol.Bech32HRP != prefix {
		log.Fatalf("unexpected prefix. expected: %s, actual: %s", parameters.L1().Protocol.Bech32HRP, prefix)
	}
	return govControllerAddr
}

func initDeployCmd() *cobra.Command {
	var (
		node             string
		peers            []string
		quorum           int
		evmChainID       uint16
		blockKeepAmount  int32
		govControllerStr string
		chainName        string
	)

	cmd := &cobra.Command{
		Use:   "deploy --chain=<name>",
		Short: "Deploy a new chain",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			node = waspcmd.DefaultWaspNodeFallback(node)
			chainName = defaultChainFallback(chainName)

			if !util.IsSlug(chainName) {
				log.Fatalf("invalid chain name: %s, must be in slug format, only lowercase and hyphens, example: foo-bar", chainName)
			}

			l1Client := cliclients.L1Client()

			govController := controllerAddrDefaultFallback(govControllerStr)

			stateController := doDKG(node, peers, quorum)

			par := apilib.CreateChainParams{
				Layer1Client:         l1Client,
				CommitteeAPIHosts:    config.NodeAPIURLs([]string{node}),
				N:                    uint16(len(node)),
				T:                    uint16(quorum),
				OriginatorKeyPair:    wallet.Load(),
				Textout:              os.Stdout,
				GovernanceController: govController,
				InitParams: dict.Dict{
					origin.ParamChainOwner:      isc.NewAgentID(govController).Bytes(),
					origin.ParamEVMChainID:      codec.EncodeUint16(evmChainID),
					origin.ParamBlockKeepAmount: codec.EncodeInt32(blockKeepAmount),
					origin.ParamWaspVersion:     codec.EncodeString(app.Version),
				},
			}

			chainID, err := apilib.DeployChain(par, stateController, govController)
			log.Check(err)

			config.AddChain(chainName, chainID.String())

			activateChain(node, chainName, chainID)
		},
	}

	waspcmd.WithWaspNodeFlag(cmd, &node)
	waspcmd.WithPeersFlag(cmd, &peers)
	cmd.Flags().Uint16VarP(&evmChainID, "evm-chainid", "", evm.DefaultChainID, "ChainID")
	cmd.Flags().Int32VarP(&blockKeepAmount, "block-keep-amount", "", governance.DefaultBlockKeepAmount, "Amount of blocks to keep in the blocklog (-1 to keep all blocks)")
	cmd.Flags().StringVar(&chainName, "chain", "", "name of the chain")
	log.Check(cmd.MarkFlagRequired("chain"))
	cmd.Flags().IntVar(&quorum, "quorum", 0, "quorum (default: 3/4s of the number of committee nodes)")
	cmd.Flags().StringVar(&govControllerStr, "gov-controller", "", "governance controller address")
	return cmd
}
