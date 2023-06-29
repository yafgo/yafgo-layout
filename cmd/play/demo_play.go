package play

import (
	"fmt"

	"github.com/spf13/cobra"
)

// 需要执行的方法记得在这里注册下
func init() {
	addSubPlayCommand(
		demo,
	)
}

var demo = &cobra.Command{
	Use:   "demo",
	Short: "play demo演示",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("play demo")
	},
}
