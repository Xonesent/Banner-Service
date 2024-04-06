package banners_http

import (
	"avito/assignment/config"
	"avito/assignment/internal/banners"
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
