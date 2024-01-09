package providers

import (
	"yafgo/yafgo-layout/pkg/sys/ycfg"

	"github.com/spf13/cobra"
)

// YCfgName
type YCfgName string

func NewYCfg(cfgName YCfgName) (cfg *ycfg.Config) {
	cfg = ycfg.New(string(cfgName),
		ycfg.WithType("yaml"),
		ycfg.WithEnvPrefix("YAFGO"),
		ycfg.WithNacosEnabled(true),
	)
	return
}

// ParseConfigName 解析启动参数中的配置名称
func ParseConfigName() (confName YCfgName) {

	// 仅负责解析配置参数的 cobra 实例
	var preCmd = &cobra.Command{}

	// 禁用 -h 标志的响应
	preCmd.SetHelpFunc(func(c *cobra.Command, s []string) {})

	// 解析配置参数
	var _cfgName string
	preCmd.PersistentFlags().StringVarP(&_cfgName, "conf", "c", "dev", "ConfigName: dev, prod...")

	// 执行主命令(这里忽略Execute的错误)
	preCmd.Execute()
	confName = YCfgName(_cfgName)

	return
}
