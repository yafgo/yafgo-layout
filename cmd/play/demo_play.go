package play

import (
	"context"
	"fmt"
	"time"
	"yafgo/yafgo-layout/internal/g"
	"yafgo/yafgo-layout/internal/model"
	"yafgo/yafgo-layout/internal/query"
	"yafgo/yafgo-layout/pkg/sys/ylog"

	"github.com/spf13/cobra"
)

// 需要执行的方法记得在这里注册下
func init() {
	addSubPlayCommand(
		demo,
		demoGorm,
	)
}

var demo = &cobra.Command{
	Use:   "demo",
	Short: "play demo演示",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("play demo")
	},
}

var demoGorm = &cobra.Command{
	Use:   "gorm",
	Short: "gorm测试",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		ylog.Debug(ctx, "GORM 测试")
		res := make([]map[string]any, 0)
		tx := g.Mysql().Raw("show tables").Find(&res)
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
		userDO := query.User.WithContext(ctx)
		userDO.Create(user)
		ylog.Debug(ctx, user)

		user, _ = userDO.First()
		ylog.Debug(ctx, user)

		user, _ = userDO.First()
		ylog.Debug(ctx, user)
	},
}
