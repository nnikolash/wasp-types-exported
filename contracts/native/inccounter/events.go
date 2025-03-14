package inccounter

import (
	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/util/rwutil"
)

func eventCounter(ctx isc.Sandbox, val int64) {
	ww := rwutil.NewBytesWriter()
	ww.WriteInt64(val)
	ctx.Event("inccounter.counter", ww.Bytes())
}
