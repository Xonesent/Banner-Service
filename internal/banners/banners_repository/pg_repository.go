package banners_repository

import (
	"avito/assignment/config"
	"avito/assignment/internal/models"
	"avito/assignment/internal/store/sql_queries"
	"avito/assignment/pkg/errlst"
	"avito/assignment/pkg/traces"
	"context"
	"database/sql"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	trmsqlx "github.com/avito-tech/go-transaction-manager/sqlx"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel"
	"sort"
	"strings"
	"time"
)

type BannersRepo struct {
	cfg      *config.Config
	db       *sqlx.DB
	txGetter *trmsqlx.CtxGetter
}

func NewBannerRepository(cfg *config.Config, db *sqlx.DB, txGetter *trmsqlx.CtxGetter) *BannersRepo {
	return &BannersRepo{
		cfg:      cfg,
		db:       db,
		txGetter: txGetter,
	}
}

func (b *BannersRepo) GetPossibleTagIds(ctx context.Context, bannerId models.BannerId) ([]models.TagId, error) {
	ctx, span := otel.Tracer("").Start(ctx, "BannersRepo.GetPossibleTagIds")
	defer span.End()

	query, args, err := sq.Select(sql_queries.TagIdColumnName).
		From(sql_queries.BannersXTagsTableName).
		Where(sq.Eq{sql_queries.BannerIdColumnName: bannerId}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, traces.SpanSetErrWrap(span, errlst.HttpServerError, err, "BannersRepo.GetPossibleTagIds.Select")
	}

	var tagIds []models.TagId

	tr := b.txGetter.DefaultTrOrDB(ctx, b.db)
	err = tr.SelectContext(ctx, &tagIds, query, args...)
	if err != nil {
		return nil, traces.SpanSetErrWrap(span, errlst.HttpServerError, err, "BannersRepo.GetPossibleTagIds.SelectContext")
	}

	return tagIds, nil
}

func (b *BannersRepo) GetBanner(ctx context.Context, featureId models.FeatureId, tagId models.TagId) (*models.Banner, error) {
	ctx, span := otel.Tracer("").Start(ctx, "BannersRepo.GetBanner")
	defer span.End()

	query, args, err := sq.Select(sql_queries.GetBannerColumnsWithInnerJoin...).
		From(fmt.Sprintf("%s b", sql_queries.BannersTableName)).
		InnerJoin(fmt.Sprintf("%s bxt ON bxt.banner_id = b.banner_id", sql_queries.BannersXTagsTableName)).
		Where(
			sq.And{
				sq.Eq{sql_queries.TagIdColumnName: tagId},
				sq.Eq{sql_queries.FeatureIdColumnName: featureId},
			},
		).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, traces.SpanSetErrWrap(span, errlst.HttpServerError, err, "BannersRepo.GetBanner.Select")
	}

	var banner models.Banner

	tr := b.txGetter.DefaultTrOrDB(ctx, b.db)
	err = tr.GetContext(ctx, &banner, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, traces.SpanSetErrWrap(span, errlst.HttpErrNotFound, err, "BannersRepo.GetBanner.ErrNoRows")
		}
		return nil, traces.SpanSetErrWrap(span, errlst.HttpServerError, err, "BannersRepo.GetBanner.GetContext")
	}

	return &banner, nil
}

