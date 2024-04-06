package banners

import (
	"avito/assignment/internal/models/banner_models"
	"context"
)

type PostgresRepository interface {
	GetPossibleBannerIds(ctx context.Context, tagId int) ([]int, error)
	GetBanner(ctx context.Context, params banner_models.GetPostgresBanner) (*banner_models.BannerContent, error)
}

type RedisRepository interface {
	PutBanner(ctx context.Context, params banner_models.PutRedisBanner) error
	GetBanner(ctx context.Context, params banner_models.GetRedisBanner) (*banner_models.BannerContent, error)
}
