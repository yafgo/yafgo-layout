package ylog

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Debug implements ILogger.
func (l *Logger) Debug(ctx context.Context, v ...any) {
	l.sugaredLogger(ctx).Debug(v...)
}

// Debugf implements ILogger.
func (l *Logger) Debugf(ctx context.Context, format string, v ...any) {
	l.sugaredLogger(ctx).Debugf(format, v...)
}

// Info implements ILogger.
func (l *Logger) Info(ctx context.Context, v ...any) {
	l.sugaredLogger(ctx).Info(v...)
}

// Infof implements ILogger.
func (l *Logger) Infof(ctx context.Context, format string, v ...any) {
	l.sugaredLogger(ctx).Infof(format, v...)
}

// Warn implements ILogger.
func (l *Logger) Warn(ctx context.Context, v ...any) {
	l.sugaredLogger(ctx).Warn(v...)
}

// Warnf implements ILogger.
func (l *Logger) Warnf(ctx context.Context, format string, v ...any) {
	l.sugaredLogger(ctx).Warnf(format, v...)
}

// Error implements ILogger.
func (l *Logger) Error(ctx context.Context, v ...any) {
	l.sugaredLogger(ctx).Error(v...)
}

// Errorf implements ILogger.
func (l *Logger) Errorf(ctx context.Context, format string, v ...any) {
	l.sugaredLogger(ctx).Errorf(format, v...)
}

// Panic implements ILogger.
func (l *Logger) Panic(ctx context.Context, v ...any) {
	l.sugaredLogger(ctx).Panic(v...)
}

// Panicf implements ILogger.
func (l *Logger) Panicf(ctx context.Context, format string, v ...any) {
	l.sugaredLogger(ctx).Panicf(format, v...)
}

// Fatal implements ILogger.
func (l *Logger) Fatal(ctx context.Context, v ...any) {
	l.sugaredLogger(ctx).Fatal(v...)
}

// Fatalf implements ILogger.
func (l *Logger) Fatalf(ctx context.Context, format string, v ...any) {
	l.sugaredLogger(ctx).Fatalf(format, v...)
}

// With creates a child logger and adds structured context to it. Fields added
// to the child don't affect the parent, and vice versa.
func (l *Logger) With(fields map[string]any) ILogger {
	var zapFields []zapcore.Field
	if len(fields) > 0 {
		zapFields = make([]zapcore.Field, 0, len(fields))
		for k, v := range fields {
			zapFields = append(zapFields, zap.Any(k, v))
		}
	}
	_l := &Logger{
		zl:     l.zl.With(zapFields...),
		config: l.config,
	}
	return _l
}

// WithCallerSkip creates a child logger
func (l *Logger) WithCallerSkip(skip int) ILogger {
	_l := &Logger{
		zl:     l.zl.WithOptions(zap.AddCallerSkip(skip)),
		config: l.config,
	}
	return _l
}

func (l *Logger) sugaredLogger(ctx context.Context) (zl *zap.SugaredLogger) {
	zlFields := make([]zapcore.Field, 0)
	for _, ctxKey := range l.config.CtxKeys {
		if val := ctx.Value(ctxKey); val != nil {
			zlFields = append(zlFields, zap.Any(ctxKey, val))
		}
	}
	if len(zlFields) > 0 {
		zl = l.zl.With(zlFields...).Sugar()
	} else {
		zl = l.zl.Sugar()
	}
	return
}
