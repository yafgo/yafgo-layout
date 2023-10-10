package app

import (
	"yafgo/yafgo-layout/internal/play"
	"yafgo/yafgo-layout/internal/server"
	"yafgo/yafgo-layout/pkg/sys/ycfg"
	"yafgo/yafgo-layout/pkg/sys/ylog"
)

type application struct {
	cfg        *ycfg.Config
	logger     *ylog.Logger
	webService *server.WebService
	playground *play.Playground
}

func newApplication(
	logger *ylog.Logger,
	cfg *ycfg.Config,
	webService *server.WebService,
	playground *play.Playground,
) (app *application) {
	return &application{
		cfg:        cfg,
		logger:     logger,
		webService: webService,
		playground: playground,
	}
}
