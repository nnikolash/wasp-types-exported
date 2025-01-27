package logger

import (
	"github.com/iotaledger/hive.go/app"
	"github.com/nnikolash/wasp-types-exported/packages/evm/evmlogger"
)

func init() {
	Component = &app.Component{
		Name:      "Logger",
		Configure: configure,
	}
}

var Component *app.Component

func configure() error {
	evmlogger.Init(Component.App().NewLogger("go-ethereum"))
	return nil
}
