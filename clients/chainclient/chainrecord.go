package chainclient

import (
	"context"

	"github.com/nnikolash/wasp-types-exported/clients/apiclient"
)

// GetChainRecord fetches the chain's Record
func (c *Client) GetChainRecord(ctx context.Context) (*apiclient.ChainInfoResponse, error) {
	chainInfo, _, err := c.WaspClient.ChainsApi.GetChainInfo(ctx, c.ChainID.String()).Execute()

	return chainInfo, err
}
