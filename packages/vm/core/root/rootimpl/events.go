package rootimpl

import (
	"github.com/nnikolash/wasp-types-exported/packages/hashing"
	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/util/rwutil"
)

func eventDeploy(ctx isc.Sandbox, progHash hashing.HashValue, name string) {
	ww := rwutil.NewBytesWriter()
	ww.Write(&progHash)
	ww.WriteString(name)
	ctx.Event("coreroot.deploy", ww.Bytes())
}

func eventGrant(ctx isc.Sandbox, deployer isc.AgentID) {
	ww := rwutil.NewBytesWriter()
	ww.Write(deployer)
	ctx.Event("coreroot.grant", ww.Bytes())
}

func eventRevoke(ctx isc.Sandbox, deployer isc.AgentID) {
	ww := rwutil.NewBytesWriter()
	ww.Write(deployer)
	ctx.Event("coreroot.revoke", ww.Bytes())
}
