package cmd

import (
	"yafgo/yafgo-layout/cmd/play"
	"yafgo/yafgo-layout/pkg/app"
)

func init() {
	// 用于play演示的命令
	playCmd := play.PlayCommand()
	if playCmd == nil {
		return
	}
	app.App().AddSubCommand(app.SubCommand{
		Cmd:       playCmd,
		IsDefault: false,
	})
}
