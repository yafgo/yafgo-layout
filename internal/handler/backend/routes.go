package backend

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
	}
}

func (p *Router) Register(router *gin.Engine) {

	rApi := router.Group("/api/backend")

	{
		r := rApi.Group("/v1")

		// index
		r.GET("/index", p.indexHandler.Index)
	}
}
