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
		var possibleBannerIds []models.BannerId
		if getManyBannerParams.TagId != nil {
			bannerIds, err := b.bannersPGRepo.GetPossibleBannerIds(ctx, *getManyBannerParams.TagId)
			if err != nil {
				return err
			}
			possibleBannerIds = bannerIds
		}

		manyBanner, err := b.bannersPGRepo.GetManyBannerPostgres(ctx, getManyBannerParams.ToGetManyPostgresBanner(possibleBannerIds))
		if err != nil {
			return err
		}

		for _, banner := range *manyBanner {
			//TODO Оптимизировать добавление тэгайдишников
			tagIds, err := b.bannersPGRepo.GetPossibleTagIds(ctx, banner.BannerId)
			if err != nil {
				return err
			}
			*manyBannerInfo = append(*manyBannerInfo, *banner.ToFullBanner(tagIds))
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
		//TODO обдумать более быструю логику добавления баннера = На текущий момент мы
		// 1) Проверяем можем ли мы добавить баннер (постоянный innerjoin и цикл)
		// 2) Добавляем баннер
		// 3) Добавляем тэгайдишники в баннер (по циклу)
		for _, tagId := range addBannerParams.TagIds {
			err := b.bannersPGRepo.CheckExist(ctx, &banners_repository.CheckExistBanner{TagId: tagId, FeatureId: addBannerParams.FeatureId})
			if err != nil {
				return err
			}
		}

		insertParams, err := b.bannersPGRepo.AddBannerPostgres(ctx, addBannerParams.ToAddBannerPostgres())
		if err != nil {
			return err
		}

		for _, tagId := range addBannerParams.TagIds {
			err = b.bannersPGRepo.AddTags(ctx, &banners_repository.AddTagsPostgres{TagId: tagId, BannerId: insertParams.BannerId})
			if err != nil {
				return err
			}
		}

		if err = b.bannersRedisRepo.PutBannerRedis(ctx, &banners_repository.PutRedisBanner{}); err != nil {
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
