package cmd

import (
	"yafgo/yafgo-layout/pkg/app"
	"yafgo/yafgo-layout/pkg/migration"
)

func init() {
	// 用于执行数据迁移的命令
	var subCmd = app.SubCommand{
		Cmd:       migration.CmdMigration,
		IsDefault: false,
	}
	app.App().AddSubCommand(subCmd)
}
