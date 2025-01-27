package vmimpl

import (
	"github.com/nnikolash/wasp-types-exported/packages/isc"
)

var _ isc.LogInterface = &requestContext{}

func (reqctx *requestContext) Infof(format string, params ...interface{}) {
	reqctx.vm.task.Log.Infof(format, params...)
}

func (reqctx *requestContext) Debugf(format string, params ...interface{}) {
	reqctx.vm.task.Log.Debugf(format, params...)
}

func (reqctx *requestContext) Panicf(format string, params ...interface{}) {
	reqctx.vm.task.Log.Panicf(format, params...)
}

func (reqctx *requestContext) Warnf(format string, params ...interface{}) {
	reqctx.vm.task.Log.Warnf(format, params...)
}
