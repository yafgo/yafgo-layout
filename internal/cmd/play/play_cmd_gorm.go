package play

import (
	"context"
	"fmt"
	"time"
	"yafgo/yafgo-layout/internal/model"
	"yafgo/yafgo-layout/pkg/sys/ylog"

	"github.com/spf13/cobra"
)

func (p *Playground) playGorm() *cobra.Command {
	return &cobra.Command{
		Use:   "gorm",
		Short: "gorm测试",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			ylog.Debug(ctx, "GORM 测试")
			res := make([]map[string]any, 0)
			tx := p.db.Raw("show tables").Find(&res)
			if tx.Error != nil {
				ylog.Errorf(ctx, "err: %v", tx.Error)
				return
			}
			ylog.Debug(ctx, res)

			now := time.Now()
			user := &model.User{
				Name:     "张三",
				Phone:    fmt.Sprintf("1%d", now.Unix()),
				Username: fmt.Sprintf("yuser_%d", now.Unix()),
				Password: "123456",
			}
			userDO := p.q.User.WithContext(ctx)
			userDO.Create(user)
			ylog.Debug(ctx, user)

			user, _ = userDO.First()
			ylog.Debug(ctx, user)

			user, _ = userDO.First()
			ylog.Debug(ctx, user)
		},
	}
}
