package banners_http

import (
	"avito/assignment/internal/banners/banners_usecase"
	"avito/assignment/internal/models"
	"avito/assignment/pkg/utilities"
	"time"
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

type PatchBannerRequest struct {
	TagIds    *[]models.TagId   `json:"tag_ids"`
	FeatureId *models.FeatureId `json:"feature_id"`
	Content   *struct {
		Title *string `json:"title"`
		Text  *string `json:"text"`
		Url   *string `json:"url"`
	} `json:"content"`
	IsActive *bool `json:"is_active"`
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
		TagIds:    utilities.RemoveDuplicates[models.TagId](b.TagIds),
		FeatureId: b.FeatureId,
		Content: struct {
			Title string
			Text  string
			Url   string
		}{Title: b.Content.Title, Text: b.Content.Text, Url: b.Content.Url},
		IsActive: b.IsActive,
	}
}

func (b *PatchBannerRequest) ToPatchBanner(bannerId models.BannerId) *banners_usecase.PatchBanner {
	if b.Content == nil {
		return &banners_usecase.PatchBanner{
			TagIds: func(*[]models.TagId) *[]models.TagId {
				if b.TagIds != nil {
					ids := utilities.RemoveDuplicates[models.TagId](*b.TagIds)
					return &ids
				}
				return nil
			}(b.TagIds),
			FeatureId: b.FeatureId,
			IsActive:  b.IsActive,
			BannerId:  bannerId,
		}
	}
	return &banners_usecase.PatchBanner{
		TagIds: func(*[]models.TagId) *[]models.TagId {
			if b.TagIds != nil {
				ids := utilities.RemoveDuplicates[models.TagId](*b.TagIds)
				return &ids
			}
			return nil
		}(b.TagIds),
		FeatureId: b.FeatureId,
		Title:     b.Content.Title,
		Text:      b.Content.Text,
		Url:       b.Content.Url,
		IsActive:  b.IsActive,
		BannerId:  bannerId,
	}
}

type GetBannerResponse struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	Url   string `json:"url"`
}

func ToGetBannerResponse(b *models.FullBanner) *GetBannerResponse {
	return &GetBannerResponse{
		Title: b.Content.Title,
		Text:  b.Content.Text,
		Url:   b.Content.Url,
	}
}

type GetManyBannerResponse struct {
	BannerId  models.BannerId  `json:"banner_id"`
	TagIds    []models.TagId   `json:"tag_ids"`
	FeatureId models.FeatureId `json:"feature_id"`
	Content   struct {
		Title string `json:"title"`
		Text  string `json:"text"`
		Url   string `json:"url"`
	} `json:"content"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Version   int64     `json:"version"`
}

func ToGetManyBannerResponse(b *[]models.FullBanner) *[]GetManyBannerResponse {
	getManyBannerResponse := make([]GetManyBannerResponse, len(*b))

	for i, value := range *b {
		getManyBannerResponse[i] = GetManyBannerResponse(value)
	}

	return &getManyBannerResponse
}
