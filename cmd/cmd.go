package cmd

import (
	"time"
	"yafgo/yafgo-layout/internal/g"
	"yafgo/yafgo-layout/pkg/app"
	"yafgo/yafgo-layout/pkg/migration"
)

func RunApp() {
	_app := app.App()
	_app.PreRun = func() {
		initTimeZone()

		// 初始化 migration
		migration.Setup(g.Cfg().Viper)

		// 初始化飞书等通知
		g.SetupNotify()
	}
	app.App().Run()
}

func initTimeZone() {
	var cstZone = time.FixedZone("CST", 8*3600) // 东八区
	time.Local = cstZone
}
