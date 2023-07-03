package migration

import (
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

// CmdMigrate 命令: 执行迁移
func (p *migrator) CmdMake() *cobra.Command {
	return &cobra.Command{
		Use:   "make",
		Short: "Create migration files",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			_ = p.MakeMigration(args[0])
		},
	}
}

// CmdMigrate 命令: 执行迁移
func (p *migrator) CmdMigrate() *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "Run migrate",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			_ = p.RunMigrate()
		},
	}
}

// CmdMigrateRollback 命令: 回滚迁移
func (p *migrator) CmdMigrateRollback() *cobra.Command {
	return &cobra.Command{
		Use:   "rollback",
		Short: "Run migrate rollback",
		Args:  cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			var step int
			if len(args) > 0 {
				step = cast.ToInt(args[0])
			}
			if step <= 0 {
				step = 1
			}
			_ = p.RunRollback(step)
		},
	}
}

// CmdForceVersion 命令: 强制版本(用于迁移出错并修复后继续)
func (p *migrator) CmdForceVersion() *cobra.Command {
	return &cobra.Command{
		Use:   "force",
		Short: "Force migrate version",
		Args:  cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			var ver int
			if len(args) > 0 {
				ver = cast.ToInt(args[0])
			}
			_ = p.ForceVersion(ver)
		},
	}
}
