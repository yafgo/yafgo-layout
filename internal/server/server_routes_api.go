package server

import (
	"github.com/gin-gonic/gin"
)

func (s *WebService) registerRoutesApi(router *gin.Engine) {

	rApi := router.Group("/api")
	{
		rApi.GET("", s.indexHandler.Root)
	}

	// v1
	rV1 := rApi.Group("/v1")
	{
		// index
		{
			rV1.GET("", s.indexHandler.Index)
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
