package sbtestsc

import (
	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/kv/dict"
)

// spawn deploys new contract and calls it
func spawn(ctx isc.Sandbox) dict.Dict {
	ctx.Log().Debugf(FuncSpawn.Name)
	progHash := ctx.Params().MustGetHashValue(ParamProgHash)
	name := Contract.Name + "_spawned"
	hname := isc.Hn(name)
	ctx.DeployContract(progHash, name, nil)

	for i := 0; i < 5; i++ {
		ctx.Call(hname, FuncIncCounter.Hname(), nil, nil)
	}
	ctx.Log().Debugf("sbtestsc.spawn: new contract name = %s hname = %s", name, hname.String())
	return nil
}
