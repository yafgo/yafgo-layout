package server

import (
	"yafgo/yafgo-layout/internal/app/http/controllers/api"
	v1 "yafgo/yafgo-layout/internal/app/http/controllers/api/v1"

	"github.com/gin-gonic/gin"
)

func (s *WebService) registerRoutesApi(router *gin.Engine) {

	rApi := router.Group("/api")

	ic := new(api.IndexController)
	{
		rApi.GET("", ic.Index)
	}

	// v1
	rV1 := rApi.Group("/v1")
	{
		// index
		ic := new(v1.IndexController)
		{
			rV1.GET("", ic.Index)
			rV1.GET("todo", todo)
		}

		// auth
		{
			rV1.POST("/auth/register/username", s.userHandler.RegisterByUsername)
			rV1.POST("/auth/login/username", s.userHandler.LoginByUsername)
		}
	}
}

func todo(c *gin.Context) {
	reqUri := c.Request.RequestURI
	c.JSON(200, gin.H{"todo": reqUri})
}
