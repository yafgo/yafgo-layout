package cmd

import (
	"go-toy/toy-layout/cmd/play"
)

func init() {
	// 用于play演示的命令
	playCmd := play.PlayCommand()
	if playCmd == nil {
		return
	}
	addSubCommand(subCommand{
		cmd:       playCmd,
		isDefault: false,
	})
}
