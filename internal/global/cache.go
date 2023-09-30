package global

import (
	"context"
	"yafgo/yafgo-layout/pkg/cache"
	"yafgo/yafgo-layout/pkg/sys/ycfg"
	"yafgo/yafgo-layout/pkg/sys/ylog"
)

var Cache cache.Store

func SetupCache(ctx context.Context, cfg *ycfg.Config) {
	subCfg := cfg.Sub("redis.cache")
	if subCfg == nil {
		ylog.Errorf(ctx, "初始化cache出错, cache redis 配置不存在")
		return
	}

	_cache, err := cache.NewRedis(
		subCfg.GetString("addr"),
		subCfg.GetString("password"),
		subCfg.GetInt("db"),
		subCfg.GetString("prefix")+":",
	)
	if err != nil {
		ylog.Errorf(ctx, "初始化cache出错: %v", err)
		return
	}
	Cache = _cache
}