func (b *BannersRepo) GetManyBanner(ctx context.Context, getManyPostgresBannerParams *GetManyPostgresBanner) (*[]models.Banner, error) {
	ctx, span := otel.Tracer("").Start(ctx, "BannersRepo.GetManyBanner")
	defer span.End()

	var conditions sq.And
	var sqlBuilder sq.SelectBuilder
	if getManyPostgresBannerParams.TagId != nil {
		conditions = append(conditions, sq.Eq{sql_queries.TagIdColumnName: getManyPostgresBannerParams.TagId})
		sqlBuilder = sq.Select(sql_queries.GetBannerColumnsWithInnerJoin...).
			From(fmt.Sprintf("%s b", sql_queries.BannersTableName)).
			InnerJoin(fmt.Sprintf("%s bxt ON bxt.banner_id = b.banner_id", sql_queries.BannersXTagsTableName))
	} else {
		sqlBuilder = sq.Select(sql_queries.SelectBannerColumns...).
			From(sql_queries.BannersTableName)
	}
	if getManyPostgresBannerParams.FeatureId != nil {
		conditions = append(conditions, sq.Eq{sql_queries.FeatureIdColumnName: getManyPostgresBannerParams.FeatureId})
	}

	queryBuilder := sqlBuilder.Where(conditions).
		PlaceholderFormat(sq.Dollar).
		OrderBy(sql_queries.BannerIdColumnName)

	if getManyPostgresBannerParams.Offset != nil {
		queryBuilder = queryBuilder.Offset(uint64(*getManyPostgresBannerParams.Offset))
	}

	if getManyPostgresBannerParams.Limit != nil && *getManyPostgresBannerParams.Limit != 0 {
		queryBuilder = queryBuilder.Limit(uint64(*getManyPostgresBannerParams.Limit))
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, traces.SpanSetErrWrap(span, errlst.HttpServerError, err, "BannersRepo.GetManyBanner.Select")
	}

	var banners []models.Banner

	tr := b.txGetter.DefaultTrOrDB(ctx, b.db)
	err = tr.SelectContext(ctx, &banners, query, args...)
	if err != nil {
		return nil, traces.SpanSetErrWrap(span, errlst.HttpServerError, err, "BannersRepo.GetManyBanner.SelectContext")
	}

	return &banners, nil
}

func (b *BannersRepo) GetManyPossibleTagIds(ctx context.Context, bannerIds []models.BannerId, manyBanner *[]models.Banner) (*[]models.FullBanner, error) {
	ctx, span := otel.Tracer("").Start(ctx, "BannersRepo.GetManyPossibleTagIds")
	defer span.End()

	manyFullBanner := make([]models.FullBanner, len(*manyBanner))
	manyFullBanner[0] = *(*manyBanner)[0].ToFullBannerWithoutTagIds()

	sort.SliceStable(bannerIds, func(i, j int) bool {
		return bannerIds[i] < bannerIds[j]
	})

	query, args, err := sq.Select(sql_queries.BannerIdColumnName, sql_queries.TagIdColumnName).
		From(sql_queries.BannersXTagsTableName).
		Where(sq.Eq{sql_queries.BannerIdColumnName: bannerIds}).
		PlaceholderFormat(sq.Dollar).
		OrderBy(sql_queries.BannerIdColumnName).
		ToSql()
	if err != nil {
		return nil, traces.SpanSetErrWrap(span, errlst.HttpServerError, err, "BannersRepo.GetManyPossibleTagIds.Select")
	}

	tr := b.txGetter.DefaultTrOrDB(ctx, b.db)
	rows, err := tr.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var currIndex int = 0

	for rows.Next() {
		var bannerId models.BannerId
		var tagId models.TagId
		if err := rows.Scan(&bannerId, &tagId); err != nil {
			return nil, traces.SpanSetErrWrap(span, errlst.HttpServerError, err, "BannersRepo.GetManyPossibleTagIds.Scan")
		}

		if bannerId != bannerIds[currIndex] {
			currIndex++
			manyFullBanner[currIndex] = *(*manyBanner)[currIndex].ToFullBannerWithoutTagIds()
		}
		manyFullBanner[currIndex].TagIds = append(manyFullBanner[currIndex].TagIds, tagId)
	}

	if err := rows.Err(); err != nil {
		return nil, traces.SpanSetErrWrap(span, errlst.HttpServerError, err, "BannersRepo.GetManyPossibleTagIds.rows.Err")
	}

	return &manyFullBanner, nil
}

