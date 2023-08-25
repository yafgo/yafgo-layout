package global

import (
	"go-toy/toy-layout/pkg/sys/ycfg"
)

// 全局 Ycfg 对象
var Ycfg *ycfg.Config

// AppName 当前应用名, 用于log前缀等
func AppName() string {
	Ycfg.SetDefault("appname", "YAFGO")
	appname := Ycfg.GetString("appname")
	return appname
}

// IsProd 是否生产环境
func IsProd() bool {
	Ycfg.SetDefault("env", "dev")
	_env := Ycfg.GetString("env")
	return _env == "prod"
}

// IsDev 是否开发环境
func IsDev() bool {
	Ycfg.SetDefault("env", "dev")
	_env := Ycfg.GetString("env")
	return _env == "dev"
}

// IsTest 是否测试环境
func IsTest() bool {
	Ycfg.SetDefault("env", "dev")
	_env := Ycfg.GetString("env")
	return _env == "test"
}
