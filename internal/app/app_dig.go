package app

import (
	"fmt"
	"yafgo/yafgo-layout/internal/g"
	"yafgo/yafgo-layout/internal/handler"
	"yafgo/yafgo-layout/internal/play"
	"yafgo/yafgo-layout/internal/repository"
	"yafgo/yafgo-layout/internal/server"
	"yafgo/yafgo-layout/internal/service"
	"yafgo/yafgo-layout/pkg/notify"

	"go.uber.org/dig"
)

func NewApp() (app *application) {
	envConf := parseEnvConfig()
	app = NewAppDig(envConf)
	// app, _ = newApp(envConf)
	return
}

// NewAppDig 使用dig初始化应用
func NewAppDig(envConf string) (app *application) {
	container := dig.New()

	// 基础依赖
	container.Provide(NewYCfgDig(envConf))
	container.Provide(NewYLog)

	// 杂项
	container.Provide(NewJwt)
	container.Provide(g.New)
	container.Provide(notify.NewFeishu)

	// db等依赖
	container.Provide(NewRedis)
	container.Provide(NewCache)
	container.Provide(NewDB)
	container.Provide(NewGormQuery)

	// handler
	container.Provide(handler.NewHandler)

	// service
	container.Provide(service.NewService)
	container.Provide(service.NewUserService)

	// repository
	container.Provide(repository.NewRepository)
	container.Provide(repository.NewUserRepository)

	// playground
	container.Provide(play.NewPlayground)
	// web服务
	container.Provide(server.NewWebService)

	// 主应用
	container.Provide(newApplication)
	err := container.Invoke(func(_app *application) {
		app = _app
	})
	if err != nil {
		fmt.Printf("err: %+v", err)
	}

	return
}
