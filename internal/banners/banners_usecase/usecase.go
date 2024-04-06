package banners_usecase

import (
	"avito/assignment/config"
	banners_postgres "avito/assignment/internal/banners/banners_repository/postgres"
	banners_redis "avito/assignment/internal/banners/banners_repository/redis"
	"github.com/avito-tech/go-transaction-manager/trm/manager"
)

type BannersUC struct {
	cfg         *config.Config
	trManager   *manager.Manager
	bannersRepo *banners_postgres.BannersRepo
	redisClient *banners_redis.ClientRedisRepo
}

func NewBannersUC(cfg *config.Config, trManager *manager.Manager, bannersRepo *banners_postgres.BannersRepo, redisClient *banners_redis.ClientRedisRepo) *BannersUC {
	return &BannersUC{
		cfg:         cfg,
		trManager:   trManager,
		bannersRepo: bannersRepo,
		redisClient: redisClient,
	}
}
