package server

import (
	"yafgo/yafgo-layout/internal/g"
	"yafgo/yafgo-layout/internal/handler"
	"yafgo/yafgo-layout/internal/handler/backend"
	"yafgo/yafgo-layout/internal/handler/frontend"
	"yafgo/yafgo-layout/pkg/sys/ycfg"
	"yafgo/yafgo-layout/pkg/sys/ylog"
)

type WebService struct {
	logger *ylog.Logger
	cfg    *ycfg.Config
	g      *g.GlobalObj

	// 路由
	routerBackend  *backend.Router
	routerFrontend *frontend.Router
}

func NewWebService(
	logger *ylog.Logger,
	cfg *ycfg.Config,
	g *g.GlobalObj,
	hdl *handler.Handler,

	routerBackend *backend.Router,
	routerFrontend *frontend.Router,
) *WebService {
	return &WebService{
		logger: logger,
		cfg:    cfg,
		g:      g,

		// 路由
		routerBackend:  routerBackend,
		routerFrontend: routerFrontend,
	}
}
