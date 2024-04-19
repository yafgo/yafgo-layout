package frontend

import (
	"yafgo/yafgo-layout/internal/handler"
	"yafgo/yafgo-layout/pkg/yview"

	"github.com/gin-gonic/gin"
)

type WebHandler interface {
	Root(ctx *gin.Context)
	Index(ctx *gin.Context)
}

type webHandler struct {
	*handler.Handler

	staticDir string
}

func NewWebHandler(
	handler *handler.Handler,
) WebHandler {
	return &webHandler{
		Handler: handler,

		staticDir: "/static",
	}
}

// Root implements WebHandler.
func (h *webHandler) Root(ctx *gin.Context) {
	h.Index(ctx)
}

// Index implements WebHandler.
func (h *webHandler) Index(ctx *gin.Context) {
	// 绑定模板数据
	var data = map[string]any{
		"StaticDir":  h.staticDir,
		"StaticHash": "this-is-hash",
		"IsDev":      h.G.IsDev(),
		"appName":    h.G.AppName(),
	}

	viewTpl := yview.ViewTpl{
		Files: []string{"./resource/views/index.tmpl"},
		Name:  "layout.tmpl",
		Data:  data,
	}
	if err := yview.HandleView(ctx, viewTpl); err != nil {
		h.Logger.Warnf(ctx, "页面渲染出错 %#v", err)
		ctx.String(400, "出错了: %v", err)
	}
}
