package api

import "github.com/gin-gonic/gin"

type IndexController struct {
	BaseAPIController
}

// Index index
func (ctrl *IndexController) Index(c *gin.Context) {

	ctrl.JSON(c, gin.H{
		"data": "data",
	})
}
