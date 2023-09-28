package global

import (
	"yafgo/yafgo-layout/pkg/sys/ycfg"
)

// 全局 Ycfg 对象
var Ycfg *ycfg.Config

// AppName 当前应用名, 用于log前缀等
func AppName() string {
	Ycfg.SetDefault("appname", "YAFGO")
	appname := Ycfg.GetString("appname")
	return appname
}

// AppEnv 当前环境
func AppEnv() string {
	Ycfg.SetDefault("env", "dev")
	return Ycfg.GetString("env")
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
