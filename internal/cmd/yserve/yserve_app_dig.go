package yserve

import (
	"fmt"
	"yafgo/yafgo-layout/internal/g"
	"yafgo/yafgo-layout/internal/handler"
	"yafgo/yafgo-layout/internal/handler/backend"
	"yafgo/yafgo-layout/internal/handler/frontend"
	"yafgo/yafgo-layout/internal/providers"
	"yafgo/yafgo-layout/internal/repository"
	"yafgo/yafgo-layout/internal/server"
	"yafgo/yafgo-layout/internal/service"
	"yafgo/yafgo-layout/pkg/file/ossutil"
	"yafgo/yafgo-layout/pkg/notify"

	"go.uber.org/dig"
)

// NewApp 使用dig初始化应用
func NewApp() (app *application) {
	container := dig.New()

	// 基础依赖
	container.Provide(providers.ParseConfigName)
	container.Provide(providers.NewYCfg)
	container.Provide(providers.NewYLog)

	// 杂项
	container.Provide(providers.NewJwt)
	container.Provide(g.New)
	container.Provide(notify.NewFeishu)
	container.Provide(ossutil.NewOssUtil)

	// db等依赖
	container.Provide(providers.NewRedis)
	container.Provide(providers.NewCache)
	container.Provide(providers.NewDB)
	container.Provide(providers.NewGormQuery)

	// service
	container.Provide(service.NewService)
	container.Provide(service.NewUserService)

	// repository
	container.Provide(repository.NewRepository)
	container.Provide(repository.NewUserRepository)

	// handler
	container.Provide(handler.NewHandler)

	// 分模块router
	container.Provide(backend.NewRouter)
	container.Provide(frontend.NewRouter)

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
