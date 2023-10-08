package repository

import (
	"yafgo/yafgo-layout/internal/query"
	"yafgo/yafgo-layout/pkg/sys/ylog"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Repository struct {
	db     *gorm.DB
	rdb    *redis.Client
	q      *query.Query
	logger *ylog.Logger
}

func NewRepository(
	db *gorm.DB,
	rdb *redis.Client,
	q *query.Query,
	logger *ylog.Logger,
) *Repository {
	return &Repository{
		db:     db,
		rdb:    rdb,
		q:      q,
		logger: logger,
	}
}
