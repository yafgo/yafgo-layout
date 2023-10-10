//go:build wireinject
// +build wireinject

package main

import (
	"path/filepath"
	"yafgo/yafgo-layout/internal/app"
	"yafgo/yafgo-layout/internal/play"
	"yafgo/yafgo-layout/pkg/notify"
	"yafgo/yafgo-layout/pkg/sys/ycfg"

	"github.com/google/wire"
)

var playgroundSet = wire.NewSet(play.NewPlayground)

var notifySet = wire.NewSet(
	notify.NewFeishu,
)

var dbSet = wire.NewSet(
	app.NewRedis,
	app.NewCache,
	app.NewDB,
	app.NewGormQuery,
)

var yLogSet = wire.NewSet(app.NewYLog)

var yCfgSet = wire.NewSet(newYCfg)

func newPlay(envConf string) (pg *play.Playground) {
	panic(wire.Build(
		playgroundSet,
		notifySet,
		dbSet,
		yLogSet,
		yCfgSet,
	))
}

func newYCfg(envConf string) (cfg *ycfg.Config) {
	cfg = ycfg.New(envConf,
		ycfg.WithType("yaml"),
		ycfg.WithEnvPrefix("YAFGO"),
		ycfg.WithDir(filepath.Join("../../")),
		ycfg.WithDir(filepath.Join("../../config/")),
	)
	return
}
