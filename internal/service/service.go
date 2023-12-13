package service

import (
	"yafgo/yafgo-layout/internal/query"
	"yafgo/yafgo-layout/pkg/jwtutil"
	"yafgo/yafgo-layout/pkg/sys/ylog"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Service struct {
	Logger *ylog.Logger
	Jwt    *jwtutil.JwtUtil
	DB     *gorm.DB
	Redis  *redis.Client
	Q      *query.Query
}

func NewService(
	logger *ylog.Logger,
	jwt *jwtutil.JwtUtil,
	db *gorm.DB,
	rdb *redis.Client,
	q *query.Query,
) *Service {
	return &Service{
		Logger: logger,
		Jwt:    jwt,
		DB:     db,
		Redis:  rdb,
		Q:      q,
	}
}
