package banner_models

import "time"

type GetBanner struct {
	TagId          int    `json:"tag_id" validate:"required"`
	FeatureId      int    `json:"feature_id" validate:"required"`
	UseLastVersion bool   `json:"use_last_version"`
	AuthToken      string `json:"-"`
}

type PutRedisBanner struct {
	TagIds    []int
	FeatureId int
	Content   BannerContent
	CreatedAt time.Time
	UpdatedAt time.Time
}

type GetPostgresBanner struct {
	FeatureId         int
	PossibleBannerIds []int
}

type GetRedisBanner struct {
	TagId     int
	FeatureId int
}

type GetManyBanner struct {
	FeatureId int `json:"feature_id"`
	TagId     int `json:"tag_id"`
	Limit     int `json:"limit"`
	Offset    int `json:"offset"`
}

type BannerContent struct {
	Title    string `db:"title"`
	Text     string `db:"text"`
	Url      string `db:"url"`
	IsActive bool   `db:"is_active"`
}