func (b *BannersRepo) AddBanner(ctx context.Context, addPostgresBannerParams *AddPostgresBanner) (*GetInsertParams, error) {
	ctx, span := otel.Tracer("").Start(ctx, "BannersRepo.AddBanner")
	defer span.End()

	query, args, err := sq.Insert(sql_queries.BannersTableName).
		Columns(sql_queries.InsertBannerColumns...).
		Values(
			addPostgresBannerParams.Title,
			addPostgresBannerParams.Text,
			addPostgresBannerParams.Url,
			addPostgresBannerParams.FeatureId,
			time.Now(),
			time.Now(),
			addPostgresBannerParams.IsActive,
			1,
		).
		Suffix("RETURNING banner_id,created_at,updated_at").
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, traces.SpanSetErrWrap(span, errlst.HttpServerError, err, "BannersRepo.AddBanner.Insert")
	}

	var bannerId models.BannerId
	var createdAt time.Time
	var updatedAt time.Time

	tr := b.txGetter.DefaultTrOrDB(ctx, b.db)
	err = tr.QueryRowContext(ctx, query, args...).Scan(&bannerId, &createdAt, &updatedAt)
	if err != nil {
		return nil, traces.SpanSetErrWrap(span, errlst.HttpServerError, err, "BannersRepo.AddBanner.QueryRowContext")
	}

	return &GetInsertParams{
		BannerId:  bannerId,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

func (b *BannersRepo) CheckExist(ctx context.Context, tagIds []models.TagId, featureId models.FeatureId) (*[]ExistBanner, error) {
	ctx, span := otel.Tracer("").Start(ctx, "BannersRepo.CheckExist")
	defer span.End()

	query, args, err := sq.Select(sql_queries.TagIdColumnName, sql_queries.FeatureIdColumnName).
		From(fmt.Sprintf("%s b", sql_queries.BannersTableName)).
		InnerJoin(fmt.Sprintf("%s bxt ON bxt.banner_id = b.banner_id", sql_queries.BannersXTagsTableName)).
		Where(
			sq.And{
				sq.Eq{sql_queries.TagIdColumnName: tagIds},
				sq.Eq{sql_queries.FeatureIdColumnName: featureId},
			},
		).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, traces.SpanSetErrWrap(span, errlst.HttpServerError, err, "BannersRepo.CheckExist.Select")
	}

	var existParams []ExistBanner

	tr := b.txGetter.DefaultTrOrDB(ctx, b.db)
	if err = tr.SelectContext(ctx, &existParams, query, args...); err != nil {
		return nil, traces.SpanSetErrWrap(span, errlst.HttpServerError, err, "BannersRepo.CheckExist.SelectContext")
	}
	return &existParams, nil
}

func (b *BannersRepo) AddTags(ctx context.Context, tagIds []models.TagId, bannerId models.BannerId) error {
	ctx, span := otel.Tracer("").Start(ctx, "BannersRepo.AddTags")
	defer span.End()

	sqlBuilder := sq.Insert(sql_queries.BannersXTagsTableName).Columns(sql_queries.InsertTagColumns...)
	for _, tagId := range tagIds {
		sqlBuilder = sqlBuilder.Values(bannerId, tagId)
	}

	query, args, err := sqlBuilder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return traces.SpanSetErrWrap(span, errlst.HttpServerError, err, "BannersRepo.AddTags.Insert")
	}

	tr := b.txGetter.DefaultTrOrDB(ctx, b.db)
	if _, err = tr.ExecContext(ctx, query, args...); err != nil {
		return traces.SpanSetErrWrap(span, errlst.HttpServerError, err, "BannersRepo.AddTags.ExecContext")
	}

	return nil
}

func (b *BannersRepo) GetBannerById(ctx context.Context, bannerId models.BannerId) (*models.FullBanner, error) {
	ctx, span := otel.Tracer("").Start(ctx, "BannersRepo.GetBannerById")
	defer span.End()

	fullBanner := &models.FullBanner{}

	query, args, err := sq.Select(sql_queries.GetFullBannerColumns...).
		From(fmt.Sprintf("%s b", sql_queries.BannersTableName)).
		InnerJoin(fmt.Sprintf("%s bxt ON bxt.banner_id = b.banner_id", sql_queries.BannersXTagsTableName)).
		Where(
			sq.And{
				sq.Eq{"b." + sql_queries.BannerIdColumnName: bannerId},
			},
		).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, traces.SpanSetErrWrap(span, errlst.HttpServerError, err, "BannersRepo.GetBannerById.Select")
	}

	var banner models.Banner
	var firstRow bool = true

	tr := b.txGetter.DefaultTrOrDB(ctx, b.db)
	rows, err := tr.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.StructScan(&banner); err != nil {
			return nil, err
		}
		if firstRow {
			fullBanner = banner.ToFullBanner([]models.TagId{banner.TagId})
			firstRow = false
		} else {
			fullBanner.TagIds = append(fullBanner.TagIds, banner.TagId)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, traces.SpanSetErrWrap(span, errlst.HttpServerError, err, "BannersRepo.GetBannerById.rows.Err")
	}

	return fullBanner, nil
}

