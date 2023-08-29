package ylog

import "context"

// ILogger is the API interface for logger.
type ILogger interface {
	Debug(ctx context.Context, v ...any)
	Debugf(ctx context.Context, format string, v ...any)
	Info(ctx context.Context, v ...any)
	Infof(ctx context.Context, format string, v ...any)
	Warn(ctx context.Context, v ...any)
	Warnf(ctx context.Context, format string, v ...any)
	Error(ctx context.Context, v ...any)
	Errorf(ctx context.Context, format string, v ...any)
	Panic(ctx context.Context, v ...any)
	Panicf(ctx context.Context, format string, v ...any)
	Fatal(ctx context.Context, v ...any)
	Fatalf(ctx context.Context, format string, v ...any)

	With(fields ...Field) ILogger
	WithCallerSkip(skip int) ILogger
	AddCommonField(key string, val any) ILogger
}

var (
	// Ensure Logger implements ILogger.
	_ ILogger = &Logger{}

	// Default logger object, for package method usage.
	defaultLogger = New()
)

// DefaultLogger returns the default logger.
func DefaultLogger() *Logger {
	return defaultLogger
}

// SetDefaultLogger sets the default logger for package glog.
// Note that there might be concurrent safety issue if calls this function
// in different goroutines.
func SetDefaultLogger(l *Logger) {
	defaultLogger = l
}
