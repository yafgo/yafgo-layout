package server

import (
	"go-toy/toy-layout/internal/middleware"
	"strings"

	"github.com/gin-gonic/gin"
)

// NewGinEngine
func NewGinEngine(isProd bool) *gin.Engine {
	// 设置 gin 的运行模式，支持 debug, release, test, 生产环境请使用 release 模式
	if isProd {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// 初始化 Gin 实例
	r := gin.New()

	// 注册路由
	registerRoutes(r)

	return r
}

// registerGlobalMiddleware 注册全局中间件
func registerGlobalMiddleware(router *gin.Engine) {
	router.Use(
		gin.Logger(),
		gin.Recovery(),
		middleware.CORS(),
		/* func(ctx *gin.Context) {
			ctx.Header("Access-Control-Expose-Headers", "custom-header")
		}, */
	)
}

func handle404(c *gin.Context) {
	acceptString := c.Request.Header.Get("Accept")
	if strings.Contains(acceptString, "text/html") {
		// c.String(404, "404")
	} else {
		c.JSON(404, gin.H{
			"code":    404,
			"message": "路由未定义",
		})
	}
}
