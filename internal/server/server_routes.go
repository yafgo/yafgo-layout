package server

import (
	"strings"
	"yafgo/yafgo-layout/internal/middleware"
	"yafgo/yafgo-layout/resource/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
)

func (s *WebService) registerRoutes(router *gin.Engine) {

	// 静态文件
	router.StaticFile("/favicon.ico", "resource/public/favicon.ico")
	router.Static("/static", "public/static/")

	// 前台路由
	s.routerFrontend.Register(router)
	// 后台路由
	s.routerBackend.Register(router)

	// swagger
	s.handleSwagger(router)

	// 处理 404
	router.NoRoute(handle404)
}

// handleSwagger 启用 swagger
func (s *WebService) handleSwagger(router *gin.Engine) {
	apiGroup := router.Group("/api/docs")

	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	if !s.g.IsDev() {
		docs.SwaggerInfo.Schemes = []string{"https", "http"}
		// 非开发环境启用 BasicAuth 验证
		apiGroup.Use(middleware.BasicAuth(s.cfg)("swagger"))
	}
	apiGroup.GET("/*any", ginswagger.WrapHandler(
		swaggerfiles.Handler,
		ginswagger.PersistAuthorization(true),
	))
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
