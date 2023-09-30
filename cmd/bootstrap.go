package cmd

import (
	"context"
	"time"
	"yafgo/yafgo-layout/internal/global"
	"yafgo/yafgo-layout/internal/query"
	"yafgo/yafgo-layout/pkg/database"
	"yafgo/yafgo-layout/pkg/logger"
	"yafgo/yafgo-layout/pkg/migration"
	"yafgo/yafgo-layout/pkg/sys/ycfg"
	"yafgo/yafgo-layout/pkg/sys/ylog"
)

// preRun 前置操作
func (app *Application) preRun() {
	ctx := context.Background()
	// 初始化配置
	// 由于大多逻辑都可能用到配置, 所以配置初始化应该首先被执行
	app.setupConfig()
	app.initTimeZone()

	// 初始化 logger
	app.setupLogger(global.Ycfg)

	// 初始化 cache
	global.SetupCache(ctx, global.Ycfg)

	// 初始化 gorm
	app.setupGorm(global.Ycfg)

	// 初始化 migration
	migration.Setup(global.Ycfg.Viper)

	// 初始化飞书等通知
	global.SetupNotify()
}

// setupConfig 初始化配置
func (app *Application) setupConfig() {
	global.Ycfg = ycfg.New(app.Mode,
		ycfg.WithType("yaml"),
		ycfg.WithEnvPrefix("YAFGO"),
		// ycfg.WithUnmarshalObj(global.Config),
	)
}

// setupGorm 初始化 gorm
func (app *Application) setupGorm(cfg *ycfg.Config) {
	gormLogger := logger.NewGormLogger(ylog.DefaultLogger())
	gormDB, err := database.NewGormMysql(cfg, gormLogger)
	if err != nil {
		panic(err)
	}
	// 赋值全局 mysql 对象
	global.Mysql = gormDB
	// 设置 query 使用的默认 db 对象
	query.SetDefault(gormDB)
}

func (app *Application) initTimeZone() {
	var cstZone = time.FixedZone("CST", 8*3600) // 东八区
	time.Local = cstZone
}
