package server

import (
	"yafgo/yafgo-layout/internal/g"
	"yafgo/yafgo-layout/internal/handler"
	"yafgo/yafgo-layout/pkg/sys/ycfg"
	"yafgo/yafgo-layout/pkg/sys/ylog"
)

type WebService struct {
	logger *ylog.Logger
	cfg    *ycfg.Config
	g      *g.GlobalObj

	// handler
	hdl          *handler.Handler
	webHandler   handler.WebHandler
	indexHandler handler.IndexHandler
	userHandler  handler.UserHandler
}

func NewWebService(
	logger *ylog.Logger,
	cfg *ycfg.Config,
	g *g.GlobalObj,
	hdl *handler.Handler,
) *WebService {
	return &WebService{
		logger: logger,
		cfg:    cfg,
		g:      g,

		// handler
		hdl:          hdl,
		webHandler:   handler.NewWebHandler(hdl),
		indexHandler: handler.NewIndexHandler(hdl),
		userHandler:  handler.NewUserHandler(hdl),
	}
}
