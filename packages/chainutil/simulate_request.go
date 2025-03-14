package chainutil

import (
	"fmt"
	"time"

	"github.com/nnikolash/wasp-types-exported/packages/chain"
	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/blocklog"
)

func SimulateRequest(
	ch chain.ChainCore,
	req isc.Request,
	estimateGas bool,
) (*blocklog.RequestReceipt, error) {
	aliasOutput, err := ch.LatestAliasOutput(chain.ActiveOrCommittedState)
	if err != nil {
		return nil, fmt.Errorf("could not get latest AliasOutput: %w", err)
	}
	res, err := runISCRequest(ch, aliasOutput, time.Now(), req, estimateGas)
	if err != nil {
		return nil, err
	}
	return res.Receipt, nil
}
