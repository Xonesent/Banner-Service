package sql_queries

const (
	BannersTableName         = "banner_schema.banners"
	BannersXTagsTableName    = "banner_schema.banners_X_tags"
	BannersVersionsTableName = "banner_schema.banners_versions"
	BannerIdColumnName       = "banner_id"
	TitleColumnName          = "title"
	TextColumnName           = "text"
	UrlColumnName            = "url"
	FeatureIdColumnName      = "feature_id"
	CreatedAtColumnName      = "created_at"
	UpdatedAtColumnName      = "updated_at"
	IsActiveColumnName       = "is_active"
	IdColumnName             = "id"
	TagIdColumnName          = "tag_id"
	TagIdsColumnName         = "tag_ids"
	VersionColumnName        = "version"
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
		VersionColumnName,
	}
	GetFullBannerColumns = []string{
		"b.banner_id",
		TitleColumnName,
		TextColumnName,
		UrlColumnName,
		TagIdColumnName,
		FeatureIdColumnName,
		CreatedAtColumnName,
		UpdatedAtColumnName,
		IsActiveColumnName,
		VersionColumnName,
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
		VersionColumnName,
	}
	SelectVersionColumns = []string{
		BannerIdColumnName,
		FeatureIdColumnName,
		TitleColumnName,
		TextColumnName,
		UrlColumnName,
		TagIdsColumnName,
		IsActiveColumnName,
		CreatedAtColumnName,
		UpdatedAtColumnName,
		VersionColumnName,
	}
	InsertBannerColumns = []string{
		TitleColumnName,
		TextColumnName,
		UrlColumnName,
		FeatureIdColumnName,
		CreatedAtColumnName,
		UpdatedAtColumnName,
		IsActiveColumnName,
		VersionColumnName,
	}
	InsertVersionColumns = []string{
		BannerIdColumnName,
		TitleColumnName,
		TextColumnName,
		UrlColumnName,
		FeatureIdColumnName,
		TagIdsColumnName,
		CreatedAtColumnName,
		UpdatedAtColumnName,
		IsActiveColumnName,
		VersionColumnName,
	}
	InsertTagColumns = []string{
		BannerIdColumnName,
		TagIdColumnName,
	}
)
