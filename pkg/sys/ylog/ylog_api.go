package ylog

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Debug implements ILogger.
func Debug(ctx context.Context, v ...any) {
	defaultLogger.Debug(ctx, v...)
}

// Debugf implements ILogger.
func Debugf(ctx context.Context, format string, v ...any) {
	defaultLogger.Debugf(ctx, format, v...)
}

// Info implements ILogger.
func Info(ctx context.Context, v ...any) {
	defaultLogger.Info(ctx, v...)
}

// Infof implements ILogger.
func Infof(ctx context.Context, format string, v ...any) {
	defaultLogger.Infof(ctx, format, v...)
}

// Warn implements ILogger.
func Warn(ctx context.Context, v ...any) {
	defaultLogger.Warn(ctx, v...)
}

// Warnf implements ILogger.
func Warnf(ctx context.Context, format string, v ...any) {
	defaultLogger.Warnf(ctx, format, v...)
}

// Error implements ILogger.
func Error(ctx context.Context, v ...any) {
	defaultLogger.Error(ctx, v...)
}

// Errorf implements ILogger.
func Errorf(ctx context.Context, format string, v ...any) {
	defaultLogger.Errorf(ctx, format, v...)
}

// Panic implements ILogger.
func Panic(ctx context.Context, v ...any) {
	defaultLogger.Panic(ctx, v...)
}

// Panicf implements ILogger.
func Panicf(ctx context.Context, format string, v ...any) {
	defaultLogger.Panicf(ctx, format, v...)
}

// Fatal implements ILogger.
func Fatal(ctx context.Context, v ...any) {
	defaultLogger.Fatal(ctx, v...)
}

// Fatalf implements ILogger.
func Fatalf(ctx context.Context, format string, v ...any) {
	defaultLogger.Fatalf(ctx, format, v...)
}

// With creates a child logger and adds structured context to it. Fields added
// to the child don't affect the parent, and vice versa.
func With(fields map[string]any) *Logger {
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
