package ycfg

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/gookit/color"
	"github.com/spf13/viper"
)

type Config struct {
	*viper.Viper

	shouldWatch   bool
	unmarshalObj  any
	onChangeFuncs []func()

	redisStore *redisStore // redis配置存储支持
}

// RegisterOnChangeHandler 注册 onChange 回调处理
func (p *Config) RegisterOnChangeHandler(fn func()) {
	p.onChangeFuncs = append(p.onChangeFuncs, fn)
}

// EmitChange 触发 change 事件
func (p *Config) EmitChange() {
	for _, fn := range p.onChangeFuncs {
		go fn()
	}
}

// MergeConfig 从 io.Reader 合并配置, 例如从nacos读到的配置等
func (p *Config) MergeConfig(in io.Reader) (err error) {
	if err = p.Viper.MergeConfig(in); err != nil {
		return
	}
	p.EmitChange()
	return
}

func (p *Config) doUnmarshal() {
	if p.unmarshalObj == nil {
		return
	}
	if err := p.Viper.Unmarshal(p.unmarshalObj); err != nil {
		color.Errorln(err)
	}
}

type ycfg struct {
	envPrefix  string
	configType string
	configDir  string

	shouldWatch  bool
	unmarshalObj any
	nacosEnabled bool // 是否启用nacos配置支持

	configDirs []string
}

type ycfgOption func(*ycfg)

// New 获取新的配置实例
//
//	默认会加载配置目录下的 "config.{配置类型}"
//	name 不为空时会继续尝试加载配置目录下的 "{name}.{配置类型}"
func New(name string, opts ...ycfgOption) *Config {
	var builder = &ycfg{
		envPrefix:  "yafgo",
		configDir:  ".",
		configType: "yaml",

		configDirs: []string{"config", "."},
	}

	for _, opt := range opts {
		opt(builder)
	}

	var cfg = builder.setupConfig(name)
	// 启用nacos配置支持
	if builder.nacosEnabled {
		cfg.setupNacos()
	}
	return cfg
}

// WithDir 自定义配置文件查找目录
//
//	默认依次在 "{cwd}/config/"、"{cwd}/" 目录下查找
func WithDir(val string) ycfgOption {
	return func(p *ycfg) {
		if val == "" || val == "." {
			return
		}
		p.configDir = val
		p.configDirs = append(p.configDirs, val)
	}
}

// WithType 自定义配置文件类型
//
//	[]string{"json", "toml", "yaml", "yml", "properties", "props", "prop", "hcl", "tfvars", "dotenv", "env", "ini"}
func WithType(val string) ycfgOption {
	return func(p *ycfg) {
		if val == "" {
			return
		}
		for _, ext := range viper.SupportedExts {
			if ext == val {
				p.configType = val
				break
			}
		}
	}
}

// WithEnvPrefix 自定义环境变量前缀, 默认: "yafgo"
//
//	如前缀为 "myapp" 时: export MYAPP_SERVER_PORT=8080
//	_viper.GetInt64("SERVER_PORT") 即可获取
//	_viper.GetString("server_port") 也可以获取
func WithEnvPrefix(val string) ycfgOption {
	return func(p *ycfg) {
		if val == "" {
			return
		}
		p.envPrefix = val
	}
}

// WithUnmarshal 解析配置到指定变量
func WithUnmarshalObj(obj any) ycfgOption {
	return func(p *ycfg) {
		if obj == nil {
			return
		}
		p.unmarshalObj = obj
	}
}

// WithShouldWatch 是否监听文件变动, 默认否
func WithShouldWatch(flag bool) ycfgOption {
	return func(p *ycfg) {
		p.shouldWatch = flag
	}
}

// WithNacosEnabled 是否启用nacos配置支持, 默认否
func WithNacosEnabled(enable bool) ycfgOption {
	return func(p *ycfg) {
		p.nacosEnabled = enable
	}
}

// SetupViper 初始化配置
//
//	配置加载顺序: [mode].yaml -> [mode].local.yaml
//	故配置优先级: [mode].local.yaml  >  [mode].yaml
//	可以自定义任意 mode, 只需要新增 [mode].yaml 文件即可
func (p *ycfg) setupConfig(name ...string) *Config {
	var cfg = &Config{
		unmarshalObj:  p.unmarshalObj,
		shouldWatch:   p.shouldWatch,
		onChangeFuncs: []func(){},
	}
	cfg.onChangeFuncs = append(cfg.onChangeFuncs, cfg.doUnmarshal)
	cfg.Viper = viper.New()
	var _viper = cfg.Viper

	// 配置
	var configFileNames = []string{"config"}
	if len(name) > 0 && name[0] != "" {
		configFileNames = append(configFileNames, name[0])
	}

	// 初始化默认配置
	_viper.SetConfigType(p.configType)
	_viper.AddConfigPath(p.configDir)

	// 读取环境变量, 设置前缀, 后续获取的时候不需要加前缀
	//   如前缀为 MYAPP 时: export MYAPP_WS_ADDR=ws://127.0.0.1:8080
	//   _viper.GetString("WS_ADDR") 即可获取
	//   _viper.GetString("ws_addr") 也可以获取
	_viper.SetEnvPrefix(p.envPrefix)
	_viper.AutomaticEnv()

	// [本地配置支持]
	// 如果local配置存在, 读取 local 配置, 该配置文件不进版本库, 如: dev.local.yaml, prod.local.yaml
	var _cfgFileNames = make([]string, 0, len(configFileNames)*2)
	for _, v := range configFileNames {
		_cfgFileNames = append(_cfgFileNames, v, v+".local")
	}

	// 最后生效的配置文件
	var configNameUsed string
	var configFileUsed string

	// 循环尝试读取各个路径下可能存在的配置文件
	for _, dir := range p.configDirs {
		for _, cfgFileName := range _cfgFileNames {
			cfgName := filepath.Join(dir, cfgFileName+"."+p.configType)
			_viper.SetConfigName(cfgName)
			if err := _viper.MergeInConfig(); err != nil {
				// 恢复最后一次正确配置文件名, 否则后续监听文件变化不可用
				_viper.SetConfigName(configNameUsed)
				if _, ok := err.(viper.ConfigFileNotFoundError); ok {
					// 尝试读取的配置文件不存在
					continue
				} else {
					panic(fmt.Errorf("MergeInConfig Error [%s]: %+v", cfgName, err))
				}
			}
			configNameUsed = cfgName
			configFileUsed = _viper.ConfigFileUsed()
			color.Successln("MergeInConfig:", configFileUsed)
		}
	}

	if configFileUsed == "" {
		color.Warnln("没有加载任何配置")
		return cfg
	}

	cfg.doUnmarshal()
	if p.shouldWatch {
		// 监听配置文件变化
		_viper.WatchConfig()
		_viper.OnConfigChange(func(e fsnotify.Event) {
			color.Infoln("config file changed:", e.Name)
			cfg.EmitChange()
		})
	}

	return cfg
}
