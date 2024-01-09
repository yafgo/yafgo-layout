package yserve

import (
	"context"
	"time"
	"yafgo/yafgo-layout/internal/providers"
	"yafgo/yafgo-layout/internal/server"
	"yafgo/yafgo-layout/pkg/sys/ycfg"
	"yafgo/yafgo-layout/pkg/sys/ylog"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
)

type application struct {
	providers.BaseApplication

	cfg        *ycfg.Config
	logger     *ylog.Logger
	rdb        *redis.Client
	webService *server.WebService
}

func newApplication(
	logger *ylog.Logger,
	cfg *ycfg.Config,
	rdb *redis.Client,
	webService *server.WebService,
) (app *application) {
	return &application{
		cfg:        cfg,
		logger:     logger,
		rdb:        rdb,
		webService: webService,
	}
}

// Command 返回主命令
func (p *application) Command() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "serve",
		Short: "Run WebServer",
		// 前置
		PreRunE: func(cmd *cobra.Command, args []string) error {
			p.actionBefore()
			return nil
		},
		// 后置
		PostRunE: func(cmd *cobra.Command, args []string) error {
			time.Sleep(time.Second * 1)
			return nil
		},
		RunE: p.run,
	}
	p.RegisterGlobalFlags(rootCmd)
	return rootCmd
}

func (p *application) run(cmd *cobra.Command, args []string) error {
	p.webService.CmdRun(cmd, args)
	return nil
}

// actionBefore 启动前操作
func (p *application) actionBefore() {
	ctx := context.Background()

	// 启用redis配置存储
	p.cfg.EnableRedis(ctx, p.rdb, "yafgo")
}
