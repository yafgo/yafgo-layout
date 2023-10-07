package app

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/spf13/cobra"
)

type bootstrap struct {
	app *application

	mode        string // dev, prod, ...
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
	bs.rootCmd = &cobra.Command{
		Use:   "yafgo",
		Short: "yafgo",
		Long:  `You can use "-h" flag to see all subcommands`,

		// rootCmd 的所有子命令都会执行以下前置代码
		PersistentPreRun: func(command *cobra.Command, args []string) {
			// [前置执行] 这里已经可以获取到命令行参数了
			app, _, err := newApp(bs.mode)
			if err != nil {
				fmt.Printf("failed to create app: %s\n", err)
				os.Exit(2)
			}
			bs.app = app
		},
	}

	var rootCmd = bs.rootCmd
	var cmdDefault *cobra.Command

	// 注册子命令
	bs.registerSubCommand()
	for _, v := range bs.subCommands {
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
		bs.registerDefaultCmd(cmdDefault)
	}

	// 注册全局选项
	bs.registerGlobalFlags()

	bs.initialized = true
	return bs
}

// registerGlobalFlags 注册全局选项（flag）
func (bs *bootstrap) registerGlobalFlags() {
	var rootCmd = bs.rootCmd
	rootCmd.PersistentFlags().StringVarP(&bs.mode, "conf", "c", "dev", "Set app config, eg: dev, prod...")
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
