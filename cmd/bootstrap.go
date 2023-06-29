package cmd

import "go-toy/toy-layout/pkg/config"

// preRun 前置操作
func (app *Application) preRun() {
	// 初始化配置
	app.setupConfig()
}

// setupConfig 初始化配置
func (app *Application) setupConfig() {
	config.SetConfigDir("etc")
	config.SetEnvPrefix("GOTOY")
	config.SetupConfig(app.Mode)
}
