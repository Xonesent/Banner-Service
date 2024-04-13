package banners_http

import (
	"avito/assignment/internal/banners/banners_usecase"
	"avito/assignment/internal/models"
	"context"
	"github.com/gofiber/fiber/v2"
)

type Handlers interface {
	GetBanner() fiber.Handler
	GetManyBanner() fiber.Handler
	AddBanner() fiber.Handler
	PatchBanner() fiber.Handler
	DeleteBanner() fiber.Handler
}

type BannersUseCase interface {
	GetBanner(ctx context.Context, getBannerParams *banners_usecase.GetBanner) (*models.FullBanner, error)
	GetManyBanner(ctx context.Context, getManyBannerParams *banners_usecase.GetManyBanner) (*[]models.FullBanner, error)
	AddBanner(ctx context.Context, addBannerParams *banners_usecase.AddBanner) (models.BannerId, error)
	PatchBanner(ctx context.Context, patchBannerParams *banners_usecase.PatchBanner) error
	DeleteBanner(ctx context.Context, bannerId models.BannerId) error
}
