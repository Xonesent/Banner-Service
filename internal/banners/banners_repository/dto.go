package banners_repository

import (
	"avito/assignment/internal/models"
	"time"
)

type GetPostgresBanner struct {
	FeatureId models.FeatureId
	TagId     models.TagId
}

type GetManyPostgresBanner struct {
	TagId             *models.TagId
	FeatureId         *models.FeatureId
	PossibleBannerIds []models.BannerId
	Limit             *int
	Offset            *int
}

type AddPostgresBanner struct {
	FeatureId models.FeatureId
	Title     string
	Text      string
	Url       string
	IsActive  bool
}

type GetInsertParams struct {
	BannerId  models.BannerId
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CheckExistBanner struct {
	TagId     models.TagId
	FeatureId models.FeatureId
}

type AddTagsPostgres struct {
	TagId    models.TagId
	BannerId models.BannerId
}

type PutRedisBanner struct {
	BannerId  models.BannerId
	TagIds    []models.TagId
	FeatureId models.FeatureId
	Content   struct {
		Title string
		Text  string
		Url   string
	}
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type GetRedisBanner struct {
	TagId     models.TagId
	FeatureId models.FeatureId
}
