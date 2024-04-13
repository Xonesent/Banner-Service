package banners_usecase

import (
	"avito/assignment/internal/banners/banners_repository"
	"avito/assignment/internal/models"
	"context"
)

type PostgresRepository interface {
	GetPossibleTagIds(ctx context.Context, bannerId models.BannerId) ([]models.TagId, error)
	GetManyPossibleTagIds(ctx context.Context, bannerIds []models.BannerId, manyBanner *[]models.Banner) (*[]models.FullBanner, error)
	CheckExist(ctx context.Context, tagIds []models.TagId, featureId models.FeatureId) (*[]banners_repository.ExistBanner, error)

	GetBanner(ctx context.Context, featureId models.FeatureId, tagId models.TagId) (*models.Banner, error)
	GetBannerById(ctx context.Context, bannerId models.BannerId) (*models.FullBanner, error)
	GetManyBanner(ctx context.Context, getManyPostgresBannerParams *banners_repository.GetManyPostgresBanner) (*[]models.Banner, error)
	GetBannerVersions(ctx context.Context, bannerId models.BannerId, versions []int64) (*[]models.FullBanner, error)

	AddBanner(ctx context.Context, addPostgresBannerParams *banners_repository.AddPostgresBanner) (*banners_repository.GetInsertParams, error)
	AddTags(ctx context.Context, tagIds []models.TagId, bannerId models.BannerId) error
	AddVersion(ctx context.Context, prevBanner *models.FullBanner) error

	UpdateBannerById(ctx context.Context, bannerId *banners_repository.UpdateBannerById) error

	DeleteBannerById(ctx context.Context, bannerId models.BannerId) error
	DeleteTags(ctx context.Context, tagIds []models.TagId, bannerId models.BannerId) error
	DeleteVersion(ctx context.Context, versions []int64, bannerId models.BannerId) error
}

type RedisRepository interface {
	PutBannerRedis(ctx context.Context, putRedisBannerParams *banners_repository.PutRedisBanner) error
	GetBannerRedis(ctx context.Context, getRedisParams *banners_repository.GetRedisBanner) (*models.FullBanner, error)
	DelBannerRedis(ctx context.Context, getRedisParams *banners_repository.GetRedisBanner) error
}
