package services

import (
	"github.com/nnikolash/wasp-types-exported/packages/chains"
	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/registry"
	"github.com/nnikolash/wasp-types-exported/packages/webapi/interfaces"
)

type RegistryService struct {
	chainsProvider              chains.Provider
	chainRecordRegistryProvider registry.ChainRecordRegistryProvider
}

func NewRegistryService(chainsProvider chains.Provider, chainRecordRegistryProvider registry.ChainRecordRegistryProvider) interfaces.RegistryService {
	return &RegistryService{
		chainsProvider:              chainsProvider,
		chainRecordRegistryProvider: chainRecordRegistryProvider,
	}
}

func (c *RegistryService) GetChainRecordByChainID(chainID isc.ChainID) (*registry.ChainRecord, error) {
	return c.chainRecordRegistryProvider.ChainRecord(chainID)
}
