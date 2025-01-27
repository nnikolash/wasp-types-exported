package sbtestsc

import (
	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/util/rwutil"
)

func eventCounter(ctx isc.Sandbox, value uint64) {
	ww := rwutil.NewBytesWriter()
	ww.WriteUint64(value)
	ctx.Event("testcore.counter", ww.Bytes())
}

func eventTest(ctx isc.Sandbox) {
	ww := rwutil.NewBytesWriter()
	ctx.Event("testcore.test", ww.Bytes())
}
