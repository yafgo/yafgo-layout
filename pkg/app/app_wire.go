//go:build wireinject
// +build wireinject

package app

import (
	"yafgo/yafgo-layout/internal/handler"
	"yafgo/yafgo-layout/internal/repository"
	"yafgo/yafgo-layout/internal/server"
	"yafgo/yafgo-layout/internal/service"

	"github.com/google/wire"
)

var handlerSet = wire.NewSet(
	handler.NewHandler,
	handler.NewUserHandler,
)

var serviceSet = wire.NewSet(
	service.NewService,
	service.NewUserService,
)

var repositorySet = wire.NewSet(
	newRedis,
	newCache,
	newDB,
	newGormQuery,
	repository.NewRepository,
	repository.NewUserRepository,
)

var jwtSet = wire.NewSet(newJwt)

var yCfgSet = wire.NewSet(NewYCfg)

var yLogSet = wire.NewSet(NewYLog)

func newApp(envConf string) (app *application, cleanUp func(), err error) {
	panic(wire.Build(
		newApplication,
		server.NewWebService,
		handlerSet,
		serviceSet,
		repositorySet,
		jwtSet,
		yCfgSet,
		yLogSet,
	))
}
