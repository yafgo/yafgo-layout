package cmd

import (
	"go-toy/toy-layout/pkg/config"
	"go-toy/toy-layout/pkg/logger"
)

// preRun 前置操作
func (app *Application) preRun() {
	// 初始化配置
	// 由于大多逻辑都可能用到配置, 所以配置初始化应该首先被执行
	app.setupConfig()

	// 初始化 logger
	app.setupLogger()
}

// setupConfig 初始化配置
func (app *Application) setupConfig() {
	config.SetConfigDir("etc")
	config.SetEnvPrefix("GOTOY")
	config.SetupConfig(app.Mode)
}

// setupLogger 初始化 logger
func (app *Application) setupLogger() {
	conf := config.Config()
	logger.SetIsProd(false)
	logger.SetPrefix("gotoy")
	logger.SetupLogger(conf)
}
