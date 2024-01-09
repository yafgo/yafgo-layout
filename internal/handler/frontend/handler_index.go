package frontend

import (
	"yafgo/yafgo-layout/internal/handler"

	"github.com/gin-gonic/gin"
)

type IndexHandler interface {
	Root(ctx *gin.Context)
	Index(ctx *gin.Context)
}

func NewIndexHandler(
	handler *handler.Handler,
) IndexHandler {
	return &indexHandler{
		Handler: handler,
	}
}

type indexHandler struct {
	*handler.Handler
}

// Root implements WebHandler.
func (h *indexHandler) Root(ctx *gin.Context) {
	ctx.String(200, "Yafgo")
}

// Index implements WebHandler.
func (h *indexHandler) Index(ctx *gin.Context) {
	ctx.String(200, "Yafgo Index")
}
