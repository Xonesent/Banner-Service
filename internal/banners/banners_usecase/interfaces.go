package banners_usecase

import (
	"avito/assignment/internal/banners/banners_repository"
	"avito/assignment/internal/models"
	"context"
)

type PostgresRepository interface {
	GetPossibleBannerIds(ctx context.Context, tagId models.TagId) ([]models.BannerId, error)
	GetPossibleTagIds(ctx context.Context, bannerId models.BannerId) ([]models.TagId, error)

	GetBannerPostgres(ctx context.Context, getPostgresqlBannerParams *banners_repository.GetPostgresBanner) (*models.Banner, error)
	GetManyBannerPostgres(ctx context.Context, getManyPostgresBannerParams *banners_repository.GetManyPostgresBanner) (*[]models.Banner, error)
	GetManyPossibleTagIds(ctx context.Context, bannerIds []models.BannerId, manyBanner *[]models.Banner) (*[]models.FullBanner, error)

	CheckExist(ctx context.Context, checkExistBannerParams *banners_repository.CheckExistBanner) (*[]banners_repository.ExistBanner, error)
	AddBannerPostgres(ctx context.Context, addPostgresBannerParams *banners_repository.AddPostgresBanner) (*banners_repository.GetInsertParams, error)
	AddTags(ctx context.Context, addTagsPostgresParams *banners_repository.AddTagsPostgres) error

	GetBannerById(ctx context.Context, bannerId models.BannerId) (*models.FullBanner, error)
	UpdateBannerById(ctx context.Context, bannerId *banners_repository.UpdateBannerById) error
	AddVersion(ctx context.Context, prevBanner *models.FullBanner) error
	DeleteTags(ctx context.Context, deleteTagsPostgresParams *banners_repository.DeleteTagsPostgres) error
	DeleteVersion(ctx context.Context, deleteVersionPostgresParams *banners_repository.DeleteVersionPostgres) error

	DeleteBannerById(ctx context.Context, bannerId models.BannerId) error

	GetBannerVersions(ctx context.Context, bannerId models.BannerId, versions []int64) (*[]models.FullBanner, error)
}

type RedisRepository interface {
	PutBannerRedis(ctx context.Context, putRedisBannerParams *banners_repository.PutRedisBanner) error
	GetBannerRedis(ctx context.Context, getRedisParams *banners_repository.GetRedisBanner) (*models.FullBanner, error)
	DelBannerRedis(ctx context.Context, getRedisParams *banners_repository.GetRedisBanner) error
}
