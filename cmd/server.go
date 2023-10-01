package cmd

import (
	"yafgo/yafgo-layout/internal/server"
	"yafgo/yafgo-layout/pkg/app"

	"github.com/spf13/cobra"
)

func init() {
	// web 服务
	srv := server.NewWebService()
	var subCmd = &cobra.Command{
		Use:   "serve",
		Short: "Run WebServer",
		Args:  cobra.NoArgs,
		Run:   srv.CmdRun,
	}

	app.App().AddSubCommand(app.SubCommand{
		Cmd:       subCmd,
		IsDefault: true,
	})
}
