package play

import (
	"context"
	"yafgo/yafgo-layout/pkg/hash"
	"yafgo/yafgo-layout/pkg/sys/ylog"

	"github.com/spf13/cobra"
)

func (p *Playground) playPassword() *cobra.Command {
	return &cobra.Command{
		Use:   "password",
		Short: "密码加密",
		Args:  cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			if len(args) == 0 {
				p.logger.Warn(ctx, "请输入原始密码")
				return
			}
			password := args[0]
			passwordHashed, err := hash.BcryptHash(args[0])
			p.logger.With(
				ylog.Any("原始", password),
				ylog.Any("加密", passwordHashed),
				ylog.Any("err", err),
			).Info(ctx, "密码bcrypt处理")
		},
	}
}