func (b *BannersRepo) UpdateBannerById(ctx context.Context, updateBannerByIdParams *UpdateBannerById) error {
	ctx, span := otel.Tracer("").Start(ctx, "BannersRepo.UpdateBannerById")
	defer span.End()

	sqlBuilder := sq.Update(sql_queries.BannersTableName)
	sqlBuilder = filter(sqlBuilder, updateBannerByIdParams)

	query, args, err := sqlBuilder.
		Where(sq.Eq{sql_queries.BannerIdColumnName: updateBannerByIdParams.BannerId}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return err
	}

	tr := b.txGetter.DefaultTrOrDB(ctx, b.db)

	_, err = tr.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func filter(sqlBuilder sq.UpdateBuilder, updateBannerByIdParams *UpdateBannerById) sq.UpdateBuilder {
	if updateBannerByIdParams.FeatureId != nil {
		sqlBuilder = sqlBuilder.Set(sql_queries.FeatureIdColumnName, updateBannerByIdParams.FeatureId)
	}
	if updateBannerByIdParams.Title != nil {
		sqlBuilder = sqlBuilder.Set(sql_queries.TitleColumnName, updateBannerByIdParams.Title)
	}
	if updateBannerByIdParams.Text != nil {
		sqlBuilder = sqlBuilder.Set(sql_queries.TextColumnName, updateBannerByIdParams.Text)
	}
	if updateBannerByIdParams.Url != nil {
		sqlBuilder = sqlBuilder.Set(sql_queries.UrlColumnName, updateBannerByIdParams.Url)
	}
	if updateBannerByIdParams.IsActive != nil {
		sqlBuilder = sqlBuilder.Set(sql_queries.IsActiveColumnName, updateBannerByIdParams.IsActive)
	}
	sqlBuilder = sqlBuilder.Set(sql_queries.UpdatedAtColumnName, time.Now())
	sqlBuilder = sqlBuilder.Set(sql_queries.VersionColumnName, updateBannerByIdParams.Version+1)

	return sqlBuilder
}

func (b *BannersRepo) AddVersion(ctx context.Context, prevBanner *models.FullBanner) error {
	ctx, span := otel.Tracer("").Start(ctx, "BannersRepo.AddVersion")
	defer span.End()

	var tagIdsStr string
	for _, id := range prevBanner.TagIds {
		tagIdsStr += fmt.Sprintf("%d,", id)
	}
	tagIdsStr = "{" + strings.TrimSuffix(tagIdsStr, ",") + "}"

	query, args, err := sq.Insert(sql_queries.BannersVersionsTableName).
		Columns(sql_queries.InsertVersionColumns...).
		Values(
			prevBanner.BannerId,
			prevBanner.Content.Title,
			prevBanner.Content.Text,
			prevBanner.Content.Url,
			prevBanner.FeatureId,
			tagIdsStr,
			prevBanner.CreatedAt,
			prevBanner.UpdatedAt,
			prevBanner.IsActive,
			prevBanner.Version,
		).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return traces.SpanSetErrWrap(span, errlst.HttpServerError, err, "BannersRepo.AddVersion.Insert")
	}

	tr := b.txGetter.DefaultTrOrDB(ctx, b.db)
	if _, err = tr.ExecContext(ctx, query, args...); err != nil {
		return traces.SpanSetErrWrap(span, errlst.HttpServerError, err, "BannersRepo.AddVersion.ExecContext")
	}

	return nil
}

