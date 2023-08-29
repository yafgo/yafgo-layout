package ylog

import (
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	parent *Logger
	zl     *zap.Logger
	config Config

	commonFields []Field
	fields       []Field
}

func New(cfg ...Config) *Logger {
	var _cfg Config = DefaultConfig()
	if len(cfg) > 0 {
		_cfg = cfg[0]
	}
	lg := newZap(_cfg)
	return lg
}

func newZap(cfg Config) *Logger {
	lg := &Logger{config: cfg}

	// 日志级别: debug < info < warn < error < panic < fatal
	var level = mapZapLevel(lg.config.Level)

	// 初始化 core
	encoder := lg.getEncoder(cfg.Stdout)
	writeSyncer := lg.getLogWriter()
	core := zapcore.NewCore(encoder, writeSyncer, level)

	// 初始化 Logger
	var zapLogger *zap.Logger
	var zapOpts = []zap.Option{
		zap.AddCaller(),                   // 调用文件和行号，内部使用 runtime.Caller
		zap.AddCallerSkip(1),              // 封装了一层，调用文件去除一层(runtime.Caller(1))
		zap.AddStacktrace(zap.ErrorLevel), // Error 时才会显示 stacktrace
	}
	// zapOpts = append(zapOpts, zap.Development())
	zapLogger = zap.New(core, zapOpts...)
	zap.ReplaceGlobals(zapLogger)
	lg.zl = zapLogger
	return lg
}

// getEncoder 设置日志存储格式
func (lg *Logger) getEncoder(isConsole bool) zapcore.Encoder {
	// 日志格式规则
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "Logger",
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

	// JSON 编码器
	if lg.config.Encoding == LogTypeJson {
		return zapcore.NewJSONEncoder(encoderConfig)
	}

	// console 模式使用 Console 编码器
	// 终端输出的关键词高亮
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	// 内置的 Console 编码器（支持 stacktrace 换行）
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// timeEncoder 自定义时间格式
func (lg *Logger) timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	prefix := lg.config.Prefix
	if prefix != "" {
		prefix = "[" + prefix + "]"
	}
	timeFormat := lg.config.TimeFormat
	if timeFormat == "" {
		timeFormat = time.DateTime
	}
	if lg.config.Encoding == LogTypeConsole {
		enc.AppendString(prefix + t.Format(timeFormat))
	} else {
		enc.AppendString(t.Format(timeFormat))
	}
}

// getLogWriter 日志记录介质
func (lg *Logger) getLogWriter() zapcore.WriteSyncer {
	conf := lg.config

	// 滚动日志
	filename := filepath.Join(conf.Path, conf.Filename)
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,               // 日志文件路径
		MaxSize:    conf.Rotate.MaxSize,    // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: conf.Rotate.MaxBackups, // 日志文件最多保存多少个备份
		MaxAge:     conf.Rotate.MaxAge,     // 文件最多保存多少天
		LocalTime:  conf.Rotate.LocalTime,  // 是否本地时间
		Compress:   conf.Rotate.Compress,   // 是否压缩
	}

	// 配置输出介质
	writeSyncers := []zapcore.WriteSyncer{
		zapcore.AddSync(lumberJackLogger),
	}
	if conf.Stdout {
		writeSyncers = append(writeSyncers, zapcore.AddSync(os.Stdout))
	}
	return zapcore.NewMultiWriteSyncer(writeSyncers...)
}
