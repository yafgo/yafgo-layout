package yorm

import (
	"fmt"
	"yafgo/yafgo-layout/internal/providers"

	"go.uber.org/dig"
)

// NewApp 使用dig初始化应用
func NewApp() (app *application) {
	container := dig.New()

	// 基础依赖
	container.Provide(providers.ParseConfigName)
	container.Provide(providers.NewYCfg)

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
