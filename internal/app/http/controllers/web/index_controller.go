package web

import (
	"github.com/gin-gonic/gin"
)

type IndexController struct {
	BaseWebController
}

// Index index
func (ctrl *IndexController) Index(c *gin.Context) {
	c.String(200, "Hello")
}
