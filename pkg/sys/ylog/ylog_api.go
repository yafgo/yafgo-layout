package ylog

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
	_lg := &Logger{
		zl: defaultLogger.zl.WithOptions(
			// 当前文件封装了一层, 所以重置下skip, 否则打印的行号始终是当前文件的, 而非真正调用的位置
			zap.AddCallerSkip(1),
		),
		config: defaultLogger.config,
	}
	return _lg
}

// With creates a child logger and adds structured context to it. Fields added
// to the child don't affect the parent, and vice versa.
func With(fields map[string]any) ILogger {
	var zapFields []zapcore.Field
	if len(fields) > 0 {
		zapFields = make([]zapcore.Field, 0, len(fields))
		for k, v := range fields {
			zapFields = append(zapFields, zap.Any(k, v))
		}
	}
	_l := &Logger{
		zl:     defaultLogger.zl.With(zapFields...),
		config: defaultLogger.config,
	}
	return _l
}

// WithCallerSkip creates a child logger
func WithCallerSkip(skip int) ILogger {
	return defaultLogger.WithCallerSkip(skip)
}
