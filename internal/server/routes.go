package server

import (
	"go-toy/toy-layout/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func registerRoutes(r *gin.Engine) {
	// 注册全局中间件
	registerGlobalMiddleware(r)

	r.GET("", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"Hello": "Go Toy"})
		logger.Logger.WithContext(ctx).Debug("handler", zap.String("hello", "go toy"))
	})

	// 静态文件
	// r.StaticFile("/favicon.ico", "public/favicon.ico")
	// r.Static("/static", "public/static/")

	// 处理 404
	r.NoRoute(handle404)
}
