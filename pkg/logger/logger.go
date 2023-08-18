package logger

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// zap 默认配置
/* [yaml 配置项]
log:
  log_level: debug
  encoding: console               # json or console
  log_filename: "./logs/log.log"  # log 文件名
  max_backups: 30                 # 日志文件最多保存多少个备份
  max_age: 7                      # 文件最多保存多少天
  max_size: 64                    # 每个日志文件保存的最大尺寸,单位：M
  compress: true                  # 是否压缩
*/
var zapConfigDefault = map[string]any{
	"log_level":    "debug",          // log级别: debug < info < warn < error < fatal < panic
	"encoding":     "console",        // json or console
	"log_filename": "./logs/log.log", // log 文件名
	"max_backups":  30,               // 日志文件最多保存多少个备份
	"max_age":      7,                // 文件最多保存多少天
	"max_size":     64,               // 每个日志文件保存的最大尺寸,单位：Mb
	"compress":     true,             // 是否压缩
}

var Logger *logger // 全局 Logger 对象

func SetupDefault(conf *viper.Viper, opts ...loggerOption) *logger {
	Logger = New(conf, opts...)
	return Logger
}

type logger struct {
	*zap.Logger

	isProd bool         // 是否生产模式
	prefix string       // log 前缀
	cfg    *viper.Viper // 配置对象
}

type loggerOption func(*logger)

// New 获取新的logger实例
func New(conf *viper.Viper, opts ...loggerOption) *logger {
	// 初始默认配置
	conf.SetDefault("log", zapConfigDefault)
	Logger = initZap(conf, opts...)
	return Logger
}

// WithIsProd 是否生成模式
func WithIsProd(isProd bool) loggerOption {
	return func(p *logger) {
		p.isProd = isProd
	}
}

// WithPrefix 设置log前缀
func WithPrefix(prefix string) loggerOption {
	return func(p *logger) {
		p.prefix = prefix
	}
}

func initZap(conf *viper.Viper, opts ...loggerOption) *logger {
	lg := &logger{cfg: conf}
	for _, opt := range opts {
		opt(lg)
	}

	// 日志级别 DEBUG,ERROR,INFO
	lv := conf.GetString("log.log_level")
	var level zapcore.Level
	// debug < info < warn < error < fatal < panic
	switch lv {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}

	// 初始化 core
	encoder := lg.getEncoder(lg.isConsole())
	writeSyncer := lg.getLogWriter()
	core := zapcore.NewCore(encoder, writeSyncer, level)

	// 初始化 Logger
	var zapLogger *zap.Logger
	var zapOpts = []zap.Option{
		zap.AddCaller(),                   // 调用文件和行号，内部使用 runtime.Caller
		zap.AddCallerSkip(1),              // 封装了一层，调用文件去除一层(runtime.Caller(1))
		zap.AddStacktrace(zap.ErrorLevel), // Error 时才会显示 stacktrace
	}
	if !lg.isProd {
		zapOpts = append(zapOpts, zap.Development())
	}
	zapLogger = zap.New(core, zapOpts...)
	zap.ReplaceGlobals(zapLogger)
	return &logger{Logger: zapLogger}
}

// getEncoder 设置日志存储格式
func (lg *logger) getEncoder(isConsole bool) zapcore.Encoder {
	// 日志格式规则
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller", // 代码调用，如 paginator/paginator.go:148
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,      // 每行日志的结尾添加 "\n"
		EncodeLevel:    zapcore.CapitalLevelEncoder,    // 日志级别名称大写，如 ERROR、INFO
		EncodeTime:     lg.timeEncoder,                 // 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, // 执行时间，以秒为单位
		EncodeCaller:   zapcore.ShortCallerEncoder,     // Caller 短格式，如：types/converter.go:17，长格式为绝对路径
	}

	// console 模式使用 Console 编码器
	if isConsole {
		// 终端输出的关键词高亮
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		// 内置的 Console 编码器（支持 stacktrace 换行）
		return zapcore.NewConsoleEncoder(encoderConfig)
	}

	// JSON 编码器
	return zapcore.NewJSONEncoder(encoderConfig)
}

// timeEncoder 自定义时间格式
func (lg *logger) timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	prefix := lg.prefix
	if prefix != "" {
		prefix = "[" + prefix + "]"
	}
	if lg.isConsole() {
		enc.AppendString(prefix + t.Format(time.DateTime+".000000"))
	} else {
		enc.AppendString(t.Format(time.DateTime + ".000000"))
	}
}

// getLogWriter 日志记录介质
func (lg *logger) getLogWriter() zapcore.WriteSyncer {
	conf := lg.cfg

	// 日志文件
	filename := conf.GetString("log.log_filename")

	// 滚动日志
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,                       // 日志文件路径
		MaxSize:    conf.GetInt("log.max_size"),    // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: conf.GetInt("log.max_backups"), // 日志文件最多保存多少个备份
		MaxAge:     conf.GetInt("log.max_age"),     // 文件最多保存多少天
		Compress:   conf.GetBool("log.compress"),   // 是否压缩
	}

	// 配置输出介质
	if lg.isConsole() {
		// 本地开发: 终端打印和记录文件
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger))
	} else {
		// 生产环境只记录文件
		// return zapcore.AddSync(lumberJackLogger)
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger))
	}
}

// isConsole 是否控制台输出
func (lg *logger) isConsole() bool {
	return lg.cfg.GetString("log.encoding") == "console"
}

const LOGGER_CTX_KEY = "zapLogger"

// NewContext 给指定的context添加字段
func (lg *logger) NewContext(ctx *gin.Context, fields ...zapcore.Field) {
	ctx.Set(LOGGER_CTX_KEY, lg.WithContext(ctx).With(fields...))
}

// WithContext 从指定的context返回一个 logger 实例
func (lg *logger) WithContext(ctx *gin.Context) *logger {
	if ctx == nil {
		return lg
	}
	zl, exists := ctx.Get(LOGGER_CTX_KEY)
	if !exists {
		return lg
	}
	ctxLogger, ok := zl.(*zap.Logger)
	if ok {
		return &logger{
			Logger: ctxLogger,
			isProd: lg.isProd,
			prefix: lg.prefix,
			cfg:    lg.cfg,
		}
	}
	return lg
}
