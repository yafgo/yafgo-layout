package play

import (
	"context"
	"fmt"
	"yafgo/yafgo-layout/internal/global"
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
		ylog.Debug(ctx, "moduleName string")
		res := make([]map[string]any, 0)
		tx := global.Mysql.Raw("show tables").Find(&res)
		if tx.Error != nil {
			ylog.Errorf(ctx, "err: %v", tx.Error)
			return
		}
		ylog.Debug(ctx, res)

		user := &model.User{
			Name: "张三",
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
