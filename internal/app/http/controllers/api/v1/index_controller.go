package v1

import "github.com/gin-gonic/gin"

type IndexController struct {
	BaseAPIController
}

// Index index
//
//	@Summary		Index
//	@Description	Index Demo
//	@Tags			默认
//	@Success		200	{object}	any	"{"code": 200, "data": [...]}"
//	@Router			/v1/ [get]
//	@Security		ApiToken
func (ctrl *IndexController) Index(c *gin.Context) {

	ctrl.JSON(c, gin.H{
		"data": "/api/v1/",
	})
}
