package migration

import (
	"yafgo/yafgo-layout/pkg/sys/ycfg"

	"github.com/spf13/cobra"
)

// migrate 主命令
//
//	go run . migrate
//	go run . migrate migrate
//	go run . migrate rollback [rollback_step]
//	go run . migrate make <migration_name>
func NewMigrateCmd(cfg *ycfg.Config) *cobra.Command {
	m := NewMigrator(cfg)
	cmdMigrate := m.CmdMigrate()

	var cmdMigration = &cobra.Command{
		Use:   "migrate",
		Short: "Database migration management",
		Args:  cobra.NoArgs,
	}
	// cmdMigration.Run = cmdMigrate.Run
	cmdMigration.AddCommand(
		m.CmdMake(),            // migrate make <migration_name>
		cmdMigrate,             // migrate migrate
		m.CmdMigrateRollback(), // migrate rollback <step>
		m.CmdForceVersion(),    // migrate force <verion>
	)
	return cmdMigration
}
