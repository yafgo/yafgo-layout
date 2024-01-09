package frontend

import (
	"yafgo/yafgo-layout/internal/handler"
	"yafgo/yafgo-layout/pkg/sys/ylog"

	"github.com/gin-gonic/gin"
)

type Router struct {
	logger *ylog.Logger

	// handler
	hdl          *handler.Handler
	indexHandler IndexHandler
	webHandler   WebHandler
	userHandler  UserHandler
}

func NewRouter(
	logger *ylog.Logger,
	hdl *handler.Handler,
) *Router {
	return &Router{
		logger: logger,
		hdl:    hdl,

		// handler
		indexHandler: NewIndexHandler(hdl),
		webHandler:   NewWebHandler(hdl),
		userHandler:  NewUserHandler(hdl),
	}
}

func (p *Router) Register(router *gin.Engine) {

	// 前台web
	{
		r := router.Group("/")

		// root
		router.GET("", p.webHandler.Root)
		// index
		r.GET("/index", p.webHandler.Index)
	}

	rApi := router.Group("/api")

	// 前台接口 [v1]
	{
		r := rApi.Group("/v1")

		// index
		rApi.GET("", p.indexHandler.Root)
		r.GET("", p.indexHandler.Index)
		r.GET("todo", todo)

		// auth
		{
			r.POST("/auth/register/username", p.userHandler.RegisterByUsername)
			r.POST("/auth/login/username", p.userHandler.LoginByUsername)
		}
	}

}

func todo(c *gin.Context) {
	reqUri := c.Request.RequestURI
	c.JSON(200, gin.H{"todo": reqUri, "success": true, "code": 0})
}
