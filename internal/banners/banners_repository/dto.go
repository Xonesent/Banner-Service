package banners_repository

import (
	"avito/assignment/internal/models"
	"strconv"
	"strings"
	"time"
)

type GetManyPostgresBanner struct {
	TagId     *models.TagId
	FeatureId *models.FeatureId
	Limit     *int
	Offset    *int
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

type ExistBanner struct {
	TagId     models.TagId     `db:"tag_id"`
	FeatureId models.FeatureId `db:"feature_id"`
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

type UpdateBannerById struct {
	FeatureId *models.FeatureId
	Title     *string
	Text      *string
	Url       *string
	IsActive  *bool
	BannerId  models.BannerId
	Version   int64
}

type FullBanner struct {
	BannerId  models.BannerId  `db:"banner_id"`
	TagIds    []uint8          `db:"tag_ids"`
	FeatureId models.FeatureId `db:"feature_id"`
	Title     string           `db:"title"`
	Text      string           `db:"text"`
	Url       string           `db:"url"`
	IsActive  bool             `db:"is_active"`
	CreatedAt time.Time        `db:"created_at"`
	UpdatedAt time.Time        `db:"updated_at"`
	Version   int64            `db:"version"`
}

func (b *FullBanner) ToFullBanners() models.FullBanner {
	values := strings.Split(strings.Trim(string(b.TagIds), "{}"), ",")
	var tagIds []models.TagId
	for _, value := range values {
		num, _ := strconv.ParseInt(value, 10, 64)
		tagIds = append(tagIds, models.TagId(num))
	}

	return models.FullBanner{
		BannerId:  b.BannerId,
		TagIds:    tagIds,
		FeatureId: b.FeatureId,
		Content: struct {
			Title string
			Text  string
			Url   string
		}{Title: b.Title, Text: b.Text, Url: b.Url},
		IsActive:  b.IsActive,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
		Version:   b.Version,
	}
}
