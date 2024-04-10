package banners_usecase

import (
	"avito/assignment/config"
	"avito/assignment/internal/banners/banners_repository"
	"avito/assignment/internal/models"
	"avito/assignment/pkg/constant"
	"context"
	"errors"
	"fmt"
	"github.com/avito-tech/go-transaction-manager/trm/manager"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
)

type BannersUC struct {
	cfg              *config.Config
	trManager        *manager.Manager
	bannersPGRepo    PostgresRepository
	bannersRedisRepo RedisRepository
}

func NewBannersUC(cfg *config.Config, trManager *manager.Manager, bannersRepo PostgresRepository, redisClient RedisRepository) *BannersUC {
	return &BannersUC{
		cfg:              cfg,
		trManager:        trManager,
		bannersPGRepo:    bannersRepo,
		bannersRedisRepo: redisClient,
	}
}

func (b *BannersUC) GetBanner(ctx context.Context, getBannerParams *GetBanner) (*models.FullBanner, error) {
	ctx, span := otel.Tracer("").Start(ctx, "BannersUC.GetBanner")
	defer span.End()

	if !getBannerParams.UseLastVersion {
		fullBanner, err := b.bannersRedisRepo.GetBannerRedis(ctx, getBannerParams.ToGetBannerRedis())
		if err != nil && !errors.Is(err, fiber.ErrNotFound) {
			return nil, err
		}
		if fullBanner != nil {
			if fullBanner.IsActive == false && getBannerParams.AuthToken == constant.UserToken {
				return nil, fiber.NewError(fiber.ErrNotFound.Code, fmt.Sprintf("BannersUC.GetBanner.NotAdminRedis; err = %s", err.Error()))
			}
			return fullBanner, nil
		}
	}

	fullBanner := &models.FullBanner{}
	err := b.trManager.Do(ctx, func(ctx context.Context) error {
		banner, err := b.bannersPGRepo.GetBannerPostgres(ctx, &banners_repository.GetPostgresBanner{
			FeatureId: getBannerParams.FeatureId,
			TagId:     getBannerParams.TagId,
		})
		if err != nil {
			return err
		}
		possibleTagIds, err := b.bannersPGRepo.GetPossibleTagIds(ctx, banner.BannerId)
		if err != nil {
			return err
		}
		fullBanner = banner.ToFullBanner(possibleTagIds)
		return nil
	})
	if err != nil {
		return nil, err
	}

	if !getBannerParams.UseLastVersion {
		if err = b.bannersRedisRepo.PutBannerRedis(ctx, ToPutRedisBanner(fullBanner)); err != nil {
			return nil, err
		}
	}

	if fullBanner.IsActive == false && getBannerParams.AuthToken == constant.UserToken {
		return nil, fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("BannersUC.GetBanner.NotAdminPostgres; err = %s", err.Error()))
	}
	return fullBanner, nil
}

func (b *BannersUC) GetManyBanner(ctx context.Context, getManyBannerParams *GetManyBanner) (*[]models.FullBanner, error) {
	ctx, span := otel.Tracer("").Start(ctx, "BannersUC.GetManyBanner")
	defer span.End()

	manyBannerInfo := &[]models.FullBanner{}
	err := b.trManager.Do(ctx, func(ctx context.Context) error {
		//TODO объединить запросы getManyBanner + getTagIds
		manyBanner, err := b.bannersPGRepo.GetManyBannerPostgres(ctx, getManyBannerParams.ToGetManyPostgresBanner())
		if err != nil {
			return err
		}

		if *manyBanner == nil {
			return nil
		}

		var bannerIds []models.BannerId
		for _, banner := range *manyBanner {
			bannerIds = append(bannerIds, banner.BannerId)
		}

		manyBannerInfo, err = b.bannersPGRepo.GetManyPossibleTagIds(ctx, bannerIds, manyBanner)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return manyBannerInfo, nil
}

func (b *BannersUC) AddBanner(ctx context.Context, addBannerParams *AddBanner) (models.BannerId, error) {
	ctx, span := otel.Tracer("").Start(ctx, "BannersUC.GetManyBanner")
	defer span.End()

	var bannerId models.BannerId
	err := b.trManager.Do(ctx, func(ctx context.Context) error {
		existBanners, err := b.bannersPGRepo.CheckExist(ctx, &banners_repository.CheckExistBanner{TagId: addBannerParams.TagIds, FeatureId: addBannerParams.FeatureId})
		if err != nil {
			return err
		}

		if len(*existBanners) != 0 {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("BannersUC.GetBanner.NotAdminRedis; err = these banners already exists %v", *existBanners))
		}

		insertParams, err := b.bannersPGRepo.AddBannerPostgres(ctx, addBannerParams.ToAddBannerPostgres())
		if err != nil {
			return err
		}

		err = b.bannersPGRepo.AddTags(ctx, &banners_repository.AddTagsPostgres{TagId: addBannerParams.TagIds, BannerId: insertParams.BannerId})
		if err != nil {
			return err
		}

		bannerId = insertParams.BannerId
		return nil
	})
	if err != nil {
		return -1, err
	}

	return bannerId, nil
}
