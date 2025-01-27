package chainclient

import (
	"context"

	"github.com/nnikolash/wasp-types-exported/clients/apiclient"
	"github.com/nnikolash/wasp-types-exported/clients/apiextensions"
	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/kv/dict"
)

// CallView sends a request to call a view function of a given contract, and returns the result of the call
func (c *Client) CallView(ctx context.Context, hContract isc.Hname, functionName string, args dict.Dict, blockNumberOrHash ...string) (dict.Dict, error) {
	viewCall := apiclient.ContractCallViewRequest{
		ContractHName: hContract.String(),
		FunctionName:  functionName,
		Arguments:     apiextensions.JSONDictToAPIJSONDict(args.JSONDict()),
	}
	if len(blockNumberOrHash) > 0 {
		viewCall.Block = &blockNumberOrHash[0]
	}

	return apiextensions.CallView(ctx, c.WaspClient, c.ChainID.String(), viewCall)
}
