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

func (b *BannersUC) GetManyBanner(ctx context.Context, params banner_models.GetManyBanner) (*[]banner_models.EditedFullBannerContent, error) {
	ctx, span := otel.Tracer("").Start(ctx, "BannersUC.GetManyBanner")
	defer span.End()

	manyBannerInfo := &[]banner_models.EditedFullBannerContent{}
	err := b.trManager.Do(ctx, func(ctx context.Context) error {
		var possibleBannerIds []int
		if params.TagId != 0 {
			bannerIds, err := b.bannersRepo.GetPossibleBannerIds(ctx, params.TagId)
			if err != nil {
				return err
			}
			possibleBannerIds = bannerIds
		}

		manyBanner, err := b.bannersRepo.SelectBanner(ctx, banner_models.SelectPostgresBanner{
			TagId:             params.TagId,
			FeatureId:         params.FeatureId,
			PossibleBannerIds: possibleBannerIds,
			Offset:            params.Offset,
			Limit:             params.Limit,
		})
		if err != nil {
			return err
		}

		for _, banner := range *manyBanner {
			//TODO Оптимизировать добавление тэгайдишников
			tagIds, err := b.bannersRepo.GetPossibleTagIds(ctx, banner.BannerId)
			if err != nil {
				return err
			}
			*manyBannerInfo = append(*manyBannerInfo, banner_models.EditBannerContent(banner, tagIds))
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return manyBannerInfo, nil
}

func (b *BannersUC) AddBanner(ctx context.Context, params banner_models.AddBanner) (*int, error) {
	ctx, span := otel.Tracer("").Start(ctx, "BannersUC.GetManyBanner")
	defer span.End()

	var bannerId *int
	err := b.trManager.Do(ctx, func(ctx context.Context) error {
		//TODO обдумать более быструю логику добавления баннера = На текущий момент мы
		// 1) Проверяем можем ли мы добавить баннер (постоянный innerjoin и цикл)
		// 2) Добавляем баннер
		// 3) Добавляем тэгайдишники в баннер (по циклу)
		for _, tagId := range params.TagIds {
			err := b.bannersRepo.CheckExist(ctx, tagId, params.FeatureId)
			if err != nil {
				return err
			}
		}

		id, err := b.bannersRepo.AddBanner(ctx, params)
		if err != nil {
			return err
		}

		for _, tagId := range params.TagIds {
			err = b.bannersRepo.AddTags(ctx, *id, tagId)
			if err != nil {
				return err
			}
		}

		if err := b.redisClient.PutBanner(ctx, banner_models.PutRedisBanner{
			TagIds:    params.TagIds,
			FeatureId: params.FeatureId,
			Content: banner_models.BannerContent{
				Title:    params.Content.Title,
				Text:     params.Content.Text,
				Url:      params.Content.Url,
				IsActive: params.IsActive,
			},
		}); err != nil {
			return err
		}

		bannerId = id
		return nil
	})
	if err != nil {
		return nil, err
	}

	return bannerId, nil
}
