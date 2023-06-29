package play

import (
	"fmt"

	"github.com/spf13/cobra"
)

// subPlayCmds 本包下需要注册到 play 命令的子命令
var subPlayCmds = make([]*cobra.Command, 0)

// addSubPlayCommand 添加二级命令
func addSubPlayCommand(cmds ...*cobra.Command) {
	subPlayCmds = append(subPlayCmds, cmds...)
}

// PlayCommand 返回 play 主命令
func PlayCommand() *cobra.Command {
	playCmd := &cobra.Command{
		Use:   "play",
		Short: "A playground for testing",
		Run:   _runPlay,
	}
	playCmd.AddCommand(subPlayCmds...)
	return playCmd
}

func _runPlay(cmd *cobra.Command, args []string) {

	fmt.Println("\n[play test case...]")

	// 可以临时在这里调用一些代码
}
