package banners_usecase

import (
	"avito/assignment/internal/banners/banners_repository"
	"avito/assignment/internal/models"
	"avito/assignment/pkg/utilities"
)

type GetBanner struct {
	TagId          models.TagId
	FeatureId      models.FeatureId
	UseLastVersion bool
	AuthToken      string
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

type PatchBanner struct {
	TagIds    *[]models.TagId
	FeatureId *models.FeatureId
	Title     *string
	Text      *string
	Url       *string
	IsActive  *bool
	BannerId  models.BannerId
}

func (b *PatchBanner) ToPatchBanner(version int64) *banners_repository.UpdateBannerById {
	return &banners_repository.UpdateBannerById{
		FeatureId: b.FeatureId,
		Title:     b.Title,
		Text:      b.Text,
		Url:       b.Url,
		IsActive:  b.IsActive,
		BannerId:  b.BannerId,
		Version:   version,
	}
}

func (b *PatchBanner) Check(banner *models.FullBanner) bool {
	if b.FeatureId == nil && b.IsActive == nil && b.TagIds == nil && b.Url == nil && b.Text == nil && b.Title == nil {
		return true
	}
	var maxCoincidence, currCoincidence int = 6, 0
	if b.FeatureId == nil || b.FeatureId != nil && *b.FeatureId == banner.FeatureId {
		currCoincidence++
	}
	if b.IsActive == nil || b.IsActive != nil && *b.IsActive == banner.IsActive {
		currCoincidence++
	}
	if b.Url == nil || b.Url != nil && *b.Url == banner.Content.Url {
		currCoincidence++
	}
	if b.Text == nil || b.Text != nil && *b.Text == banner.Content.Text {
		currCoincidence++
	}
	if b.Title == nil || b.Title != nil && *b.Title == banner.Content.Title {
		currCoincidence++
	}
	if b.TagIds == nil {
		currCoincidence++
	} else {
		if utilities.AreSlicesEqual(*b.TagIds, banner.TagIds) {
			currCoincidence++
		}
	}
	return currCoincidence == maxCoincidence
}

func ToPatchBanner(banner models.FullBanner) *PatchBanner {
	return &PatchBanner{
		TagIds:    &banner.TagIds,
		FeatureId: &banner.FeatureId,
		Title:     &banner.Content.Title,
		Text:      &banner.Content.Text,
		Url:       &banner.Content.Url,
		IsActive:  &banner.IsActive,
		BannerId:  banner.BannerId,
	}
}
