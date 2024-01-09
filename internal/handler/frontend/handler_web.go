package frontend

import (
	"yafgo/yafgo-layout/internal/handler"

	"github.com/gin-gonic/gin"
)

type WebHandler interface {
	Root(ctx *gin.Context)
	Index(ctx *gin.Context)
}

type webHandler struct {
	*handler.Handler
}

func NewWebHandler(
	handler *handler.Handler,
) WebHandler {
	return &webHandler{
		Handler: handler,
	}
}

// Root implements WebHandler.
func (h *webHandler) Root(ctx *gin.Context) {
	ctx.String(200, "Yafgo")
}

// Index implements WebHandler.
func (h *webHandler) Index(ctx *gin.Context) {
	ctx.String(200, "Yafgo Index")
}
