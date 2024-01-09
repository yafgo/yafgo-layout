package providers

import (
	"context"
	"time"
	"yafgo/yafgo-layout/internal/query"
	"yafgo/yafgo-layout/pkg/cache"
	"yafgo/yafgo-layout/pkg/database"
	"yafgo/yafgo-layout/pkg/logger"
	"yafgo/yafgo-layout/pkg/sys/ycfg"
	"yafgo/yafgo-layout/pkg/sys/ylog"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func NewDB(cfg *ycfg.Config, lg *ylog.Logger) *gorm.DB {
	gormLogger := logger.NewGormLogger(lg)
	db, err := database.NewGormMysql(cfg, gormLogger)
	if err != nil {
		panic(err)
	}
	return db
}

func NewGormQuery(db *gorm.DB) *query.Query {
	query.SetDefault(db)
	return query.Q
}

func NewRedis(cfg *ycfg.Config) *redis.Client {
	subCfg := cfg.Sub("redis.default")
	if subCfg == nil {
		panic("初始化redis出错, redis.default 配置不存在")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     subCfg.GetString("addr"),
		Password: subCfg.GetString("password"),
		DB:       subCfg.GetInt("db"),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic("redis error: " + err.Error())
	}
	return rdb
}

func NewCache(cfg *ycfg.Config) cache.Store {
	subCfg := cfg.Sub("redis.cache")
	if subCfg == nil {
		panic("初始化cache出错, redis.cache 配置不存在")
	}

	cc, err := cache.NewRedis(
		subCfg.GetString("addr"),
		subCfg.GetString("password"),
		subCfg.GetInt("db"),
		subCfg.GetString("prefix")+":",
	)

	if err != nil {
		panic("初始化cache出错: " + err.Error())
	}
	return cc
}
