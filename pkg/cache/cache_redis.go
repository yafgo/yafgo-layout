package cache

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type storeRedis struct {
	prefix string
	client *redis.Client
}

func NewRedis(addr string, password string, database int, prefix string) (*storeRedis, error) {

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       database,
	})

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return nil, errors.WithMessage(err, "init connection error")
	}

	return &storeRedis{
		prefix: prefix,
		client: client,
	}, nil
}

func (s *storeRedis) Get(ctx context.Context, key string) (any, error) {
	key = s.prefix + key
	val, err := s.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, errors.Wrap(err, "notFound")
	}

	return val, err
}

func (s *storeRedis) GetWithTTL(ctx context.Context, key string) (any, time.Duration, error) {
	key = s.prefix + key
	val, err := s.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, 0, errors.Wrap(err, "notFound")
	}
	if err != nil {
		return "", 0, err
	}

	ttl, err := s.client.TTL(ctx, key).Result()
	if err != nil {
		return nil, 0, err
	}

	return val, ttl, err
}

func (s *storeRedis) Set(ctx context.Context, key string, value any, ttl ...time.Duration) error {
	var _ttl time.Duration = -1
	if len(ttl) > 0 {
		_ttl = ttl[0]
	}
	err := s.client.Set(ctx, s.prefix+key, value, _ttl).Err()
	if err != nil {
		return err
	}

	return nil
}

func (s *storeRedis) Delete(ctx context.Context, key string) error {
	_, err := s.client.Del(ctx, s.prefix+key).Result()

	return err
}

func (s *storeRedis) Exists(ctx context.Context, key string) bool {
	value, err := s.client.Exists(ctx, s.prefix+key).Result()

	if err != nil || value == 0 {
		return false
	}

	return true
}

func (s *storeRedis) Clear(ctx context.Context) error {
	if err := s.client.FlushAll(ctx).Err(); err != nil {
		return err
	}

	return nil
}

func (s *storeRedis) GetType() string {
	return "redis"
}
