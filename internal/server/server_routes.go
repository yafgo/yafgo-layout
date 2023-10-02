package server

import (
	"yafgo/yafgo-layout/internal/g"
	"yafgo/yafgo-layout/internal/middleware"
	"yafgo/yafgo-layout/resource/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
)

func registerRoutes(router *gin.Engine) {
	// 注册全局中间件
	registerGlobalMiddleware(router)

	router.GET("", func(ctx *gin.Context) {
		appname := g.AppName()
		ctx.JSON(200, gin.H{"Hello": appname})
	})

	// 静态文件
	router.StaticFile("/favicon.ico", "resource/public/favicon.ico")
	router.Static("/static", "public/static/")

	// api 路由
	registerRoutesApi(router)

	// web 路由
	registerRoutesWeb(router)

	// swagger
	handleSwagger(router)

	// 处理 404
	router.NoRoute(handle404)
}

// handleSwagger 启用 swagger
func handleSwagger(router *gin.Engine) {
	apiGroup := router.Group("/api/docs")

	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	if !g.IsDev() {
		docs.SwaggerInfo.Schemes = []string{"https", "http"}
		// 非开发环境启用 BasicAuth 验证
		apiGroup.Use(middleware.BasicAuth("swagger"))
	}
	apiGroup.GET("/*any", ginswagger.WrapHandler(
		swaggerfiles.Handler,
		ginswagger.PersistAuthorization(true),
	))
}
