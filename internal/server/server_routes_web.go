package server

import (
	"yafgo/yafgo-layout/internal/app/http/controllers/web"

	"github.com/gin-gonic/gin"
)

func (s *WebService) registerRoutesWeb(router *gin.Engine) {

	rWeb := router.Group("/")

	// index
	ic := new(web.IndexController)
	rWeb.GET("/index", ic.Index)
}
