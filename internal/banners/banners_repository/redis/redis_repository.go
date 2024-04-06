package banners_redis

import (
	"avito/assignment/config"
	"github.com/redis/go-redis/v9"
)

type ClientRedisRepo struct {
	db  *redis.Client
	cfg *config.Config
}

func NewClientRedisRepository(db *redis.Client, cfg *config.Config) *ClientRedisRepo {
	return &ClientRedisRepo{
		db:  db,
		cfg: cfg,
	}
}
