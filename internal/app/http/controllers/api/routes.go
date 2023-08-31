package api

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {

	rApi := router.Group("/api")

	// index
	ic := new(IndexController)
	{
		rApi.GET("", ic.Index)
		rApi.GET("todo", todo)
	}
}

func todo(c *gin.Context) {
	c.JSON(200, gin.H{"todo": "hello world"})
}
