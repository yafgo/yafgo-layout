package app

import (
	"log"
	"os"
	"sync"

	"github.com/spf13/cobra"
)

type bootstrap struct {
	app *application

	envConf     string // dev, prod, ...
	mutex       sync.Mutex
	initialized bool
	rootCmd     *cobra.Command
	subCommands []subCommand
}

func Bootstrap() *bootstrap {
	return new(bootstrap)
}

// Run 执行cmd
func (bs *bootstrap) Run() {
	bs.mutex.Lock()
	if bs.initialized {
		log.Fatalln("The App is already running...")
		bs.mutex.Unlock()
		return
	}
	bs.initialize()
	bs.mutex.Unlock()

	// 执行主命令
	if err := bs.rootCmd.Execute(); err != nil {
		log.Fatalf("Failed to run app with %v: %+v\n", os.Args, err)
	}
}

// initialize 初始化
func (bs *bootstrap) initialize() *bootstrap {
	// 前置解析配置名
	bs.parseConfigFlag()
	app, err := newApp(bs.envConf)
	if err != nil {
		log.Printf("failed to create app: %s\n", err)
		os.Exit(2)
	}
	bs.app = app

	// 开始正式解析命令参数
	bs.rootCmd = &cobra.Command{
		Use:   "yafgo",
		Short: "yafgo",
		Long:  `You can use "-h" flag to see all subcommands`,

		// rootCmd 的所有子命令都会执行以下前置代码
		PersistentPreRun: func(command *cobra.Command, args []string) {
			// [前置执行] 这里已经可以获取到命令行参数了
		},
	}

	// 注册子命令
	bs.registerSubCommand()
	var cmdDefault *cobra.Command
	for _, v := range bs.subCommands {
		if v.Cmd != nil {
			bs.rootCmd.AddCommand(v.Cmd)
			if v.IsDefault && cmdDefault == nil {
				cmdDefault = v.Cmd
				bs.rootCmd.Long = `Default will run "` + cmdDefault.Use + `" command, you can use "-h" flag to see all subcommands`
			}
		}
	}

	// 配置默认命令
	if cmdDefault != nil {
		bs.registerDefaultCmd(cmdDefault)
	}

	// 注册全局选项
	bs.registerGlobalFlags()

	bs.initialized = true
	return bs
}

// parseConfigFlag 前置解析 configName
func (bs *bootstrap) parseConfigFlag() {
	// 正式 rootCmd 之前的仅负责解析配置参数的 cobra 实例
	var preCmd = &cobra.Command{}

	// 禁用 -h 标志的响应
	preCmd.SetHelpFunc(func(c *cobra.Command, s []string) {})

	// 解析配置参数
	preCmd.PersistentFlags().StringVarP(&bs.envConf, "conf", "c", "dev", "ConfigName: dev, prod...")

	// 执行主命令(这里忽略Execute的错误)
	preCmd.Execute()
}

// registerGlobalFlags 注册全局选项（flag）
func (bs *bootstrap) registerGlobalFlags() {
	var rootCmd = bs.rootCmd
	rootCmd.PersistentFlags().StringVarP(&bs.envConf, "conf", "c", "dev", "Set app config, eg: dev, prod...")
}

// registerDefaultCmd 注册默认命令
func (bs *bootstrap) registerDefaultCmd(subCmd *cobra.Command) {
	firstArg := firstElement(os.Args[1:])
	if firstArg == "-h" || firstArg == "--help" {
		return
	}
	rootCmd := bs.rootCmd
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
