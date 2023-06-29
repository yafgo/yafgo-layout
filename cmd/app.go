package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

type subCommand struct {
	cmd       *cobra.Command
	isDefault bool
}

var subCommands []subCommand = make([]subCommand, 0, 5)

// addSubCommand 添加二级命令
func addSubCommand(cmd subCommand) {
	subCommands = append(subCommands, cmd)
}

type Application struct {
	Mode string // dev, prod, ...

	initialized bool
	rootCmd     *cobra.Command
}

// Initialize 初始化
func (app *Application) Initialize() *Application {
	app.rootCmd = &cobra.Command{
		Use:   "toy",
		Short: "go-toy",
		Long:  `You can use "-h" flag to see all subcommands`,

		// rootCmd 的所有子命令都会执行以下前置代码
		PersistentPreRun: func(command *cobra.Command, args []string) {
			// [前置执行] 这里已经可以获取到命令行参数了
			app.preRun()
		},
	}

	var rootCmd = app.rootCmd
	var cmdDefault *cobra.Command

	// 注册子命令
	for _, v := range subCommands {
		if v.cmd != nil {
			rootCmd.AddCommand(v.cmd)
			if v.isDefault && cmdDefault == nil {
				cmdDefault = v.cmd
				rootCmd.Long = fmt.Sprintf(`Default will run "%s" command, you can use "-h" flag to see all subcommands`, cmdDefault.Use)
			}
		}
	}

	// 配置默认命令
	if cmdDefault != nil {
		app.registerDefaultCmd(cmdDefault)
	}

	// 注册全局选项
	app.registerGlobalFlags()

	app.initialized = true
	return app
}

// Run 执行cmd
func (app *Application) Run() {
	if !app.initialized {
		app.Initialize()
	}
	// 执行主命令
	if err := app.rootCmd.Execute(); err != nil {
		log.Fatalf("Failed to run app with %v: %+v\n", os.Args, err)
	}
}

// registerGlobalFlags 注册全局选项（flag）
func (app *Application) registerGlobalFlags() {
	var rootCmd = app.rootCmd
	rootCmd.PersistentFlags().StringVarP(&app.Mode, "mode", "m", "dev", "Set app mode, eg: dev, prod...")
}

// registerDefaultCmd 注册默认命令
func (app *Application) registerDefaultCmd(subCmd *cobra.Command) {
	firstArg := firstElement(os.Args[1:])
	if firstArg == "-h" || firstArg == "--help" {
		return
	}
	rootCmd := app.rootCmd
	cmd, _, err := rootCmd.Find(os.Args[1:])
	if err == nil && cmd.Use == rootCmd.Use {
		args := append([]string{subCmd.Use}, os.Args[1:]...)
		rootCmd.SetArgs(args)
	}
}

// FirstElement 安全地获取 args[0]，避免 panic: runtime error: index out of range
func firstElement(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return ""
}
