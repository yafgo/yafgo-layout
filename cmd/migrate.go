package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	// 用于执行数据迁移的命令
	var subCmd = subCommand{
		cmd: &cobra.Command{
			Use:   "migrate",
			Short: "Run Auto Migrate for database",
			Run:   runMigrate,
			Args:  cobra.NoArgs, // 不需要参数
		},
		isDefault: true,
	}
	addSubCommand(subCmd)
}

func runMigrate(cmd *cobra.Command, args []string) {

	fmt.Println("run migrate...")

}
