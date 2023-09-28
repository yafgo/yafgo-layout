package play

import (
	"fmt"
	"time"
	"yafgo/yafgo-layout/pkg/notify"

	"github.com/spf13/cobra"
)

// 需要执行的方法记得在这里注册下
func init() {
	addSubPlayCommand(
		feishuPost,
	)
}

var feishuPost = &cobra.Command{
	Use:   "notify",
	Short: "飞书发送测试",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("play 飞书通知")

		feishu := notify.Feishu()
		// feishu.WithRobot(robotUrl) // 自定义robot
		feishu.AtAll() // at 所有人
		// err := feishu.SendText("测试发送文本消息: %s, %d", "张三", 18)
		err := feishu.SendPost("测试富文本标题", "测试发送富文本消息: %s, %d", "张三", 18)
		fmt.Println(err)
		time.Sleep(time.Second * 1)
	},
}
