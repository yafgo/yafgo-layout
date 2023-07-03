package migration

import (
	"github.com/spf13/cobra"
)

// migrate 主命令
//
//	go run . migrate
//	go run . migrate migrate
//	go run . migrate rollback [rollback_step]
//	go run . migrate make <migration_name>
var CmdMigration = &cobra.Command{
	Use:   "migrate",
	Short: "Database migration management",
	Args:  cobra.NoArgs,
}

func init() {
	m := NewMigrator()
	cmdMigrate := m.CmdMigrate()
	CmdMigration.Run = cmdMigrate.Run
	CmdMigration.AddCommand(
		m.CmdMake(),            // migrate make <migration_name>
		cmdMigrate,             // migrate migrate
		m.CmdMigrateRollback(), // migrate rollback <step>
		m.CmdForceVersion(),    // migrate force <verion>
	)
}
