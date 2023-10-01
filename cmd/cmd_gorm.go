package cmd

import (
	"yafgo/yafgo-layout/internal/g"
	"yafgo/yafgo-layout/internal/gorm_gen"
	"yafgo/yafgo-layout/pkg/app"

	"github.com/gookit/color"
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

	app.App().AddSubCommand(app.SubCommand{
		Cmd:       subCmd,
		IsDefault: false,
	})
}

func runGormGen(cmd *cobra.Command, args []string) {

	color.Successln("Run gorm_gen...")
	conf := g.Cfg()
	dsn := conf.GetString("data.mysql.default")
	gorm_gen.RunGenerate(dsn)

}
