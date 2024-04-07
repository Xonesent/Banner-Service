package banners_http

import (
	"avito/assignment/config"
	"avito/assignment/internal/banners"
	"avito/assignment/internal/models/banner_models"
	"avito/assignment/pkg/traces"
	reqvalidator "avito/assignment/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

type BannersHandlers struct {
	bannersUC banners.Usecase
	cfg       *config.Config
}

func NewUserHandler(bannersUC banners.Usecase, cfg *config.Config) *BannersHandlers {
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

		params := banner_models.GetBanner{}
		if err := reqvalidator.ReadRequest(c, &params); err != nil {
			return traces.SpanSetErrWrapf(
				span,
				fiber.ErrBadRequest,
				"BannersHandlers.GetBanner.ReadRequest(args:%v)",
				params,
			)
		}
		params.AuthToken = token

		bannerInfo, err := b.bannersUC.GetBanner(ctx, params)
		if err != nil {
			return err
		}

		return c.JSON(fiber.Map{
			"title": bannerInfo.Title,
			"text":  bannerInfo.Text,
			"url":   bannerInfo.Url,
		})
	}
}

func (b *BannersHandlers) GetManyBanner() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, span := traces.StartFiberTrace(c, "BannersHandlers.GetBanner")
		defer span.End()

		params := banner_models.GetManyBanner{}
		if err := reqvalidator.ReadRequest(c, &params); err != nil {
			return traces.SpanSetErrWrapf(
				span,
				fiber.ErrBadRequest,
				"BannersHandlers.GetBanner.ReadRequest(args:%v)",
				params,
			)
		}

		manyBannerInfo, err := b.bannersUC.GetManyBanner(ctx, params)
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

		params := banner_models.AddBanner{}
		if err := reqvalidator.ReadRequest(c, &params); err != nil {
			return traces.SpanSetErrWrapf(
				span,
				fiber.ErrBadRequest,
				"BannersHandlers.GetBanner.ReadRequest(args:%v)",
				params,
			)
		}

		bannerId, err := b.bannersUC.AddBanner(ctx, params)
		if err != nil {
			return err
		}

		return c.JSON(fiber.Map{
			"banner_id": bannerId,
		})
	}
}
