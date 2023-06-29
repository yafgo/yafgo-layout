package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/gookit/color"
	"github.com/spf13/viper"
)

var conf *viper.Viper
var (
	envPrefix = ""
	configDir = "."
)

func Config() *viper.Viper {
	if conf == nil {
		color.Errorln("Config尚未初始化")
		os.Exit(1)
	}
	return conf
}

func SetEnvPrefix(val string) {
	envPrefix = val
}

func SetConfigDir(val string) {
	configDir = val
}

// SetupViper 初始化配置
//
//	配置加载顺序: [mode].yaml -> [mode].local.yaml
//	故配置优先级: [mode].local.yaml  >  [mode].yaml
//	可以自定义任意 mode, 只需要新增 [mode].yaml 文件即可
func SetupConfig(mode string) *viper.Viper {
	var configName string

	// 初始化默认配置
	conf = viper.New()
	conf.SetConfigType("yaml")
	conf.AddConfigPath(configDir)

	// 读取环境变量, 设置前缀, 后续获取的时候不需要加前缀
	//   如前缀为 MYAPP 时: export MYAPP_WS_ADDR=ws://127.0.0.1:8080
	//   _viper.GetString("WS_ADDR") 即可获取
	//   _viper.GetString("ws_addr") 也可以获取
	conf.SetEnvPrefix(envPrefix)
	conf.AutomaticEnv()

	defer func() {
		// 监听配置文件变化
		conf.WatchConfig()
		conf.OnConfigChange(func(e fsnotify.Event) {
			color.Infoln("config file changed:", e.Name)
		})
	}()

	// 1. 读取指定 mode 配置
	if mode == "" {
		panic("Error: mode is empty")
	}
	if _f := filepath.Join(configDir, mode+".yaml"); fileExist(_f) {
		configName = mode
		// 读取指定 mode 配置
		conf.SetConfigName(configName)
		color.Successln("MergeInConfig:", _f)
		if err := conf.MergeInConfig(); err != nil {
			panic(fmt.Errorf("MergeInConfig Error [%s]: %+v", configName, err))
		}
	} else {
		panic(fmt.Sprintf("Error: [%s] not exist", configName))
	}

	// 2. [如果local配置存在]
	// 读取 local 配置, 不进版本库, 如: dev.local.yaml, prod.local.yaml
	if _f := filepath.Join(configDir, mode+".local.yaml"); fileExist(_f) {
		configName = mode + ".local"
		// 读取 local 配置
		conf.SetConfigName(configName)
		color.Successln("MergeInConfig:", _f)
		if err := conf.MergeInConfig(); err != nil {
			panic(fmt.Errorf("MergeInConfig Error [%s]: %+v", configName, err))
		}
	}

	return conf
}

func fileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}
