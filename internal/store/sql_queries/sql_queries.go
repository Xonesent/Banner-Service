package sql_queries

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
	GetBannerColumnsWithInnerJoin = []string{
		"b.banner_id",
		TitleColumnName,
		TextColumnName,
		UrlColumnName,
		FeatureIdColumnName,
		CreatedAtColumnName,
		UpdatedAtColumnName,
		IsActiveColumnName,
	}
	GetBannerColumns = []string{
		BannerIdColumnName,
		TitleColumnName,
		TextColumnName,
		UrlColumnName,
		FeatureIdColumnName,
		CreatedAtColumnName,
		UpdatedAtColumnName,
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
		CreatedAtColumnName,
		UpdatedAtColumnName,
		IsActiveColumnName,
	}
	InsertTagColumns = []string{
		BannerIdColumnName,
		TagIdColumnName,
	}
)
