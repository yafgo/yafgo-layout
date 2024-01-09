package play

import (
	"fmt"
	"yafgo/yafgo-layout/internal/g"
	"yafgo/yafgo-layout/internal/providers"
	"yafgo/yafgo-layout/internal/repository"
	"yafgo/yafgo-layout/internal/service"
	"yafgo/yafgo-layout/pkg/file/ossutil"
	"yafgo/yafgo-layout/pkg/notify"

	"go.uber.org/dig"
)

// NewApp 使用dig初始化应用
func NewApp() (app *Playground) {
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

	// 主应用
	container.Provide(NewPlayground)
	err := container.Invoke(func(_pg *Playground) {
		app = _pg
	})
	if err != nil {
		fmt.Printf("err: %+v", err)
	}

	return
}
