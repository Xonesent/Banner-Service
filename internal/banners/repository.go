package banners

import (
	"avito/assignment/internal/models/banner_models"
	"context"
)

type PostgresRepository interface {
	GetPossibleBannerIds(ctx context.Context, tagId int) ([]int, error)
	GetPossibleTagIds(ctx context.Context, bannerId int) ([]int, error)
	GetBanner(ctx context.Context, params banner_models.GetPostgresBanner) (*banner_models.BannerContent, error)
	SelectBanner(ctx context.Context, params banner_models.SelectPostgresBanner) (*[]banner_models.FullBannerContent, error)
	AddBanner(ctx context.Context, params banner_models.AddBanner) (*int, error)
	AddTags(ctx context.Context, bannerId int, tagIds int) error
}

type RedisRepository interface {
	PutBanner(ctx context.Context, params banner_models.PutRedisBanner) error
	GetBanner(ctx context.Context, params banner_models.GetRedisBanner) (*banner_models.BannerContent, error)
}
