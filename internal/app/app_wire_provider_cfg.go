package app

import (
	"yafgo/yafgo-layout/pkg/sys/ycfg"
)

func NewYCfg(envConf string) (cfg *ycfg.Config) {
	cfg = ycfg.New(envConf,
		ycfg.WithType("yaml"),
		ycfg.WithEnvPrefix("YAFGO"),
		ycfg.WithNacosEnabled(true),
	)
	return
}

func NewYCfgDig(envConf string) func() *ycfg.Config {
	return func() (cfg *ycfg.Config) {
		return NewYCfg(envConf)
	}
}
