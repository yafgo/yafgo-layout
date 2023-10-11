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
	webHandler   handler.WebHandler
	indexHandler handler.IndexHandler
	userHandler  handler.UserHandler
}

func NewWebService(
	logger *ylog.Logger,
	cfg *ycfg.Config,
	g *g.GlobalObj,
	webHandler handler.WebHandler,
	indexHandler handler.IndexHandler,
	userHandler handler.UserHandler,
) *WebService {
	return &WebService{
		logger: logger,
		cfg:    cfg,
		g:      g,

		// handler
		webHandler:   webHandler,
		indexHandler: indexHandler,
		userHandler:  userHandler,
	}
}
