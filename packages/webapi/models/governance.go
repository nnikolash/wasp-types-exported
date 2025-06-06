package models

import "github.com/nnikolash/wasp-types-exported/packages/vm/gas"

/*
Both Gov* structs should be removed at some point.
The corecontract implementations should be moved outside the webapi, therefore using the webapi ChainInfo/Metadata structs should be avoided.
*/
type GovPublicChainMetadata struct {
	EVMJsonRPCURL   string `json:"evmJsonRpcURL" swagger:"desc(The EVM json rpc url),required"`
	EVMWebSocketURL string `json:"evmWebSocketURL" swagger:"desc(The EVM websocket url)),required"`

	Name        string `json:"name" swagger:"desc(The name of the chain),required"`
	Description string `json:"description" swagger:"desc(The description of the chain.),required"`
	Website     string `json:"website" swagger:"desc(The official website of the chain.),required"`
}

type GovChainInfoResponse struct {
	ChainID      string                 `json:"chainID" swagger:"desc(ChainID (Bech32-encoded).),required"`
	ChainOwnerID string                 `json:"chainOwnerId" swagger:"desc(The chain owner address (Bech32-encoded).),required"`
	GasFeePolicy *gas.FeePolicy         `json:"gasFeePolicy" swagger:"desc(The gas fee policy),required"`
	GasLimits    *gas.Limits            `json:"gasLimits" swagger:"desc(The gas limits),required"`
	PublicURL    string                 `json:"publicURL" swagger:"desc(The fully qualified public url leading to the chains metadata),required"`
	Metadata     GovPublicChainMetadata `json:"metadata" swagger:"desc(The metadata of the chain),required"`
}

type GovAllowedStateControllerAddressesResponse struct {
	Addresses []string `json:"addresses" swagger:"desc(The allowed state controller addresses (Bech32-encoded))"`
}

type GovChainOwnerResponse struct {
	ChainOwner string `json:"chainOwner" swagger:"desc(The chain owner (Bech32-encoded))"`
}
