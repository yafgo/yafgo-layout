// global 全局耦合包

package g

import (
	"yafgo/yafgo-layout/pkg/sys/ycfg"
)

type GlobalObj struct {
	cfg *ycfg.Config
}

func New(cfg *ycfg.Config) *GlobalObj {
	return &GlobalObj{
		cfg: cfg,
	}
}

// AppName 当前应用名, 用于log前缀等
func (g *GlobalObj) AppName() string {
	g.cfg.SetDefault("appname", "YAFGO")
	appname := g.cfg.GetString("appname")
	return appname
}

// AppEnv 当前环境
func (g *GlobalObj) AppEnv() string {
	g.cfg.SetDefault("env", "dev")
	return g.cfg.GetString("env")
}

// IsProd 是否生产环境
func (g *GlobalObj) IsProd() bool {
	return g.AppEnv() == "prod"
}

// IsDev 是否开发环境
func (g *GlobalObj) IsDev() bool {
	return g.AppEnv() == "dev"
}

// IsTest 是否测试环境
func (g *GlobalObj) IsTest() bool {
	return g.AppEnv() == "test"
}
