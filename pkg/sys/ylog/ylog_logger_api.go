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

func (l *Logger) copy() *Logger {
	return &Logger{
		zl:           l.zl,
		parent:       l,
		config:       l.config,
		fields:       l.fields,
		commonFields: l.commonFields,
	}
}

// WithCallerSkip creates a child logger
func (l *Logger) WithCallerSkip(skip int) ILogger {
	_l := l.copy()
	_l.zl = l.zl.WithOptions(zap.AddCallerSkip(skip))
	return _l
}

func (l *Logger) AddCommonField(key string, val any) ILogger {
	l.commonFields = append(l.commonFields, Field{
		Key: key,
		Val: val,
	})
	return l
}

func (l *Logger) sugaredLogger(ctx context.Context) (zl *zap.SugaredLogger) {
	zlFields := make([]zapcore.Field, 0)
	// commonFields
	for _, field := range l.commonFields {
		zlFields = append(zlFields, zap.Any(field.Key, field.Val))
	}
	// contextKeys
	for _, ctxKey := range l.config.CtxKeys {
		if val := ctx.Value(ctxKey); val != nil {
			zlFields = append(zlFields, zap.Any(ctxKey, val))
		}
	}
	// fields
	for _, field := range l.fields {
		zlFields = append(zlFields, zap.Any(field.Key, field.Val))
	}
	if len(zlFields) > 0 {
		zl = l.zl.With(zlFields...).Sugar()
	} else {
		zl = l.zl.Sugar()
	}
	return
}
