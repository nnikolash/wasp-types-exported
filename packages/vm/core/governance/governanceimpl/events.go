package governanceimpl

import (
	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/util/rwutil"
)

func eventRotate(ctx isc.Sandbox, newAddr iotago.Address, oldAddr iotago.Address) {
	ww := rwutil.NewBytesWriter()
	isc.AddressToWriter(ww, newAddr)
	isc.AddressToWriter(ww, oldAddr)
	ctx.Event("coregovernance.rotate", ww.Bytes())
}
