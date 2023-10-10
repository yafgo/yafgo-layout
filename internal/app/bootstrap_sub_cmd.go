package app

import (
	"yafgo/yafgo-layout/internal/gorm_gen"
	"yafgo/yafgo-layout/pkg/migration"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

type subCommand struct {
	Cmd       *cobra.Command
	IsDefault bool
}

// addSubCommand 注册子命令
func (bs *bootstrap) addSubCommand(subCmd ...subCommand) {
	bs.subCommands = append(bs.subCommands, subCmd...)
}

// registerSubCommand 注册子命令
func (bs *bootstrap) registerSubCommand() {
	// web 服务
	bs.addSubCommand(subCommand{
		Cmd: &cobra.Command{
			Use:   "serve",
			Short: "Run WebServer",
			Args:  cobra.NoArgs,
			Run: func(cmd *cobra.Command, args []string) {
				bs.app.webService.CmdRun(cmd, args)
			},
		},
		IsDefault: true,
	})

	// 用于play演示的命令
	bs.addSubCommand(subCommand{
		Cmd:       bs.app.playground.PlayCommand(),
		IsDefault: false,
	})

	// gorm相关命令
	bs.addSubCommand(subCommand{
		Cmd:       bs.cmdGORM(),
		IsDefault: false,
	})

	// 用于执行数据迁移的命令
	bs.addSubCommand(subCommand{
		Cmd:       migration.NewMigrateCmd(bs.app.cfg),
		IsDefault: false,
	})
}

func (bs *bootstrap) cmdGORM() *cobra.Command {
	// 用于执行数据迁移的命令
	var subCmd = &cobra.Command{
		Use:   "gorm",
		Short: "Gorm Tools, Generate Gorm models & queries",
		Args:  cobra.NoArgs,
		Run:   bs.runGormGen,
	}
	subCmd.AddCommand(&cobra.Command{
		Use:   "gen",
		Short: "Generate Gorm models & queries",
		Args:  cobra.NoArgs,
		Run:   bs.runGormGen,
	})

	return subCmd
}

func (bs *bootstrap) runGormGen(cmd *cobra.Command, args []string) {
	color.Successln("Run gorm_gen...")
	dsn := bs.app.cfg.GetString("data.mysql.default")
	gorm_gen.RunGenerate(dsn)
}
