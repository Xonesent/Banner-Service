package banners_http

import (
	"avito/assignment/internal/banners/banners_usecase"
	"avito/assignment/internal/models"
)

type GetBannerRequest struct {
	TagId          models.TagId     `json:"tag_id" validate:"required"`
	FeatureId      models.FeatureId `json:"feature_id" validate:"required"`
	UseLastVersion bool             `json:"use_last_version"`
}

type GetManyBannerRequest struct {
	FeatureId *models.FeatureId `json:"feature_id"`
	TagId     *models.TagId     `json:"tag_id"`
	Limit     *int              `json:"limit"`
	Offset    *int              `json:"offset"`
}

type AddBannerRequest struct {
	TagIds    []models.TagId   `json:"tag_ids" validate:"required"`
	FeatureId models.FeatureId `json:"feature_id" validate:"required"`
	Content   struct {
		Title string `json:"title" validate:"required"`
		Text  string `json:"text" validate:"required"`
		Url   string `json:"url" validate:"required"`
	} `json:"content"`
	IsActive bool `json:"is_active" validate:"required"`
}

func (b *GetBannerRequest) ToGetBanner() *banners_usecase.GetBanner {
	return &banners_usecase.GetBanner{
		TagId:          b.TagId,
		FeatureId:      b.FeatureId,
		UseLastVersion: b.UseLastVersion,
	}
}

func (b *GetManyBannerRequest) ToGetManyBanner() *banners_usecase.GetManyBanner {
	return &banners_usecase.GetManyBanner{
		FeatureId: b.FeatureId,
		TagId:     b.TagId,
		Limit:     b.Limit,
		Offset:    b.Offset,
	}
}

func (b *AddBannerRequest) ToAddBanner() *banners_usecase.AddBanner {
	return &banners_usecase.AddBanner{
		TagIds:    b.TagIds,
		FeatureId: b.FeatureId,
		Content: struct {
			Title string
			Text  string
			Url   string
		}{Title: b.Content.Title, Text: b.Content.Text, Url: b.Content.Url},
		IsActive: b.IsActive,
	}
}
