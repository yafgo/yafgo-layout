package server

import (
	"github.com/gin-gonic/gin"
)

func (s *WebService) registerRoutesWeb(router *gin.Engine) {

	rWeb := router.Group("/")

	// index
	rWeb.GET("/index", s.webHandler.Index)
}
