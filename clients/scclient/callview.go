package scclient

import (
	"context"

	"github.com/nnikolash/wasp-types-exported/packages/kv/dict"
)

func (c *SCClient) CallView(context context.Context, functionName string, args dict.Dict) (dict.Dict, error) {
	return c.ChainClient.CallView(context, c.ContractHname, functionName, args)
}
