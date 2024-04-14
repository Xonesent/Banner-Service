package banners_http

import (
	"avito/assignment/config"
	"avito/assignment/internal/models"
	"avito/assignment/pkg/errlst"
	"avito/assignment/pkg/traces"
	reqvalidator "avito/assignment/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type BannersHandlers struct {
	bannersUC BannersUseCase
	cfg       *config.Config
}

func NewUserHandler(bannersUC BannersUseCase, cfg *config.Config) *BannersHandlers {
	return &BannersHandlers{
		bannersUC: bannersUC,
		cfg:       cfg,
	}
}

func (b *BannersHandlers) GetBanner() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, span := traces.StartFiberTrace(c, "BannersHandlers.GetBanner")
		defer span.End()

		token := c.Locals("token").(string)

		getBanner := GetBannerRequest{}
		if err := reqvalidator.ReadRequest(c, &getBanner); err != nil {
			return traces.SpanSetErrWrap(span, errlst.HttpErrInvalidRequest, err, "BannersHandlers.GetBanner.ReadRequest")
		}

		getBannerDTO := getBanner.ToGetBanner()
		getBannerDTO.AuthToken = token

		bannerInfo, err := b.bannersUC.GetBanner(ctx, getBannerDTO)
		if err != nil {
			return err
		}

		return c.JSON(ToGetBannerResponse(bannerInfo))
	}
}

func (b *BannersHandlers) GetManyBanner() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, span := traces.StartFiberTrace(c, "BannersHandlers.GetManyBanner")
		defer span.End()

		getManyBanner := GetManyBannerRequest{}
		if err := reqvalidator.ReadRequest(c, &getManyBanner); err != nil {
			return traces.SpanSetErrWrap(span, errlst.HttpErrInvalidRequest, err, "BannersHandlers.GetManyBanner.ReadRequest")
		}

		getManyBannerDTO := getManyBanner.ToGetManyBanner()

		manyBannerInfo, err := b.bannersUC.GetManyBanner(ctx, getManyBannerDTO)
		if err != nil {
			return err
		}

		return c.JSON(ToGetManyBannerResponse(manyBannerInfo))
	}
}

func (b *BannersHandlers) AddBanner() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, span := traces.StartFiberTrace(c, "BannersHandlers.AddBanner")
		defer span.End()

		addBanner := AddBannerRequest{}
		if err := reqvalidator.ReadRequest(c, &addBanner); err != nil {
			return traces.SpanSetErrWrap(span, errlst.HttpErrInvalidRequest, err, "BannersHandlers.AddBanner.ReadRequest")
		}
		if len(addBanner.TagIds) == 0 {
			return traces.SpanSetErrWrap(span, errlst.HttpErrInvalidRequest, nil, "BannersHandlers.AddBanner.NilTagIds")
		}

		addBannerDTO := addBanner.ToAddBanner()

		bannerId, err := b.bannersUC.AddBanner(ctx, addBannerDTO)
		if err != nil {
			return err
		}

		c.Status(fiber.StatusCreated)
		return c.JSON(fiber.Map{
			"banner_id": bannerId,
		})
	}
}

func (b *BannersHandlers) PatchBanner() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, span := traces.StartFiberTrace(c, "BannersHandlers.PatchBanner")
		defer span.End()

		patchBanner := PatchBannerRequest{}
		if err := reqvalidator.ReadRequest(c, &patchBanner); err != nil {
			return traces.SpanSetErrWrap(span, errlst.HttpErrInvalidRequest, nil, "BannersHandlers.PatchBanner.ReadRequest")
		}
		bannerId, err := strconv.Atoi(c.Params("banner_id"))
		if err != nil {
			return traces.SpanSetErrWrap(span, errlst.HttpErrInvalidRequest, nil, "BannersHandlers.PatchBanner.WrongBannerParams")
		}

		patchBannerDTO := patchBanner.ToPatchBanner(models.BannerId(bannerId))

		err = b.bannersUC.PatchBanner(ctx, patchBannerDTO)
		if err != nil {
			return err
		}

		return c.JSON(fiber.Map{
			"message": "Success",
		})
	}
}

func (b *BannersHandlers) DeleteBanner() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, span := traces.StartFiberTrace(c, "BannersHandlers.DeleteBanner")
		defer span.End()

		bannerId, err := strconv.Atoi(c.Params("banner_id"))
		if err != nil {
			return traces.SpanSetErrWrap(span, errlst.HttpErrInvalidRequest, nil, "BannersHandlers.DeleteBanner.WrongBannerParams")
		}

		err = b.bannersUC.DeleteBanner(ctx, models.BannerId(bannerId))
		if err != nil {
			return err
		}

		c.Status(fiber.StatusNoContent)
		return c.JSON(fiber.Map{
			"message": "Success",
		})
	}
}

func (b *BannersHandlers) ViewVersions() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, span := traces.StartFiberTrace(c, "BannersHandlers.ViewVersions")
		defer span.End()

		bannerId, err := strconv.Atoi(c.Params("banner_id"))
		if err != nil {
			return traces.SpanSetErrWrap(span, errlst.HttpErrInvalidRequest, nil, "BannersHandlers.ViewVersions.WrongBannerParams")
		}

		banners, err := b.bannersUC.ViewVersions(ctx, models.BannerId(bannerId))
		if err != nil {
			return err
		}

		return c.JSON(banners)
	}
}

func (b *BannersHandlers) BannerRollback() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, span := traces.StartFiberTrace(c, "BannersHandlers.PatchBanner")
		defer span.End()

		bannerId, err := strconv.Atoi(c.Params("banner_id"))
		if err != nil {
			return traces.SpanSetErrWrap(span, errlst.HttpErrInvalidRequest, nil, "BannersHandlers.ViewVersions.WrongBannerParams")
		}
		version, err := strconv.Atoi(c.Params("version"))
		if err != nil {
			return traces.SpanSetErrWrap(span, errlst.HttpErrInvalidRequest, nil, "BannersHandlers.ViewVersions.WrongVersionParams")
		}

		err = b.bannersUC.BannerRollback(ctx, models.BannerId(bannerId), int64(version))
		if err != nil {
			return err
		}

		return c.JSON(fiber.Map{
			"message": "Success",
		})
	}
}
