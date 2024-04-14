package banners_usecase

import (
	"avito/assignment/config"
	"avito/assignment/internal/banners/banners_repository"
	"avito/assignment/internal/models"
	"avito/assignment/pkg/constant"
	"avito/assignment/pkg/errlst"
	"avito/assignment/pkg/traces"
	"avito/assignment/pkg/utilities"
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

// GetBanner (берем случай с use_last_version = false, так как он сложнее)
// 1. Запрашиваем редис отдать запись, если удается - то сразу ее возвращаем
// 2. Если нам не удалось ее получить, то берем ее из постгреса и кладем в редис
func (b *BannersUC) GetBanner(ctx context.Context, getBannerParams *GetBanner) (*models.FullBanner, error) {
	ctx, span := otel.Tracer("").Start(ctx, "BannersUC.GetBanner")
	defer span.End()

	if !getBannerParams.UseLastVersion {
		fullBanner, err := b.bannersRedisRepo.GetBannerRedis(ctx, getBannerParams.FeatureId, getBannerParams.TagId)
		if err != nil && !errors.Is(err, fiber.ErrNotFound) {
			return nil, err
		}
		if fullBanner != nil {
			if fullBanner.IsActive == false && getBannerParams.AuthToken == constant.UserToken {
				return nil, traces.SpanSetErrWrap(span, errlst.HttpErrNotFound, nil, "BannersUC.GetBanner.NotAdmin")
			}
			return fullBanner, nil
		}
	}

	fullBanner := &models.FullBanner{}
	err := b.trManager.Do(ctx, func(ctx context.Context) error {
		banner, err := b.bannersPGRepo.GetBanner(ctx, getBannerParams.FeatureId, getBannerParams.TagId)
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
		return nil, traces.SpanSetErrWrap(span, errlst.HttpErrNotFound, nil, "BannersUC.GetBanner.NotAdmin")
	}
	return fullBanner, nil
}

// GetManyBanner
// 1. Просим из постгреса все записи из бд с баннерами соответствующие требованиям
// 2. Добавляем к каждой записи тэг айдишники из бд с тэгами
func (b *BannersUC) GetManyBanner(ctx context.Context, getManyBannerParams *GetManyBanner) (*[]models.FullBanner, error) {
	ctx, span := otel.Tracer("").Start(ctx, "BannersUC.GetManyBanner")
	defer span.End()

	manyBannerInfo := &[]models.FullBanner{}
	err := b.trManager.Do(ctx, func(ctx context.Context) error {
		manyBanner, err := b.bannersPGRepo.GetManyBanner(ctx, getManyBannerParams.ToGetManyPostgresBanner())
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

// AddBanner
// 1. Проверяем существует ли уже запись в бд с соответствующими фич тэг айдишниками
// 2. Добавляем запись в бд с баннерами
// 3. Добавляем записи в бд с тэгами
func (b *BannersUC) AddBanner(ctx context.Context, addBannerParams *AddBanner) (models.BannerId, error) {
	ctx, span := otel.Tracer("").Start(ctx, "BannersUC.GetManyBanner")
	defer span.End()

	var bannerId models.BannerId
	err := b.trManager.Do(ctx, func(ctx context.Context) error {
		existBanners, err := b.bannersPGRepo.CheckExist(ctx, addBannerParams.TagIds, addBannerParams.FeatureId)
		if err != nil {
			return err
		}

		if len(*existBanners) != 0 {
			return traces.SpanSetErrWrap(span, errlst.HttpErrInvalidRequest,
				errors.New(fmt.Sprintf("these banners already exists %v", *existBanners)), "BannersUC.AddBanner.AlreadyExists")
		}

		insertParams, err := b.bannersPGRepo.AddBanner(ctx, addBannerParams.ToAddBannerPostgres())
		if err != nil {
			return err
		}

		err = b.bannersPGRepo.AddTags(ctx, addBannerParams.TagIds, insertParams.BannerId)
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

// PatchBanner Нечитабельный мусор, в readme добавлю че произошло
// 1. Проверяем существует ли баннер, который надо обновить
// 2. Проверяем способны ли мы добавить в бд запись
// 3. Проверяем обновляем ли мы хоть что то в существующей записи
// 4. Обновляю баннер
// 5. Удаляю + добавляю тэги, чтобы они соответствовали запросу
// 6. Добавляю версию в бд с версиями
// 7. Удаляю пятую версию баннера в бд (костыльно, но лучше не придумать наверное)
func (b *BannersUC) PatchBanner(ctx context.Context, patchBannerParams *PatchBanner) error {
	ctx, span := otel.Tracer("").Start(ctx, "BannersUC.PatchBanner")
	defer span.End()

	err := b.trManager.Do(ctx, func(ctx context.Context) error {
		prevBanner, err := b.bannersPGRepo.GetBannerById(ctx, patchBannerParams.BannerId)
		if err != nil {
			return err
		}
		if prevBanner.BannerId == 0 {
			return traces.SpanSetErrWrap(span, errlst.HttpErrNotFound,
				errors.New(fmt.Sprintf("impossible to update, banner with id %d doesnt exist", patchBannerParams.BannerId)), "BannersUC.PatchBanner.ErrNotFound")
		}

		var addTags []models.TagId
		var existBanners *[]banners_repository.ExistBanner
		if patchBannerParams.TagIds != nil {
			addTags = utilities.FindUniqueElements(*patchBannerParams.TagIds, prevBanner.TagIds)
		}
		if patchBannerParams.FeatureId != nil && *patchBannerParams.FeatureId != prevBanner.FeatureId {
			if patchBannerParams.TagIds != nil {
				existBanners, err = b.bannersPGRepo.CheckExist(ctx, *patchBannerParams.TagIds, *patchBannerParams.FeatureId)
				if err != nil {
					return err
				}
			} else {
				existBanners, err = b.bannersPGRepo.CheckExist(ctx, prevBanner.TagIds, *patchBannerParams.FeatureId)
				if err != nil {
					return err
				}
			}
		} else {
			if patchBannerParams.TagIds != nil {
				existBanners, err = b.bannersPGRepo.CheckExist(ctx, addTags, prevBanner.FeatureId)
				if err != nil {
					return err
				}
			} else {
				existBanners = &[]banners_repository.ExistBanner{}
			}
		}

		if len(*existBanners) != 0 {
			return traces.SpanSetErrWrap(span, errlst.HttpErrInvalidRequest,
				errors.New(fmt.Sprintf("these banners already exists {tag_id feature_id} %v", *existBanners)), "BannersUC.PatchBanner.AlreadyExists")
		}

		if patchBannerParams.Check(prevBanner) {
			return traces.SpanSetErrWrap(span, errlst.HttpErrInvalidRequest,
				errors.New(fmt.Sprintf("nothing to update")), "BannersUC.PatchBanner.NothingToUpdate")
		}

		err = b.bannersPGRepo.UpdateBannerById(ctx, patchBannerParams.ToPatchBanner(prevBanner.Version))
		if err != nil {
			return err
		}
		if patchBannerParams.TagIds != nil {
			if len(addTags) != 0 {
				err = b.bannersPGRepo.AddTags(ctx, addTags, prevBanner.BannerId)
				if err != nil {
					return err
				}
			}
			deleteTags := utilities.FindUniqueElements(prevBanner.TagIds, *patchBannerParams.TagIds)
			if len(deleteTags) != 0 {
				err = b.bannersPGRepo.DeleteTags(ctx, deleteTags, prevBanner.BannerId)
				if err != nil {
					return err
				}
			}
		}

		err = b.bannersPGRepo.AddVersion(ctx, prevBanner)
		if err != nil {
			return err
		}

		err = b.bannersPGRepo.DeleteVersion(ctx, []int64{prevBanner.Version - 3}, patchBannerParams.BannerId)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// DeleteBanner
// 1. Проверяю существует ли запись, которую я хочу удалить (в readme добавлю кое че по этому поводу)
// 2. Удаляю версии, тэги и сам баннер
func (b *BannersUC) DeleteBanner(ctx context.Context, bannerId models.BannerId) error {
	ctx, span := otel.Tracer("").Start(ctx, "BannersUC.DeleteBanner")
	defer span.End()

	err := b.trManager.Do(ctx, func(ctx context.Context) error {
		prevBanner, err := b.bannersPGRepo.GetBannerById(ctx, bannerId)
		if err != nil {
			return err
		}
		if prevBanner.BannerId == 0 {
			return traces.SpanSetErrWrap(span, errlst.HttpErrNotFound,
				errors.New(fmt.Sprintf("impossible to delete, banner with id %d doesnt exist", bannerId)), "BannersUC.DeleteBanner.DoNotExist")
		}

		err = b.bannersPGRepo.DeleteVersion(ctx, []int64{}, bannerId)
		if err != nil {
			return err
		}
		err = b.bannersPGRepo.DeleteTags(ctx, []models.TagId{}, bannerId)
		if err != nil {
			return err
		}
		err = b.bannersPGRepo.DeleteBannerById(ctx, bannerId)
		if err != nil {
			return err
		}
		for _, tagId := range prevBanner.TagIds {
			err = b.bannersRedisRepo.DelBannerRedis(ctx, prevBanner.FeatureId, tagId)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// ViewVersions
// 1. Проверяем существует ли баннер, чьи версии мы хотим просмотреть
// 2. Добавляем в слайс версий текущую версию и остальные
func (b *BannersUC) ViewVersions(ctx context.Context, bannerId models.BannerId) (*[]models.FullBanner, error) {
	ctx, span := otel.Tracer("").Start(ctx, "BannersUC.ViewVersions")
	defer span.End()

	var fullBanners []models.FullBanner
	err := b.trManager.Do(ctx, func(ctx context.Context) error {
		banner, err := b.bannersPGRepo.GetBannerById(ctx, bannerId)
		if err != nil {
			return err
		}

		if banner.BannerId == 0 {
			return traces.SpanSetErrWrap(span, errlst.HttpErrNotFound,
				errors.New(fmt.Sprintf("impossible to view versions, banner with id %d doesnt exist", bannerId)), "BannersUC.ViewVersions.DoNotExist")
		}
		fullBanners = append(fullBanners, *banner)

		banners, err := b.bannersPGRepo.GetBannerVersions(ctx, bannerId, []int64{})
		if err != nil {
			return err
		}
		fullBanners = append(fullBanners, *banners...)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &fullBanners, nil
}

// BannerRollback
// 1. Проверяем существут ли баннер под таким айди и с такой версией
// 2. Закидываем баннер в ручку patch
func (b *BannersUC) BannerRollback(ctx context.Context, bannerId models.BannerId, version int64) error {
	ctx, span := otel.Tracer("").Start(ctx, "BannersUC.BannerRollback")
	defer span.End()

	err := b.trManager.Do(ctx, func(ctx context.Context) error {
		banner, err := b.bannersPGRepo.GetBannerVersions(ctx, bannerId, []int64{version})
		if err != nil {
			return err
		}

		if len(*banner) == 0 {
			return traces.SpanSetErrWrap(span, errlst.HttpErrNotFound,
				errors.New(fmt.Sprintf("impossible to rollback, banner with id %d and version %d doesnt exist", bannerId, version)), "BannersUC.BannerRollback.DoNotExist")
		}

		err = b.PatchBanner(ctx, ToPatchBanner((*banner)[0]))
		if err != nil {
			return err
		}

		for _, tagId := range (*banner)[0].TagIds {
			err = b.bannersRedisRepo.DelBannerRedis(ctx, (*banner)[0].FeatureId, tagId)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
