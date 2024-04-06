package banners

import (
	"avito/assignment/internal/models/banner_models"
	"context"
)

type Usecase interface {
	GetBanner(ctx context.Context, params banner_models.GetBanner) (*banner_models.BannerContent, error)
	//GetManyBanner(ctx context.Context, params banner_models.GetManyBanner) (*[]banner_models.BannerContent, error)
}
