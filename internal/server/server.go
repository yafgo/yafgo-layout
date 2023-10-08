package server

import (
	"yafgo/yafgo-layout/internal/handler"
	"yafgo/yafgo-layout/pkg/sys/ycfg"
	"yafgo/yafgo-layout/pkg/sys/ylog"
)

type WebService struct {
	logger *ylog.Logger
	cfg    *ycfg.Config

	// handler
	webHandler   handler.WebHandler
	indexHandler handler.IndexHandler
	userHandler  handler.UserHandler
}

func NewWebService(
	logger *ylog.Logger,
	cfg *ycfg.Config,
	webHandler handler.WebHandler,
	indexHandler handler.IndexHandler,
	userHandler handler.UserHandler,
) *WebService {
	return &WebService{
		logger: logger,
		cfg:    cfg,

		// handler
		webHandler:   webHandler,
		indexHandler: indexHandler,
		userHandler:  userHandler,
	}
}
