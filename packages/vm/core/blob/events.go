package blob

import (
	"github.com/nnikolash/wasp-types-exported/packages/hashing"
	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/util/rwutil"
)

func eventStore(ctx isc.Sandbox, blobHash hashing.HashValue) {
	ww := rwutil.NewBytesWriter()
	ww.Write(&blobHash)
	ctx.Event("coreblob.store", ww.Bytes())
}
