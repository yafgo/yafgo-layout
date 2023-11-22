package app

import (
	"yafgo/yafgo-layout/internal/gorm_gen"
	"yafgo/yafgo-layout/internal/make"
	"yafgo/yafgo-layout/pkg/migration"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

type subCommand struct {
	Cmd       *cobra.Command
	IsDefault bool
}

// addSubCommand 注册子命令
func (app *application) addSubCommand(subCmd ...subCommand) {
	app.subCommands = append(app.subCommands, subCmd...)
}

// registerSubCommand 注册子命令
func (app *application) registerSubCommand() {
	// web 服务
	app.addSubCommand(subCommand{
		Cmd: &cobra.Command{
			Use:   "serve",
			Short: "Run WebServer",
			Args:  cobra.NoArgs,
			Run: func(cmd *cobra.Command, args []string) {
				app.webService.CmdRun(cmd, args)
			},
		},
		IsDefault: true,
	})

	// 用于play演示的命令
	app.addSubCommand(subCommand{
		Cmd:       app.playground.PlayCommand(),
		IsDefault: false,
	})

	// gorm相关命令
	app.addSubCommand(subCommand{
		Cmd:       app.cmdGORM(),
		IsDefault: false,
	})

	// 用于执行代码生成的命令
	app.addSubCommand(subCommand{
		Cmd:       make.CmdMake,
		IsDefault: false,
	})

	// 用于执行数据迁移的命令
	app.addSubCommand(subCommand{
		Cmd:       migration.NewMigrateCmd(app.cfg),
		IsDefault: false,
	})
}

func (app *application) cmdGORM() *cobra.Command {
	// 用于执行数据迁移的命令
	var subCmd = &cobra.Command{
		Use:   "gorm",
		Short: "Gorm Tools, Generate Gorm models & queries",
		Args:  cobra.NoArgs,
		Run:   app.runGormGen,
	}
	subCmd.AddCommand(&cobra.Command{
		Use:   "gen",
		Short: "Generate Gorm models & queries",
		Args:  cobra.NoArgs,
		Run:   app.runGormGen,
	})

	return subCmd
}

func (app *application) runGormGen(cmd *cobra.Command, args []string) {
	color.Successln("Run gorm_gen...")
	dsn := app.cfg.GetString("data.mysql.default")
	gorm_gen.RunGenerate(dsn)
}
