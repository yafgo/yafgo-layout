package handler

import (
	"yafgo/yafgo-layout/pkg/sys/ylog"
)

type Handler struct {
	logger *ylog.Logger
}

func NewHandler(logger *ylog.Logger) *Handler {
	return &Handler{
		logger: logger,
	}
}
