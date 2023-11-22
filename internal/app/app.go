package app

import (
	"log"
	"os"
	"sync"
	"yafgo/yafgo-layout/internal/play"
	"yafgo/yafgo-layout/internal/server"
	"yafgo/yafgo-layout/pkg/sys/ycfg"
	"yafgo/yafgo-layout/pkg/sys/ylog"

	"github.com/spf13/cobra"
)

type application struct {
	cfg        *ycfg.Config
	logger     *ylog.Logger
	webService *server.WebService
	playground *play.Playground

	// 非依赖属性
	envConf     string // dev, prod, ...
	mutex       sync.Mutex
	initialized bool
	rootCmd     *cobra.Command
	subCommands []subCommand
}

func newApplication(
	logger *ylog.Logger,
	cfg *ycfg.Config,
	webService *server.WebService,
	playground *play.Playground,
) (app *application) {
	return &application{
		cfg:        cfg,
		logger:     logger,
		webService: webService,
		playground: playground,
	}
}

// Run 执行cmd
func (app *application) Run() {
	app.mutex.Lock()
	if app.initialized {
		log.Fatalln("The App is already running...")
		app.mutex.Unlock()
		return
	}
	app.initialize()
	app.mutex.Unlock()

	// 执行主命令
	if err := app.rootCmd.Execute(); err != nil {
		log.Fatalf("Failed to run app with %v: %+v\n", os.Args, err)
	}
}

// initialize 初始化
func (app *application) initialize() {
	// 解析命令参数
	app.rootCmd = &cobra.Command{
		Use:   "yafgo",
		Short: "yafgo",
		Long:  `You can use "-h" flag to see all subcommands`,

		// rootCmd 的所有子命令都会执行以下前置代码
		PersistentPreRun: func(command *cobra.Command, args []string) {
			// [前置执行] 这里已经可以获取到命令行参数了
		},
	}

	// 注册子命令
	app.registerSubCommand()
	var cmdDefault *cobra.Command
	for _, v := range app.subCommands {
		if v.Cmd != nil {
			app.rootCmd.AddCommand(v.Cmd)
			if v.IsDefault && cmdDefault == nil {
				cmdDefault = v.Cmd
				app.rootCmd.Long = `Default will run "` + cmdDefault.Use + `" command, you can use "-h" flag to see all subcommands`
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

// registerGlobalFlags 注册全局选项（flag）
func (app *application) registerGlobalFlags() {
	var rootCmd = app.rootCmd
	rootCmd.PersistentFlags().StringVarP(&app.envConf, "conf", "c", "dev", "Set app config, eg: dev, prod...")
}

// parseEnvConfig 解析配置名称
func parseEnvConfig() (envConf string) {

	// 仅负责解析配置参数的 cobra 实例
	var preCmd = &cobra.Command{}

	// 禁用 -h 标志的响应
	preCmd.SetHelpFunc(func(c *cobra.Command, s []string) {})

	// 解析配置参数
	preCmd.PersistentFlags().StringVarP(&envConf, "conf", "c", "dev", "ConfigName: dev, prod...")

	// 执行主命令(这里忽略Execute的错误)
	preCmd.Execute()

	return
}

// FirstElement 安全地获取 args[0]，避免 panic: runtime error: index out of range
func firstElement(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return ""
}
