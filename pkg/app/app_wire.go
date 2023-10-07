//go:build wireinject
// +build wireinject

package app

import (
	"yafgo/yafgo-layout/internal/server"

	"github.com/google/wire"
)

var yCfgSet = wire.NewSet(NewYCfg)

var yLogSet = wire.NewSet(NewYLog)

func newApp(envConf string) (app *Application, cleanUp func(), err error) {
	panic(wire.Build(
		newApplication,
		server.NewWebService,
		yCfgSet,
		yLogSet,
	))
}
