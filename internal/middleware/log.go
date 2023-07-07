package middleware

import (
	"bytes"
	"fmt"
	"go-toy/toy-layout/pkg/hash"
	logger_pkg "go-toy/toy-layout/pkg/logger"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Logger 记录请求日志
func Logger() gin.HandlerFunc {
	logger := logger_pkg.Logger
	return func(ctx *gin.Context) {

		// trace_id 存到 ctx 中
		trace := hash.Md5(hash.GenUUID())
		logger.NewContext(ctx, zap.String("trace", trace))

		// request info
		reqLogFields := []zap.Field{
			zap.String("req_method", ctx.Request.Method),
			zap.String("req_url", ctx.Request.URL.String()),
			zap.String("query", ctx.Request.URL.RawQuery),
			zap.String("ip", ctx.ClientIP()),
			zap.String("ua", ctx.Request.UserAgent()),
		}
		headers := ctx.Request.Header.Clone()
		reqLogFields = append(reqLogFields, zap.Any("header", headers))
		var reqBody []byte
		if ctx.Request.Body != nil {
			// ctx.Request.Body 是一个 buffer 对象，只能读取一次
			reqBody, _ = ctx.GetRawData()
			// [重要] 读取后，重新赋值 ctx.Request.Body ，以供后续的其他操作
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))
			reqLogFields = append(reqLogFields, zap.String("req_params", string(reqBody)))
		}
		// logger.WithContext(ctx).With(reqLogFields...).Info("Request")

		// 记录耗时
		t1 := time.Now()
		ctx.Next()
		cost := time.Since(t1)

		// 获取 response 内容
		/* w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: ctx.Writer}
		ctx.Writer = w */

		// 记录响应状态和耗时
		respStatus := ctx.Writer.Status()
		logFields := []zap.Field{
			zap.Int("status", respStatus),
			zap.String("elapse", fmt.Sprintf("%v", cost)),
		}
		logFields = append(logFields, reqLogFields...)
		logFields = append(logFields, zap.String("errors", ctx.Errors.ByType(gin.ErrorTypePrivate).String()))

		// 记录本次请求log
		_logger := logger.WithContext(ctx).With(logFields...)
		logMsg := "HTTP Access Log"
		if respStatus > 400 && respStatus <= 499 {
			_logger.Warn(logMsg)
		} else if respStatus >= 500 && respStatus <= 599 {
			_logger.Error(logMsg)
		} else {
			_logger.Info(logMsg)
		}
	}
}

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}
