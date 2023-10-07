package app

import (
	"yafgo/yafgo-layout/internal/server"
	"yafgo/yafgo-layout/pkg/sys/ycfg"
	"yafgo/yafgo-layout/pkg/sys/ylog"
)

type application struct {
	cfg        *ycfg.Config
	logger     *ylog.Logger
	webService *server.WebService
}

func newApplication(
	logger *ylog.Logger,
	cfg *ycfg.Config,
	webService *server.WebService,
) (app *application) {
	return &application{
		cfg:        cfg,
		logger:     logger,
		webService: webService,
	}
}

// Mode app模式: dev, prod, ...
func (app *application) Mode() string {
	return app.cfg.GetString("env")
}
