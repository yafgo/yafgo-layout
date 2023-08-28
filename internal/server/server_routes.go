package server

import (
	"go-toy/toy-layout/internal/app/http/controllers/api"
	"go-toy/toy-layout/internal/app/http/controllers/web"
	"go-toy/toy-layout/internal/global"

	"github.com/gin-gonic/gin"
)

func registerRoutes(router *gin.Engine) {
	// 注册全局中间件
	registerGlobalMiddleware(router)

	router.GET("", func(ctx *gin.Context) {
		appname := global.AppName()
		ctx.JSON(200, gin.H{"Hello": appname})
	})

	// 静态文件
	// r.StaticFile("/favicon.ico", "public/favicon.ico")
	// r.Static("/static", "public/static/")

	// api 路由
	api.RegisterRoutes(router)

	// web 路由
	web.RegisterRoutes(router)

	// 处理 404
	router.NoRoute(handle404)
}
