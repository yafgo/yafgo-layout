package ylog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	LEVEL_DEBUG = "debug" // DebugLevel logs are typically voluminous, and are usually disabled in production.
	LEVEL_INFO  = "info"  // InfoLevel is the default logging priority.
	LEVEL_WARN  = "warn"  // WarnLevel logs are more important than Info, but don't need individual human review.
	LEVEL_ERROR = "error" // ErrorLevel logs are high-priority. If an application is running smoothly, it shouldn't generate any error-level logs.
	LEVEL_PANIC = "panic" // PanicLevel logs a message, then panics.
	LEVEL_FATAL = "fatal" // FatalLevel logs a message, then calls os.Exit(1).
)

var (
	_zapLevelMap = map[string]zapcore.Level{
		LEVEL_DEBUG: zap.DebugLevel,
		LEVEL_INFO:  zap.InfoLevel,
		LEVEL_WARN:  zap.WarnLevel,
		LEVEL_ERROR: zap.ErrorLevel,
		LEVEL_PANIC: zap.PanicLevel,
		LEVEL_FATAL: zap.FatalLevel,
	}
)

func mapZapLevel(lv string) zapcore.Level {
	if val, ok := _zapLevelMap[lv]; ok {
		return val
	}
	return zap.InfoLevel
}
