package banners_postgres

const (
	BannersTableName      = "banner_schema.banners"
	BannersXTagsTableName = "banner_schema.banners_X_tags"
	BannerIdColumnName    = "banner_id"
	TitleColumnName       = "title"
	TextColumnName        = "text"
	UrlColumnName         = "url"
	FeatureIdColumnName   = "feature_id"
	CreatedAtColumnName   = "created_at"
	UpdatedAtColumnName   = "updated_at"
	IsActiveColumnName    = "is_active"
	IdColumnName          = "id"
	TagIdColumnName       = "tag_id"
)

var (
	GetBannerColumns = []string{
		TitleColumnName,
		TextColumnName,
		UrlColumnName,
		IsActiveColumnName,
	}
	SelectBannerColumns = []string{
		BannerIdColumnName,
		FeatureIdColumnName,
		TitleColumnName,
		TextColumnName,
		UrlColumnName,
		IsActiveColumnName,
		CreatedAtColumnName,
		UpdatedAtColumnName,
	}
	InsertBannerColumns = []string{
		TitleColumnName,
		TextColumnName,
		UrlColumnName,
		FeatureIdColumnName,
		IsActiveColumnName,
		CreatedAtColumnName,
		UpdatedAtColumnName,
	}
	InsertTagColumns = []string{
		BannerIdColumnName,
		TagIdColumnName,
	}
)
