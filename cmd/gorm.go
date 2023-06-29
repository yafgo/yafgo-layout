package cmd

import (
	"fmt"
	"go-toy/toy-layout/internal/gorm_gen"

	"github.com/spf13/cobra"
)

func init() {
	// 用于执行数据迁移的命令
	var subCmd = &cobra.Command{
		Use:   "gorm",
		Short: "Gorm Tools, Generate Gorm models & queries",
		Args:  cobra.NoArgs,
		Run:   runGormGen,
	}
	subCmd.AddCommand(&cobra.Command{
		Use:   "gen",
		Short: "Generate Gorm models & queries",
		Args:  cobra.NoArgs,
		Run:   runGormGen,
	})

	addSubCommand(subCommand{
		cmd:       subCmd,
		isDefault: false,
	})
}

func runGormGen(cmd *cobra.Command, args []string) {

	fmt.Println("run gorm_gen...")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		"dbUser", "dbPassword", "dbHost", "dbPort", "db")
	gorm_gen.RunGenerate(dsn)

}
