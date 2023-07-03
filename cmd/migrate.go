package cmd

import (
	"go-toy/toy-layout/pkg/migration"
)

func init() {
	// 用于执行数据迁移的命令
	var subCmd = subCommand{
		cmd:       migration.CmdMigration,
		isDefault: false,
	}
	addSubCommand(subCmd)
}
