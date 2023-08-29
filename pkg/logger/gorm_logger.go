package logger

import (
	"context"
	"errors"
	"path/filepath"
	"runtime"
	"strings"
	"time"
	"yafgo/yafgo-layout/pkg/sys/ylog"

	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// GormLogger 操作对象，实现 gormlogger.Interface
type GormLogger struct {
	Logger        *ylog.Logger
	SlowThreshold time.Duration
}

// NewGormLogger 外部调用。实例化一个 GormLogger 对象，示例：
//
//	DB, err := gorm.Open(dbConfig, &gorm.Config{
//	    Logger: logger.NewGormLogger(),
//	})
func NewGormLogger(lg *ylog.Logger) GormLogger {
	return GormLogger{
		Logger:        lg,
		SlowThreshold: 200 * time.Millisecond, // 慢查询阈值，单位为千分之一秒
	}
}

// LogMode 实现 gormlogger.Interface 的 LogMode 方法
func (l GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return GormLogger{
		Logger:        l.Logger,
		SlowThreshold: l.SlowThreshold,
	}
}

// Info 实现 gormlogger.Interface 的 Info 方法
func (l GormLogger) Info(ctx context.Context, str string, args ...interface{}) {
	l.logger().Infof(ctx, str, args...)
}

// Warn 实现 gormlogger.Interface 的 Warn 方法
func (l GormLogger) Warn(ctx context.Context, str string, args ...interface{}) {
	l.logger().Warnf(ctx, str, args...)
}

// Error 实现 gormlogger.Interface 的 Error 方法
func (l GormLogger) Error(ctx context.Context, str string, args ...interface{}) {
	l.logger().Errorf(ctx, str, args...)
}

// Trace 实现 gormlogger.Interface 的 Trace 方法
func (l GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {

	// 获取运行时间
	elapsed := time.Since(begin)
	// 获取 SQL 请求和返回条数
	sql, rows := fc()

	// 通用字段
	lg := l.logger().With(
		ylog.Any("sql", sql),
		ylog.Any("elapsed", elapsed.String()),
		ylog.Any("rows", rows),
	)

	// Gorm 错误
	if err != nil {
		// 记录未找到的错误使用 warning 等级
		if errors.Is(err, gorm.ErrRecordNotFound) {
			lg.Warn(ctx, "Database ErrRecordNotFound")
		} else {
			// 其他错误使用 error 等级
			lg.Error(ctx, err)
		}
	}

	// 慢查询日志
	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		lg.Warn(ctx, "Database Slow Log")
	} else {
		// 记录所有 SQL 请求
		lg.Info(ctx, "Database Query")
	}

}

// logger 内用的辅助方法，确保 Zap 内置信息 Caller 的准确性（如 paginator/paginator.go:148）
func (l GormLogger) logger() ylog.ILogger {

	// 跳过 gorm 内置的调用
	var (
		pkgGorm    = filepath.Join("gorm.io", "gorm")
		pkgGormGen = filepath.Join("gorm.io", "gen")
	)

	// 减去一次封装
	clone := l.Logger.WithCallerSkip(-1)

	for i := 2; i < 10; i++ {
		_, file, _, ok := runtime.Caller(i)
		switch {
		case !ok:
		case strings.HasSuffix(file, "_test.go"):
		case strings.Contains(file, pkgGorm):
		case strings.Contains(file, pkgGormGen):
		default:
			// 返回一个附带跳过行号的新的 zap logger
			return clone.WithCallerSkip(i)
		}
	}

	return clone
}
