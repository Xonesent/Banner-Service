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
	AddBannerPostgres(ctx context.Context, addPostgresBannerParams *banners_repository.AddPostgresBanner) (*banners_repository.GetInsertParams, error)
	AddTags(ctx context.Context, addTagsPostgresParams *banners_repository.AddTagsPostgres) error
	CheckExist(ctx context.Context, checkExistBannerParams *banners_repository.CheckExistBanner) error
}

type RedisRepository interface {
	PutBannerRedis(ctx context.Context, putRedisBannerParams *banners_repository.PutRedisBanner) error
	GetBannerRedis(ctx context.Context, getRedisParams *banners_repository.GetRedisBanner) (*models.FullBanner, error)
}
