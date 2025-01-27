package chainclient

import (
	"context"

	"github.com/nnikolash/wasp-types-exported/clients/apiclient"
	"github.com/nnikolash/wasp-types-exported/packages/hashing"
	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/kv/dict"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/blob"
)

// UploadBlob sends an off-ledger request to call 'store' in the blob contract.
func (c *Client) UploadBlob(ctx context.Context, fields dict.Dict) (hashing.HashValue, isc.OffLedgerRequest, *apiclient.ReceiptResponse, error) {
	blobHash := blob.MustGetBlobHash(fields)

	req, err := c.PostOffLedgerRequest(ctx,
		blob.Contract.Hname(),
		blob.FuncStoreBlob.Hname(),

		PostRequestParams{
			Args: fields,
		},
	)
	if err != nil {
		return hashing.NilHash, nil, nil, err
	}

	receipt, _, err := c.WaspClient.ChainsApi.WaitForRequest(ctx, c.ChainID.String(), req.ID().String()).Execute()

	return blobHash, req, receipt, err
}