func (b *BannersRepo) DeleteTags(ctx context.Context, tagIds []models.TagId, bannerId models.BannerId) error {
	ctx, span := otel.Tracer("").Start(ctx, "BannersRepo.DeleteTags")
	defer span.End()

	var conditions sq.And
	conditions = append(conditions, sq.Eq{sql_queries.BannerIdColumnName: bannerId})
	if len(tagIds) != 0 {
		conditions = append(conditions, sq.Eq{sql_queries.TagIdColumnName: tagIds})
	}

	query, args, err := sq.Delete(sql_queries.BannersXTagsTableName).
		Where(conditions).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return traces.SpanSetErrWrap(span, errlst.HttpServerError, err, "BannersRepo.DeleteTags.Delete")
	}

	tr := b.txGetter.DefaultTrOrDB(ctx, b.db)
	if _, err = tr.ExecContext(ctx, query, args...); err != nil {
		return traces.SpanSetErrWrap(span, errlst.HttpServerError, err, "BannersRepo.DeleteTags.ExecContext")
	}

	return nil
}

func (b *BannersRepo) DeleteVersion(ctx context.Context, versions []int64, bannerId models.BannerId) error {
	ctx, span := otel.Tracer("").Start(ctx, "BannersRepo.DeleteVersion")
	defer span.End()

	var conditions sq.And
	conditions = append(conditions, sq.Eq{sql_queries.BannerIdColumnName: bannerId})
	if len(versions) != 0 {
		conditions = append(conditions, sq.Eq{sql_queries.VersionColumnName: versions})
	}

	query, args, err := sq.Delete(sql_queries.BannersVersionsTableName).
		Where(conditions).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return traces.SpanSetErrWrap(span, errlst.HttpServerError, err, "BannersRepo.DeleteTags.Delete")
	}

	tr := b.txGetter.DefaultTrOrDB(ctx, b.db)
	if _, err = tr.ExecContext(ctx, query, args...); err != nil {
		return traces.SpanSetErrWrap(span, errlst.HttpServerError, err, "BannersRepo.DeleteTags.ExecContext")
	}

	return nil
}

func (b *BannersRepo) DeleteBannerById(ctx context.Context, bannerId models.BannerId) error {
	ctx, span := otel.Tracer("").Start(ctx, "BannersRepo.DeleteBannerById")
	defer span.End()

	query, args, err := sq.Delete(sql_queries.BannersTableName).
		Where(sq.Eq{sql_queries.BannerIdColumnName: bannerId}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return traces.SpanSetErrWrap(span, errlst.HttpServerError, err, "BannersRepo.DeleteBannerById.Delete")
	}

	tr := b.txGetter.DefaultTrOrDB(ctx, b.db)
	if _, err = tr.ExecContext(ctx, query, args...); err != nil {
		return traces.SpanSetErrWrap(span, errlst.HttpServerError, err, "BannersRepo.DeleteBannerById.ExecContext")
	}

	return nil
}

func (b *BannersRepo) GetBannerVersions(ctx context.Context, bannerId models.BannerId, versions []int64) (*[]models.FullBanner, error) {
	ctx, span := otel.Tracer("").Start(ctx, "BannersRepo.GetPossibleTagIds")
	defer span.End()

	sqlBuilder := sq.Select(sql_queries.SelectVersionColumns...).
		From(sql_queries.BannersVersionsTableName)

	conditions := sq.And{}
	conditions = append(conditions, sq.Eq{sql_queries.BannerIdColumnName: bannerId})
	if len(versions) != 0 {
		conditions = append(conditions, sq.Eq{sql_queries.VersionColumnName: versions})
	}

	query, args, err := sqlBuilder.Where(conditions).
		PlaceholderFormat(sq.Dollar).
		OrderBy(fmt.Sprintf("%s DESC", sql_queries.VersionColumnName)).
		ToSql()
	if err != nil {
		return nil, traces.SpanSetErrWrap(span, errlst.HttpServerError, err, "BannersRepo.GetBannerVersions.Select")
	}

	var banners []FullBanner

	tr := b.txGetter.DefaultTrOrDB(ctx, b.db)
	err = tr.SelectContext(ctx, &banners, query, args...)
	if err != nil {
		return nil, traces.SpanSetErrWrap(span, errlst.HttpServerError, err, "BannersRepo.GetBannerVersions.SelectContext")
	}

	fullBanners := make([]models.FullBanner, len(banners))
	for i := range banners {
		fullBanners[i] = banners[i].ToFullBanners()
	}

	return &fullBanners, nil
}
