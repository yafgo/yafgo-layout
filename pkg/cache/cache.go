package cache

import (
	"context"
	"time"
)

type Store interface {
	Get(ctx context.Context, key string) (any, error)
	GetWithTTL(ctx context.Context, key string) (any, time.Duration, error)
	Set(ctx context.Context, key string, value any, ttl ...time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) bool
	Clear(ctx context.Context) error
	GetType() string
}
