package g

import (
	"sync"
	"yafgo/yafgo-layout/pkg/app"
	"yafgo/yafgo-layout/pkg/sys/ycfg"
)

// Cfg 获取全局 Cfg 对象
func Cfg() *ycfg.Config {
	globalCfg := newCfg()(app.App().Mode())
	return globalCfg
}

func newCfg() func(configName string) *ycfg.Config {
	// 全局 Cfg 对象
	var once sync.Once
	var cfgInst *ycfg.Config

	return func(configName string) *ycfg.Config {

		// 初始化配置实例
		once.Do(func() {
			cfgInst = ycfg.New(configName,
				ycfg.WithType("yaml"),
				ycfg.WithEnvPrefix("YAFGO"),
				// ycfg.WithUnmarshalObj(g.Config),
			)
		})

		return cfgInst
	}
}

// AppName 当前应用名, 用于log前缀等
func AppName() string {
	Cfg().SetDefault("appname", "YAFGO")
	appname := Cfg().GetString("appname")
	return appname
}

// AppEnv 当前环境
func AppEnv() string {
	Cfg().SetDefault("env", "dev")
	return Cfg().GetString("env")
}

// IsProd 是否生产环境
func IsProd() bool {
	return AppEnv() == "prod"
}

// IsDev 是否开发环境
func IsDev() bool {
	return AppEnv() == "dev"
}

// IsTest 是否测试环境
func IsTest() bool {
	return AppEnv() == "test"
}
