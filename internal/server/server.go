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
	indexHandler handler.IndexHandler
	userHandler  handler.UserHandler
}

func NewWebService(
	logger *ylog.Logger,
	cfg *ycfg.Config,
	indexHandler handler.IndexHandler,
	userHandler handler.UserHandler,
) *WebService {
	return &WebService{
		logger: logger,
		cfg:    cfg,

		// handler
		indexHandler: indexHandler,
		userHandler:  userHandler,
	}
}
