package play

import (
	"fmt"
	"yafgo/yafgo-layout/internal/query"
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
	db     func() *gorm.DB
	rdb    func() *redis.Client
	q      func() *query.Query
	logger func() *ylog.Logger
}

func NewPlayground(
	db func() *gorm.DB,
	rdb func() *redis.Client,
	q func() *query.Query,
	logger func() *ylog.Logger,
) *Playground {
	pg := &Playground{
		db:     db,
		rdb:    rdb,
		q:      q,
		logger: logger,
	}
	pg.subCmds = make([]*cobra.Command, 10)
	return pg
}

// PlayCommand 返回 play 主命令
func (p *Playground) PlayCommand() *cobra.Command {
	playCmd := &cobra.Command{
		Use:   "play",
		Short: "A playground for testing",
		Run:   p.runPlay,
	}
	p.addSubCommands()
	playCmd.AddCommand(p.subCmds...)
	return playCmd
}

func (p *Playground) runPlay(cmd *cobra.Command, args []string) {

	fmt.Println("\n[play test case...]")

	// 可以临时在这里调用一些代码
}
