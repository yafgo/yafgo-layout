package server

import (
	"go-toy/toy-layout/global"
	"go-toy/toy-layout/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func registerRoutes(r *gin.Engine) {
	// 注册全局中间件
	registerGlobalMiddleware(r)

	r.GET("", func(ctx *gin.Context) {
		appname := global.AppName()
		ctx.JSON(200, gin.H{"Hello": appname})
		logger.Logger.WithContext(ctx).Debug("handler", zap.String("hello", appname))
	})

	// 静态文件
	// r.StaticFile("/favicon.ico", "public/favicon.ico")
	// r.Static("/static", "public/static/")

	// 处理 404
	r.NoRoute(handle404)
}
