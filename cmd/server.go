package cmd

import (
	"go-toy/toy-layout/internal/server"

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

	addSubCommand(subCommand{
		cmd:       subCmd,
		isDefault: true,
	})
}
