package yorm

import (
	"time"
	"yafgo/yafgo-layout/internal/database"
	"yafgo/yafgo-layout/internal/providers"
	"yafgo/yafgo-layout/pkg/sys/ycfg"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

// application 对gorm_gen的操作
type application struct {
	providers.BaseApplication

	cfg *ycfg.Config
}

func newApplication(
	cfg *ycfg.Config,
) *application {
	inst := &application{
		cfg: cfg,
	}
	return inst
}

// Command 返回主命令
func (p *application) Command() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "orm",
		Short: "Gorm Tools, Generate Gorm models & queries",
		Long:  `You can use "-h" flag to see all subcommands`,
		RunE:  p.run,
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			time.Sleep(time.Second * 1)
		},
	}
	rootCmd.AddCommand(&cobra.Command{
		Use:   "gen",
		Short: "Generate Gorm models & queries",
		Args:  cobra.NoArgs,
		RunE:  p.run,
	})
	p.RegisterGlobalFlags(rootCmd)
	return rootCmd
}

func (p *application) run(cmd *cobra.Command, args []string) error {
	color.Successln("Run gorm_gen...")
	dsn := p.cfg.GetString("data.mysql.default")
	database.RunGenerate(dsn)
	return nil
}
