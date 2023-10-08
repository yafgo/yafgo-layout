package server

import (
	"yafgo/yafgo-layout/internal/middleware"

	"github.com/gin-gonic/gin"
)

// NewGinEngine
func (s *WebService) NewGinEngine(isProd bool) *gin.Engine {
	// 设置 gin 的运行模式，支持 debug, release, test, 生产环境请使用 release 模式
	if isProd {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// 初始化 Gin 实例
	r := gin.New()

	// 注册全局中间件
	s.registerGlobalMiddleware(r)

	// 注册路由
	s.registerRoutes(r)

	return r
}

// registerGlobalMiddleware 注册全局中间件
func (s *WebService) registerGlobalMiddleware(router *gin.Engine) {
	router.Use(
		middleware.Logger(),
		gin.Recovery(),
		middleware.CORS(),
		/* func(ctx *gin.Context) {
			ctx.Header("Access-Control-Expose-Headers", "custom-header")
		}, */
	)
}
