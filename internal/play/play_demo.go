package play

import (
	"fmt"

	"github.com/spf13/cobra"
)

func (p *Playground) playDemo() *cobra.Command {
	return &cobra.Command{
		Use:   "demo",
		Short: "play demo演示",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("play demo")
		},
	}
}
