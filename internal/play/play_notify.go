package play

import (
	"fmt"

	"github.com/spf13/cobra"
)

func (p *Playground) playFeishu() *cobra.Command {
	return &cobra.Command{
		Use:   "notify",
		Short: "飞书发送测试",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("play 飞书通知")

			p.feishu.New().
				// WithRobot(robotUrl). // 自定义robot
				AtAll(). // at 所有人
				SendText("测试发送文本消息: %s, %d", "张三", 18)

			p.feishu.New().
				AtAll().
				SendPost("测试富文本标题", "测试发送富文本消息: %s, %d", "张三", 18)
			// time.Sleep(time.Second * 1)
		},
	}
}
