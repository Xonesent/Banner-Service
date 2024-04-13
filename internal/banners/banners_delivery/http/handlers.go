package banners_http

import (
	"avito/assignment/config"
	"avito/assignment/internal/models"
	"avito/assignment/pkg/traces"
	reqvalidator "avito/assignment/pkg/validator"
	"fmt"
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
			return fiber.NewError(fiber.ErrBadRequest.Code, fmt.Sprintf("BannersHandlers.GetBanner.ReadRequest; err = %s", err.Error()))
		}

		getBannerDTO := getBanner.ToGetBanner()
		getBannerDTO.AuthToken = token

		bannerInfo, err := b.bannersUC.GetBanner(ctx, getBannerDTO)
		if err != nil {
			return err
		}

		return c.JSON(bannerInfo)
	}
}

func (b *BannersHandlers) GetManyBanner() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, span := traces.StartFiberTrace(c, "BannersHandlers.GetBanner")
		defer span.End()

		getManyBanner := GetManyBannerRequest{}
		if err := reqvalidator.ReadRequest(c, &getManyBanner); err != nil {
			return fiber.NewError(fiber.ErrBadRequest.Code, fmt.Sprintf("BannersHandlers.GetBanner.ReadRequest; err = %s", err.Error()))
		}

		getManyBannerDTO := getManyBanner.ToGetManyBanner()

		manyBannerInfo, err := b.bannersUC.GetManyBanner(ctx, getManyBannerDTO)
		if err != nil {
			return err
		}

		return c.JSON(manyBannerInfo)
	}
}

func (b *BannersHandlers) AddBanner() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, span := traces.StartFiberTrace(c, "BannersHandlers.GetBanner")
		defer span.End()

		addBanner := AddBannerRequest{}
		if err := reqvalidator.ReadRequest(c, &addBanner); err != nil {
			return fiber.NewError(fiber.ErrBadRequest.Code, fmt.Sprintf("BannersHandlers.GetBanner.ReadRequest; err = %s", err.Error()))
		}
		if len(addBanner.TagIds) == 0 {
			return fiber.NewError(fiber.ErrBadRequest.Code, fmt.Sprintf("BannersHandlers.GetBanner.ReadRequest; err = NilTagIds"))
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
			return fiber.NewError(fiber.ErrBadRequest.Code, fmt.Sprintf("BannersHandlers.PatchBanner.ReadRequest; err = %s", err.Error()))
		}
		bannerId, err := strconv.Atoi(c.Params("banner_id"))
		if err != nil {
			return fiber.NewError(fiber.ErrBadRequest.Code, fmt.Sprintf("BannersHandlers.PatchBanner.Params; err = %s", err.Error()))
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
		ctx, span := traces.StartFiberTrace(c, "BannersHandlers.PatchBanner")
		defer span.End()

		bannerId, err := strconv.Atoi(c.Params("banner_id"))
		if err != nil {
			return fiber.NewError(fiber.ErrBadRequest.Code, fmt.Sprintf("BannersHandlers.PatchBanner.Params; err = %s", err.Error()))
		}

		err = b.bannersUC.DeleteBanner(ctx, models.BannerId(bannerId))
		if err != nil {
			return err
		}

		return c.JSON(fiber.Map{
			"banner_id": bannerId,
		})
	}
}
