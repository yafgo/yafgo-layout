package app

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/spf13/cobra"
)

type SubCommand struct {
	Cmd       *cobra.Command
	IsDefault bool
}

type application struct {
	mode string // dev, prod, ...

	initialized bool
	rootCmd     *cobra.Command
	subCommands []SubCommand

	// 前置执行
	PreRun func()
}

var (
	appInst *application
	appOnce sync.Once
	appMu   sync.Mutex
)

// Run 执行cmd
func (app *application) Run() {
	appMu.Lock()
	if app.initialized {
		log.Fatalln("The App is already running...")
		appMu.Unlock()
		return
	}
	app.initialize()
	appMu.Unlock()

	// 执行主命令
	if err := app.rootCmd.Execute(); err != nil {
		log.Fatalf("Failed to run app with %v: %+v\n", os.Args, err)
	}
}

// AddSubCommand 注册子命令
func (app *application) AddSubCommand(subCmd ...SubCommand) {
	app.subCommands = append(app.subCommands, subCmd...)
}

// initialize 初始化
func (app *application) initialize() *application {
	app.rootCmd = &cobra.Command{
		Use:   "yafgo",
		Short: "yafgo",
		Long:  `You can use "-h" flag to see all subcommands`,

		// rootCmd 的所有子命令都会执行以下前置代码
		PersistentPreRun: func(command *cobra.Command, args []string) {
			// [前置执行] 这里已经可以获取到命令行参数了
			if app.PreRun != nil {
				app.PreRun()
			}
		},
	}

	var rootCmd = app.rootCmd
	var cmdDefault *cobra.Command

	// 注册子命令
	for _, v := range app.subCommands {
		if v.Cmd != nil {
			rootCmd.AddCommand(v.Cmd)
			if v.IsDefault && cmdDefault == nil {
				cmdDefault = v.Cmd
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

// registerGlobalFlags 注册全局选项（flag）
func (app *application) registerGlobalFlags() {
	var rootCmd = app.rootCmd
	rootCmd.PersistentFlags().StringVarP(&app.mode, "mode", "m", "dev", "Set app mode, eg: dev, prod...")
}

// registerDefaultCmd 注册默认命令
func (app *application) registerDefaultCmd(subCmd *cobra.Command) {
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
