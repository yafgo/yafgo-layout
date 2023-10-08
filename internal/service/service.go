package service

import (
	"yafgo/yafgo-layout/pkg/sys/ylog"
)

type Service struct {
	logger *ylog.Logger
	// jwt    *jwt.JWT
}

func NewService(
	logger *ylog.Logger,
	// jwt *jwt.JWT,
) *Service {
	return &Service{
		logger: logger,
		// jwt:    jwt,
	}
}
