package providers

import (
	"yafgo/yafgo-layout/pkg/sys/ycfg"
	"yafgo/yafgo-layout/pkg/sys/ylog"
)

// zap 默认配置
/* [yaml 配置项]
log:
  level: debug
  encoding: console                    # json or console
  path: "./storage/log"                # log 文件目录
  filename: "log.log"                  # log 文件名
  prefix: "yafgo"                      # log 文件名
  time_format: "2006-01-02 15:04:05"   # 时间格式
  stdout: true                         # 是否输出到 stdout
  ctx_keys: ["req_id"]
  rotate:
    max_size: 64                    # 每个日志文件保存的最大尺寸,单位：M
    max_age: 7                      # 文件最多保存多少天
    max_backups: 30                 # 日志文件最多保存多少个备份
    local_time: true                # 日志文件名使用本地化时间
    compress: false                 # 是否压缩
*/
var loggerConfigDefault = map[string]any{
	"level":       "debug",         // log级别: debug < info < warn < error < panic < fatal
	"encoding":    "console",       // json or console
	"path":        "./storage/log", // 日志路径
	"filename":    "log.log",       // log 文件名
	"prefix":      "yafgo",         // 控制台输出时的前缀
	"time_format": "",              // 时间格式
	"stdout":      true,            // 是否输出到 stdout
	"ctx_keys": []string{
		"req_id",
	}, // context中附带的log字段
	"rotate": map[string]any{
		"max_size":    64,    // 每个日志文件保存的最大尺寸,单位：Mb
		"max_age":     7,     // 文件最多保存多少天
		"max_backups": 30,    // 日志文件最多保存多少个备份
		"local_time":  true,  // 日志文件名使用本地化时间
		"compress":    false, // 是否压缩
	}, // 日志滚动记录
}

// setupLogger 初始化 logger
func NewYLog(cfg *ycfg.Config) *ylog.Logger {
	// 初始默认配置
	cfg.SetDefault("log", loggerConfigDefault)
	subCfg := cfg.Sub("log")

	// logger 默认配置
	lgCfg := ylog.Config{
		TimeFormat: subCfg.GetString("time_format"),
		Encoding:   subCfg.GetString("encoding"),
		Path:       subCfg.GetString("path"),
		Filename:   subCfg.GetString("filename"),
		Level:      subCfg.GetString("level"),
		Prefix:     subCfg.GetString("prefix"),
		CtxKeys:    subCfg.GetStringSlice("ctx_keys"),
		Stdout:     subCfg.GetBool("stdout"),
		Rotate: ylog.ConfigRotate{
			MaxSize:    subCfg.GetInt("max_size"),
			MaxAge:     subCfg.GetInt("max_age"),
			MaxBackups: subCfg.GetInt("max_backups"),
			LocalTime:  subCfg.GetBool("local_time"),
			Compress:   subCfg.GetBool("compress"),
		},
	}
	if lgCfg.Prefix == "" {
		lgCfg.Prefix = "AppName()"
	}
	lg := ylog.New(lgCfg)
	// 替换ylog包的默认 logger
	ylog.SetDefaultLogger(lg)
	return lg
}
