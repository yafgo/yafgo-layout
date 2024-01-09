package main

import (
	"fmt"
	"yafgo/yafgo-layout/internal/providers"
	"yafgo/yafgo-layout/pkg/migration"

	"github.com/spf13/cobra"
	"go.uber.org/dig"
)

func main() {
	app := newApp()
	app.Execute()
}

func newApp() (app *cobra.Command) {
	container := dig.New()

	// 基础依赖
	container.Provide(providers.ParseConfigName)
	container.Provide(providers.NewYCfg)

	// 主应用
	container.Provide(migration.NewMigrateCmd)
	err := container.Invoke(func(_app *cobra.Command) {
		app = _app
	})
	if err != nil {
		fmt.Printf("err: %+v", err)
	}

	return
}
