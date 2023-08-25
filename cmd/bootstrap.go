package cmd

import (
	"go-toy/toy-layout/global"
	"go-toy/toy-layout/internal/query"
	"go-toy/toy-layout/pkg/database"
	"go-toy/toy-layout/pkg/logger"
	"go-toy/toy-layout/pkg/migration"
	"go-toy/toy-layout/pkg/sys/ycfg"
	"go-toy/toy-layout/pkg/sys/ylog"
	"time"

	"github.com/spf13/viper"
)

// preRun 前置操作
func (app *Application) preRun() {
	// 初始化配置
	// 由于大多逻辑都可能用到配置, 所以配置初始化应该首先被执行
	app.setupConfig()
	app.initTimeZone()

	// 初始化 logger
	app.setupLogger(global.Ycfg.Viper)

	// 初始化 gorm
	app.setupGorm(global.Ycfg.Viper)

	// 初始化 migration
	migration.Setup(global.Ycfg.Viper)
}

// setupConfig 初始化配置
func (app *Application) setupConfig() {
	global.Ycfg = ycfg.New(app.Mode,
		ycfg.WithType("yaml"),
		ycfg.WithEnvPrefix("YAFGO"),
		// ycfg.WithUnmarshalObj(global.Config),
	)
}

// setupLogger 初始化 logger
func (app *Application) setupLogger(cfg *viper.Viper) {
	lgCfg := ylog.DefaultConfig()
	lgCfg.CtxKeys = []string{"req_id"}
	lgCfg.Prefix = global.AppName()
	// lgCfg.LogType = ylog.LogTypeJson
	// lgCfg.Stdout = true
	lgCfg.Rotate = ylog.ConfigRotate{
		MaxSize:    4,
		MaxAge:     7,
		MaxBackups: 0,
		LocalTime:  true,
		Compress:   false,
	}
	lg := ylog.New(lgCfg)
	ylog.SetDefaultLogger(lg)
}

// setupGorm 初始化 gorm
func (app *Application) setupGorm(cfg *viper.Viper) {
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
