// global 全局耦合包

package g

import (
	"yafgo/yafgo-layout/pkg/sys/ycfg"
)

// 全局 Cfg 对象
var instCfg *ycfg.Config

// Cfg 获取全局 Cfg 对象
func SetCfg(cfg *ycfg.Config) {
	instCfg = cfg
}

// AppName 当前应用名, 用于log前缀等
func AppName() string {
	instCfg.SetDefault("appname", "YAFGO")
	appname := instCfg.GetString("appname")
	return appname
}

// AppEnv 当前环境
func AppEnv() string {
	instCfg.SetDefault("env", "dev")
	return instCfg.GetString("env")
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
