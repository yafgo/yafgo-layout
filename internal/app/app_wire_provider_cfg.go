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
