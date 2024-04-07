package banner_models

import "time"

type FullBannerContent struct {
	BannerId  int       `db:"banner_id"`
	FeatureId int       `db:"feature_id"`
	Title     string    `db:"title"`
	Text      string    `db:"text"`
	Url       string    `db:"url"`
	IsActive  bool      `db:"is_active"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type EditedFullBannerContent struct {
	BannerId  int   `json:"banner_id"`
	TagIds    []int `json:"tag_ids"`
	FeatureId int   `json:"feature_id"`
	Content   struct {
		Title string `json:"title"`
		Text  string `json:"text"`
		Url   string `json:"url"`
	} `json:"content"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func EditBannerContent(banner FullBannerContent, tagIds []int) EditedFullBannerContent {
	return EditedFullBannerContent{
		BannerId:  banner.BannerId,
		TagIds:    tagIds,
		FeatureId: banner.FeatureId,
		Content: struct {
			Title string `json:"title"`
			Text  string `json:"text"`
			Url   string `json:"url"`
		}(struct {
			Title string
			Text  string
			Url   string
		}{Title: banner.Title, Text: banner.Text, Url: banner.Url}),
		IsActive:  banner.IsActive,
		CreatedAt: banner.CreatedAt,
		UpdatedAt: banner.UpdatedAt,
	}
}

type GetBanner struct {
	TagId          int    `json:"tag_id" validate:"required"`
	FeatureId      int    `json:"feature_id" validate:"required"`
	UseLastVersion bool   `json:"use_last_version"`
	AuthToken      string `json:"-"`
}

type AddBanner struct {
	TagIds    []int   `json:"tag_ids" validate:"required"`
	FeatureId int     `json:"feature_id" validate:"required"`
	Content   Content `json:"content" validate:"required"`
	IsActive  bool    `json:"is_active" validate:"required"`
}

type Content struct {
	Title string `json:"title" validate:"required"`
	Text  string `json:"text" validate:"required"`
	Url   string `json:"url" validate:"required"`
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

type SelectPostgresBanner struct {
	TagId             int
	FeatureId         int
	PossibleBannerIds []int
	Limit             int
	Offset            int
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
