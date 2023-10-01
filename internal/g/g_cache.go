package g

import (
	"context"
	"sync"
	"yafgo/yafgo-layout/pkg/cache"
	"yafgo/yafgo-layout/pkg/sys/ylog"
)

// 全局 cache 对象
var instCache cache.Store
var onceCache sync.Once

// Cache 获取全局 Cache 对象
func Cache() cache.Store {
	onceCache.Do(func() {
		ctx := context.Background()
		subCfg := Cfg().Sub("redis.cache")
		if subCfg == nil {
			ylog.Fatal(ctx, "初始化cache出错, cache redis 配置不存在")
			return
		}

		_cache, err := cache.NewRedis(
			subCfg.GetString("addr"),
			subCfg.GetString("password"),
			subCfg.GetInt("db"),
			subCfg.GetString("prefix")+":",
		)
		if err != nil {
			ylog.Fatalf(ctx, "初始化cache出错: %v", err)
			return
		}
		instCache = _cache
	})

	return instCache
}
