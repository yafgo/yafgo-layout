package g

import (
	"sync"
	"yafgo/yafgo-layout/pkg/app"
	"yafgo/yafgo-layout/pkg/sys/ycfg"
)

// 全局 Cfg 对象
var instCfg *ycfg.Config
var onceCfg sync.Once

// Cfg 获取全局 Cfg 对象
func Cfg() *ycfg.Config {
	// 初始化配置实例
	onceCfg.Do(func() {
		configName := app.App().Mode()
		instCfg = ycfg.New(configName,
			ycfg.WithType("yaml"),
			ycfg.WithEnvPrefix("YAFGO"),
			// ycfg.WithUnmarshalObj(g.Config),
		)
	})

	return instCfg
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
