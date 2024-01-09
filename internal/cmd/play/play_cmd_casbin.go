package play

import (
	"context"
	"yafgo/yafgo-layout/pkg/perm"

	"github.com/spf13/cobra"
)

func (p *Playground) playCasbin() *cobra.Command {
	return &cobra.Command{
		Use:   "casbin",
		Short: "casbin测试",
		Run: func(cmd *cobra.Command, args []string) {

			ctx := context.Background()
			casbinUtil := perm.NewCasbinUtil(p.logger, p.db, p.q)
			enf, err := casbinUtil.Enforcer()
			if err != nil {
				p.logger.Fatalf(ctx, "enforcer: %+v", err)
			}

			var arg1 string
			if len(args) > 0 {
				arg1 = args[0]
			}

			switch arg1 {
			case "save":
				// 先清空现有规则
				enf.ClearPolicy()
				enf.SavePolicy()
				// 写入新规则
				add, addErr := enf.AddPolicies([][]string{
					{"zhangsan", "/app/demo", "read"},
					{"zhangsan", "/app/test", "write"},
					{"staff", "/app/test", "read"},
					{"staff", "/app/test", "write"},
				})
				p.logger.Infof(ctx, "addPolicy: %v, %+v", add, addErr)

			case "remove":
				// enf.RemovePolicy([]string{"zhangsan", "/app/demo"})
				enf.RemovePolicy([]string{"zhangsan", "/app/demo", "read"})

			case "add":
				enf.AddPolicy([]string{"zhangsan", "/app/demo", "read"})

			case "read":
				ok, e := enf.Enforce("admin", "/app/demo", "read")
				p.logger.Infof(ctx, "ok: %v, err: %+v", ok, e)
				ok, e = enf.Enforce("admin1", "/app/demo", "read")
				p.logger.Infof(ctx, "ok: %v, err: %+v", ok, e)
				ok, e = enf.Enforce("zhangsan", "/app/demo", "read")
				p.logger.Infof(ctx, "ok: %v, err: %+v", ok, e)
				ok, e = enf.Enforce("zhangsan", "/app/test", "write")
				p.logger.Infof(ctx, "ok: %v, err: %+v", ok, e)

			default:
				p.logger.Infof(ctx, "请指定参数")
			}
		},
	}
}
