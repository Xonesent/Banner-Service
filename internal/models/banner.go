package models

import (
	"time"
)

type Banner struct {
	BannerId  BannerId  `db:"banner_id"`
	TagId     TagId     `db:"tag_id"`
	FeatureId FeatureId `db:"feature_id"`
	Title     string    `db:"title"`
	Text      string    `db:"text"`
	Url       string    `db:"url"`
	IsActive  bool      `db:"is_active"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Version   int64     `db:"version"`
}

type FullBanner struct {
	BannerId  BannerId
	TagIds    []TagId
	FeatureId FeatureId
	Content   struct {
		Title string
		Text  string
		Url   string
	}
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
	Version   int64
}

func (b *Banner) ToFullBanner(tagIds []TagId) *FullBanner {
	return &FullBanner{
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

func (b *Banner) ToFullBannerWithoutTagIds() *FullBanner {
	return &FullBanner{
		BannerId:  b.BannerId,
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
