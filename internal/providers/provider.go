package providers

import "github.com/spf13/cobra"

type BaseApplication struct {
	ConfName string
}

// registerGlobalFlags 注册全局选项（flag）
func (app *BaseApplication) RegisterGlobalFlags(rootCmd *cobra.Command) {
	rootCmd.PersistentFlags().StringVarP(&app.ConfName, "conf", "c", "dev", "Set app config, eg: dev, prod...")
}
