package handler

import (
	"github.com/gin-gonic/gin"
)

type IndexHandler interface {
	Root(ctx *gin.Context)
	Index(ctx *gin.Context)
}

type indexHandler struct {
	*Handler
}

func NewIndexHandler(
	handler *Handler,
) IndexHandler {
	return &indexHandler{
		Handler: handler,
	}
}

// Root implements IndexHandler.
//
//	@Summary		ApiRoot
//	@Description	Api Root
//	@Tags			API
//	@Success		200	{object}	any	"{"code": 200, "data": [...]}"
//	@Router			/ [get]
//	@Security		ApiToken
func (h *indexHandler) Root(ctx *gin.Context) {
	h.Resp().Success(ctx, gin.H{
		"data": "/api/",
	})
}

// Index implements IndexHandler.
//
//	@Summary		ApiIndex
//	@Description	Api Index Demo
//	@Tags			API
//	@Success		200	{object}	any	"{"code": 200, "data": [...]}"
//	@Router			/v1/ [get]
//	@Security		ApiToken
func (h *indexHandler) Index(ctx *gin.Context) {
	h.Resp().Success(ctx, gin.H{
		"data": "/api/v1/",
	})
}
