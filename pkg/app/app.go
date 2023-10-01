package app

// App 全局app实例
func App() *application {
	appOnce.Do(func() {
		appInst = new(application)
		appInst.subCommands = make([]SubCommand, 0, 10)
	})
	return appInst
}

// Mode app模式: dev, prod, ...
func (app *application) Mode() string {
	return app.mode
}
