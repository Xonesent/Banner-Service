package banners_usecase

import (
	"avito/assignment/config"
	banners_postgres "avito/assignment/internal/banners/banners_repository/postgres"
	banners_redis "avito/assignment/internal/banners/banners_repository/redis"
	"avito/assignment/internal/models/banner_models"
	"avito/assignment/pkg/constant"
	"context"
	"errors"
	"github.com/avito-tech/go-transaction-manager/trm/manager"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
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

func (b *BannersUC) GetBanner(ctx context.Context, params banner_models.GetBanner) (*banner_models.BannerContent, error) {
	ctx, span := otel.Tracer("").Start(ctx, "BannersUC.GetBanner")
	defer span.End()

	bannerInfo := &banner_models.BannerContent{}
	if !params.UseLastVersion {
		bannerInfo, err := b.redisClient.GetBanner(ctx, banner_models.GetRedisBanner{
			TagId:     params.TagId,
			FeatureId: params.FeatureId,
		})
		if err != nil && !errors.Is(err, fiber.ErrNotFound) {
			return nil, err
		}
		if bannerInfo != nil {
			if bannerInfo.IsActive == false && params.AuthToken == constant.UserToken {
				return nil, fiber.ErrNotFound
			}
			return bannerInfo, nil
		}
	}

	err := b.trManager.Do(ctx, func(ctx context.Context) error {
		possibleBannerIds, err := b.bannersRepo.GetPossibleBannerIds(ctx, params.TagId)
		if err != nil {
			return err
		}
		bannerInfo, err = b.bannersRepo.GetBanner(ctx, banner_models.GetPostgresBanner{
			FeatureId:         params.FeatureId,
			PossibleBannerIds: possibleBannerIds,
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	if !params.UseLastVersion {
		if err := b.redisClient.PutBanner(ctx, banner_models.PutRedisBanner{
			TagIds:    []int{params.TagId},
			FeatureId: params.FeatureId,
			Content:   *bannerInfo,
		}); err != nil {
			return nil, err
		}
	}

	if bannerInfo.IsActive == false && params.AuthToken == constant.UserToken {
		return nil, fiber.ErrNotFound
	}
	return bannerInfo, nil
}
