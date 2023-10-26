package middleware

import (
	"bytes"
	"fmt"
	"io"
	"time"
	"yafgo/yafgo-layout/pkg/hash"
	"yafgo/yafgo-layout/pkg/sys/ylog"

	"github.com/gin-gonic/gin"
)

// Logger 记录请求日志
func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// req_id 存到 ctx 中
		traceID := hash.GenGUID()
		ctx.Set("reqid", traceID)
		ctx.Header("X-Reqid", traceID)

		// request info
		headers := ctx.Request.Header.Clone()
		reqLogFields := []ylog.Field{
			ylog.Any("req_method", ctx.Request.Method),
			ylog.Any("req_url", ctx.Request.URL.String()),
			ylog.Any("query", ctx.Request.URL.RawQuery),
			ylog.Any("ip", ctx.ClientIP()),
			ylog.Any("ua", ctx.Request.UserAgent()),
			ylog.Any("header", headers),
		}
		var reqBody []byte
		if ctx.Request.Body != nil {
			// ctx.Request.Body 是一个 buffer 对象，只能读取一次
			reqBody, _ = ctx.GetRawData()
			// [重要] 读取后，重新赋值 ctx.Request.Body ，以供后续的其他操作
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))
			reqLogFields = append(reqLogFields, ylog.Any("req_params", string(reqBody)))
		}
		// 记录请求
		ylog.With(reqLogFields...).Info(ctx, "Request")

		// 记录耗时
		t1 := time.Now()
		ctx.Next()
		cost := time.Since(t1)

		// 获取 response 内容
		/* w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: ctx.Writer}
		ctx.Writer = w */

		// 记录响应状态和耗时
		respStatus := ctx.Writer.Status()
		reqLogFields = append(reqLogFields,
			ylog.Any("status", respStatus),
			ylog.Any("elapse", fmt.Sprintf("%v", cost)),
			ylog.Any("errors", ctx.Errors.ByType(gin.ErrorTypePrivate).String()),
		)

		// 记录本次请求log
		_ylog := ylog.With(reqLogFields...)
		logMsg := "HTTP Access Log"
		if respStatus > 400 && respStatus <= 499 {
			_ylog.Warn(ctx, logMsg)
		} else if respStatus >= 500 && respStatus <= 599 {
			_ylog.Error(ctx, logMsg)
		} else {
			_ylog.Info(ctx, logMsg)
		}
	}
}

/* type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
} */
