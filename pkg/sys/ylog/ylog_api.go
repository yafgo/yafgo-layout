package ylog

import (
	"context"

	"go.uber.org/zap"
)

// Debug implements ILogger.
func Debug(ctx context.Context, v ...any) {
	_ylog().Debug(ctx, v...)
}

// Debugf implements ILogger.
func Debugf(ctx context.Context, format string, v ...any) {
	_ylog().Debugf(ctx, format, v...)
}

// Info implements ILogger.
func Info(ctx context.Context, v ...any) {
	_ylog().Info(ctx, v...)
}

// Infof implements ILogger.
func Infof(ctx context.Context, format string, v ...any) {
	_ylog().Infof(ctx, format, v...)
}

// Warn implements ILogger.
func Warn(ctx context.Context, v ...any) {
	_ylog().Warn(ctx, v...)
}

// Warnf implements ILogger.
func Warnf(ctx context.Context, format string, v ...any) {
	_ylog().Warnf(ctx, format, v...)
}

// Error implements ILogger.
func Error(ctx context.Context, v ...any) {
	_ylog().Error(ctx, v...)
}

// Errorf implements ILogger.
func Errorf(ctx context.Context, format string, v ...any) {
	_ylog().Errorf(ctx, format, v...)
}

// Panic implements ILogger.
func Panic(ctx context.Context, v ...any) {
	_ylog().Panic(ctx, v...)
}

// Panicf implements ILogger.
func Panicf(ctx context.Context, format string, v ...any) {
	_ylog().Panicf(ctx, format, v...)
}

// Fatal implements ILogger.
func Fatal(ctx context.Context, v ...any) {
	_ylog().Fatal(ctx, v...)
}

// Fatalf implements ILogger.
func Fatalf(ctx context.Context, format string, v ...any) {
	_ylog().Fatalf(ctx, format, v...)
}

func _ylog() ILogger {
	_lg := defaultLogger.copy()
	_lg.zl = defaultLogger.zl.WithOptions(
		// 当前文件封装了一层, 所以重置下skip, 否则打印的行号始终是当前文件的, 而非真正调用的位置
		zap.AddCallerSkip(1),
	)
	return _lg
}

// With 创建一个子logger, 添加到该子logger的字段不会影响父级, 反之亦然
func With(fields ...Field) ILogger {
	_l := defaultLogger.With(fields...)
	return _l
}

// WithCallerSkip creates a child logger
func WithCallerSkip(skip int) ILogger {
	return defaultLogger.WithCallerSkip(skip)
}

func AddCommonField(key string, val any) ILogger {
	defaultLogger.commonFields = append(defaultLogger.commonFields, Field{
		Key: key,
		Val: val,
	})
	return defaultLogger
}
