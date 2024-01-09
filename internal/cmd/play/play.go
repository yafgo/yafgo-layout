package play

import (
	"context"
	"fmt"
	"time"
	"yafgo/yafgo-layout/internal/providers"
	"yafgo/yafgo-layout/internal/query"
	"yafgo/yafgo-layout/pkg/file/ossutil"
	"yafgo/yafgo-layout/pkg/notify"
	"yafgo/yafgo-layout/pkg/sys/ycfg"
	"yafgo/yafgo-layout/pkg/sys/ylog"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

// Playground 演练场
type Playground struct {
	providers.BaseApplication

	// 需要注册到 play 命令的子命令
	subCmds []*cobra.Command

	// Playground 需要用到的外部组件
	db     *gorm.DB
	rdb    *redis.Client
	q      *query.Query
	logger *ylog.Logger
	cfg    *ycfg.Config
	feishu *notify.FeishuRobot
	Oss    *ossutil.OssUtil
}

func NewPlayground(
	db *gorm.DB,
	rdb *redis.Client,
	q *query.Query,
	logger *ylog.Logger,
	cfg *ycfg.Config,
	feishu *notify.FeishuRobot,
	oss *ossutil.OssUtil,
) *Playground {
	pg := &Playground{
		db:     db,
		rdb:    rdb,
		q:      q,
		logger: logger,
		cfg:    cfg,
		feishu: feishu,
		Oss:    oss,
	}
	pg.subCmds = make([]*cobra.Command, 0, 10)
	return pg
}

// Command 返回主命令
func (p *Playground) Command() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "play",
		Short: "A playground for testing",
		Long:  `You can use "-h" flag to see all subcommands`,
		Run:   p.runPlay,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			p.actionBefore(cmd.Context())
			return nil
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			time.Sleep(time.Second * 1)
		},
	}
	p.RegisterGlobalFlags(rootCmd)
	rootCmd.AddCommand(p.addSubCommands()...)
	return rootCmd
}

func (p *Playground) runPlay(cmd *cobra.Command, args []string) {

	fmt.Println("\n[play test case...]")

	// 可以临时在这里调用一些代码
}

// actionBefore 启动前操作
func (app *Playground) actionBefore(ctx context.Context) {

	// 启用redis配置存储
	app.cfg.EnableRedis(ctx, app.rdb, "yafgo")
}
