package web

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {

	// 静态目录
	router.Static("/static", "public/static/")

	rWeb := router.Group("/")

	// index
	ic := new(IndexController)
	rWeb.GET("/index", ic.Index)

}
