//go:build wireinject
// +build wireinject

package app

import (
	"yafgo/yafgo-layout/internal/g"
	"yafgo/yafgo-layout/internal/handler"
	"yafgo/yafgo-layout/internal/play"
	"yafgo/yafgo-layout/internal/repository"
	"yafgo/yafgo-layout/internal/server"
	"yafgo/yafgo-layout/internal/service"
	"yafgo/yafgo-layout/pkg/notify"

	"github.com/google/wire"
)

var playgroundSet = wire.NewSet(play.NewPlayground)

var handlerSet = wire.NewSet(
	handler.NewHandler,
)

var serviceSet = wire.NewSet(
	service.NewService,
	service.NewUserService,
)

var repositorySet = wire.NewSet(
	NewRedis,
	NewCache,
	NewDB,
	NewGormQuery,
	repository.NewRepository,
	repository.NewUserRepository,
)

var notifySet = wire.NewSet(
	notify.NewFeishu,
)

var gSet = wire.NewSet(g.New)

var jwtSet = wire.NewSet(NewJwt)

var yCfgSet = wire.NewSet(NewYCfg)

var yLogSet = wire.NewSet(NewYLog)

func newApp(envConf string) (app *application, err error) {
	panic(wire.Build(
		newApplication,
		playgroundSet,
		server.NewWebService,
		handlerSet,
		serviceSet,
		repositorySet,
		notifySet,
		gSet,
		jwtSet,
		yCfgSet,
		yLogSet,
	))
}
