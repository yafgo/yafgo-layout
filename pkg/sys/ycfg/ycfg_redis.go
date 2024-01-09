package ycfg

import (
	"context"
	"log"
	"strings"

	"github.com/gookit/color"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

// EnableRedis 启用 redis 支持
func (p *Config) EnableRedis(ctx context.Context, rdb *redis.Client, key ...string) (err error) {
	var _key string
	if len(key) > 0 {
		_key = key[0]
	}
	if _key == "" {
		_key = "key"
	}
	log.Println(color.Success.Sprint("redis配置存储已启用"))

	p.redisStore = NewRedisStore(rdb, _key)
	// 从redis读取初始配置
	content, _ := p.redisStore.Get(ctx)
	if _err := p.MergeConfig(strings.NewReader(content)); _err != nil {
		log.Println(color.Error.Sprintf("从redis读取配置出错: %+v", _err))
	}
	// 监听配置变化
	go p.redisStore.Watch(ctx, func(newContent string) {
		log.Println(color.Success.Sprint("redis中的配置改变"))
		if _err := p.MergeConfig(strings.NewReader(newContent)); _err != nil {
			log.Println(color.Error.Sprintf("从redis更新配置出错: %+v", _err))
		}
	})
	return
}

// GetRedisContent 获取redis中的配置
func (p *Config) GetRedisContent(ctx context.Context) (string, error) {
	return p.redisStore.Get(ctx)
}

// SetRedisContent 更新redis中的配置
func (p *Config) SetRedisContent(ctx context.Context, content string) error {
	return p.redisStore.Set(ctx, content)
}

func NewRedisStore(rdb *redis.Client, key string) *redisStore {
	inst := &redisStore{
		rdb: rdb,
		key: key,
	}

	return inst
}

type redisStore struct {
	rdb *redis.Client
	key string
}

// keyData 存储配置数据的key
func (rs *redisStore) keyData() string {
	return "ycfg:data:" + rs.key
}

// keyPubsub 订阅配置变动的key
func (rs *redisStore) keyPubsub() string {
	return "ycfg:pubsub:" + rs.key
}

func (rs *redisStore) Get(ctx context.Context) (string, error) {
	content, err := rs.rdb.Get(ctx, rs.keyData()).Result()
	if err == redis.Nil {
		return "", errors.Wrap(err, "notFound")
	}

	return content, err
}

func (rs *redisStore) Set(ctx context.Context, content string) error {
	err := rs.rdb.Set(ctx, rs.keyData(), content, 0).Err()
	if err != nil {
		return err
	}
	err = rs.rdb.Publish(ctx, rs.keyPubsub(), 1).Err()
	return err
}

func (rs *redisStore) Watch(ctx context.Context, cb func(newContent string)) {
	sub := rs.rdb.Subscribe(ctx, rs.keyPubsub())
	defer sub.Close()
	ch := sub.Channel()

	for {
		select {
		case <-ctx.Done():
			log.Println("redis取消订阅")
			return
		case <-ch:
			if cb != nil {
				val, _ := rs.Get(ctx)
				cb(val)
			}
		}
	}
}
