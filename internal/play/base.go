package play

import (
	"fmt"
	"time"
	"yafgo/yafgo-layout/internal/query"
	"yafgo/yafgo-layout/pkg/notify"
	"yafgo/yafgo-layout/pkg/sys/ylog"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

// Playground 演练场
type Playground struct {
	// 需要注册到 play 命令的子命令
	subCmds []*cobra.Command

	// Playground 需要用到的外部组件
	db     *gorm.DB
	rdb    *redis.Client
	q      *query.Query
	logger *ylog.Logger
	feishu *notify.FeishuRobot
}

func NewPlayground(
	db *gorm.DB,
	rdb *redis.Client,
	q *query.Query,
	logger *ylog.Logger,
	feishu *notify.FeishuRobot,
) *Playground {
	pg := &Playground{
		db:     db,
		rdb:    rdb,
		q:      q,
		logger: logger,
		feishu: feishu,
	}
	pg.subCmds = make([]*cobra.Command, 0, 10)
	return pg
}

// PlayCommand 返回 play 主命令
func (p *Playground) PlayCommand() *cobra.Command {
	playCmd := &cobra.Command{
		Use:   "play",
		Short: "A playground for testing",
		Long:  `You can use "-h" flag to see all subcommands`,
		Run:   p.runPlay,
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			time.Sleep(time.Second * 1)
		},
	}
	p.addSubCommands()
	playCmd.AddCommand(p.subCmds...)
	return playCmd
}

func (p *Playground) runPlay(cmd *cobra.Command, args []string) {

	fmt.Println("\n[play test case...]")

	// 可以临时在这里调用一些代码
}
