package banners_usecase

import (
	"avito/assignment/internal/banners/banners_repository"
	"avito/assignment/internal/models"
)

type GetBanner struct {
	TagId          models.TagId
	FeatureId      models.FeatureId
	UseLastVersion bool
	AuthToken      string
}

func (b *GetBanner) ToGetBannerRedis() *banners_repository.GetRedisBanner {
	return &banners_repository.GetRedisBanner{
		TagId:     b.TagId,
		FeatureId: b.FeatureId,
	}
}

type GetManyBanner struct {
	FeatureId *models.FeatureId
	TagId     *models.TagId
	Limit     *int
	Offset    *int
}

func (b *GetManyBanner) ToGetManyPostgresBanner() *banners_repository.GetManyPostgresBanner {
	return &banners_repository.GetManyPostgresBanner{
		TagId:     b.TagId,
		FeatureId: b.FeatureId,
		Limit:     b.Limit,
		Offset:    b.Offset,
	}
}

type AddBanner struct {
	TagIds    []models.TagId
	FeatureId models.FeatureId
	Content   struct {
		Title string
		Text  string
		Url   string
	}
	IsActive bool
}

func (b *AddBanner) ToAddBannerPostgres() *banners_repository.AddPostgresBanner {
	return &banners_repository.AddPostgresBanner{
		FeatureId: b.FeatureId,
		Title:     b.Content.Title,
		Text:      b.Content.Text,
		Url:       b.Content.Url,
		IsActive:  b.IsActive,
	}
}

func (b *AddBanner) ToPutRedisBanner(getInsertParams *banners_repository.GetInsertParams) *banners_repository.PutRedisBanner {
	return &banners_repository.PutRedisBanner{
		BannerId:  getInsertParams.BannerId,
		TagIds:    b.TagIds,
		FeatureId: b.FeatureId,
		Content: struct {
			Title string
			Text  string
			Url   string
		}{Title: b.Content.Title, Text: b.Content.Text, Url: b.Content.Url},
		IsActive:  b.IsActive,
		CreatedAt: getInsertParams.CreatedAt,
		UpdatedAt: getInsertParams.UpdatedAt,
	}
}

func ToPutRedisBanner(fullBanner *models.FullBanner) *banners_repository.PutRedisBanner {
	return &banners_repository.PutRedisBanner{
		BannerId:  fullBanner.BannerId,
		TagIds:    fullBanner.TagIds,
		FeatureId: fullBanner.FeatureId,
		Content: struct {
			Title string
			Text  string
			Url   string
		}{Title: fullBanner.Content.Title, Text: fullBanner.Content.Text, Url: fullBanner.Content.Url},
		IsActive:  fullBanner.IsActive,
		CreatedAt: fullBanner.CreatedAt,
		UpdatedAt: fullBanner.UpdatedAt,
	}
}
