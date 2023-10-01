package cmd

import (
	"context"
	"time"
	"yafgo/yafgo-layout/internal/g"
	"yafgo/yafgo-layout/pkg/app"
	"yafgo/yafgo-layout/pkg/migration"
)

func RunApp() {
	_app := app.App()
	_app.PreRun = func() {
		ctx := context.Background()
		// 初始化配置
		// 由于大多逻辑都可能用到配置, 所以配置初始化应该首先被执行
		// app.setupConfig()
		// app.initTimeZone()

		// 初始化 logger
		// app.setupLogger(g.Cfg())

		// 初始化 cache
		g.SetupCache(ctx, g.Cfg())

		// 初始化 migration
		migration.Setup(g.Cfg().Viper)

		// 初始化飞书等通知
		g.SetupNotify()
	}
	app.App().Run()
}

// preRun 前置操作
func (app *Application) preRun() {
	ctx := context.Background()
	// 初始化配置
	// 由于大多逻辑都可能用到配置, 所以配置初始化应该首先被执行
	// app.setupConfig()
	app.initTimeZone()

	// 初始化 logger
	app.setupLogger(g.Cfg())

	// 初始化 cache
	g.SetupCache(ctx, g.Cfg())

	// 初始化 migration
	migration.Setup(g.Cfg().Viper)

	// 初始化飞书等通知
	g.SetupNotify()
}

/* // setupConfig 初始化配置
func (app *Application) setupConfig() {
	g.Cfg = ycfg.New(app.Mode,
		ycfg.WithType("yaml"),
		ycfg.WithEnvPrefix("YAFGO"),
		// ycfg.WithUnmarshalObj(g.Config),
	)
} */

func (app *Application) initTimeZone() {
	var cstZone = time.FixedZone("CST", 8*3600) // 东八区
	time.Local = cstZone
}
