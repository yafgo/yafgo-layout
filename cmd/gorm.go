package cmd

import (
	"fmt"
	"go-toy/toy-layout/internal/gorm_gen"
	"go-toy/toy-layout/pkg/config"

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
	conf := config.Config()
	dsn := conf.GetString("data.mysql.default")
	gorm_gen.RunGenerate(dsn)

}
