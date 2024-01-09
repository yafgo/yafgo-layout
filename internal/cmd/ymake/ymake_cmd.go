package ymake

import (
	"github.com/spf13/cobra"
)

func (p *codeMaker) Command() *cobra.Command {
	var cmdExample = `go run . make repo user
go run . make service user
go run . make handler user
go run . make handler backend/user
go run . make handler frontend/user`
	var CmdMake = &cobra.Command{
		Use:     "make [type] [name]",
		Short:   "Make a new handler/service/repository",
		Example: cmdExample,
		Args:    cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	CmdMake.AddCommand(p.commandCreateRepo())
	CmdMake.AddCommand(p.commandCreateService())
	CmdMake.AddCommand(p.commandCreateHandler())
	return CmdMake
}

func (p *codeMaker) commandCreateRepo() *cobra.Command {
	return &cobra.Command{
		Use:     "repo [name]",
		Short:   "Make a new repository",
		Example: "go run . make repo user",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return p.Make(TypeRepository, args)
		},
	}
}

func (p *codeMaker) commandCreateService() *cobra.Command {
	return &cobra.Command{
		Use:     "service [name]",
		Short:   "Make a new service",
		Example: "go run . make service user",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return p.Make(TypeService, args)
		},
	}
}

func (p *codeMaker) commandCreateHandler() *cobra.Command {
	var cmdExample = `go run . make handler user
go run . make handler backend/user
go run . make handler frontend/user`
	return &cobra.Command{
		Use:     "handler [name]",
		Short:   "Make a new handler",
		Example: cmdExample,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return p.Make(TypeHandler, args)
		},
	}
}
