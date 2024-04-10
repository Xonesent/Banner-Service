package banners_repository

import (
	"avito/assignment/config"
	"avito/assignment/internal/models"
	"avito/assignment/internal/store/sql_queries"
	"context"
	"database/sql"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	trmsqlx "github.com/avito-tech/go-transaction-manager/sqlx"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel"
	"sort"
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

func (b *BannersRepo) GetPossibleBannerIds(ctx context.Context, tagId models.TagId) ([]models.BannerId, error) {
	ctx, span := otel.Tracer("").Start(ctx, "BannersRepo.GetPossibleBannerIds")
	defer span.End()

	query, args, err := sq.Select(sql_queries.BannerIdColumnName).
		From(sql_queries.BannersXTagsTableName).
		Where(sq.Eq{sql_queries.TagIdColumnName: tagId}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("BannersRepo.GetPossibleBannerIds.Select; err = %s", err.Error()))
	}

	var bannerIds []models.BannerId

	tr := b.txGetter.DefaultTrOrDB(ctx, b.db)
	err = tr.SelectContext(ctx, &bannerIds, query, args...)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("BannersRepo.GetPossibleBannerIds.SelectContext; err = %s", err.Error()))
	}

	return bannerIds, nil
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
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("BannersRepo.GetPossibleTagIds.Select; err = %s", err.Error()))
	}

	var tagIds []models.TagId

	tr := b.txGetter.DefaultTrOrDB(ctx, b.db)
	err = tr.SelectContext(ctx, &tagIds, query, args...)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("BannersRepo.GetPossibleTagIds.SelectContext; err = %s", err.Error()))
	}

	return tagIds, nil
}

func (b *BannersRepo) GetBannerPostgres(ctx context.Context, getPostgresqlBannerParams *GetPostgresBanner) (*models.Banner, error) {
	ctx, span := otel.Tracer("").Start(ctx, "BannersRepo.GetBannerPostgres")
	defer span.End()

	query, args, err := sq.Select(sql_queries.GetBannerColumnsWithInnerJoin...).
		From(fmt.Sprintf("%s b", sql_queries.BannersTableName)).
		InnerJoin(fmt.Sprintf("%s bxt ON bxt.banner_id = b.banner_id", sql_queries.BannersXTagsTableName)).
		Where(
			sq.And{
				sq.Eq{sql_queries.TagIdColumnName: getPostgresqlBannerParams.TagId},
				sq.Eq{sql_queries.FeatureIdColumnName: getPostgresqlBannerParams.FeatureId},
			},
		).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("BannersRepo.GetBannerPostgres.SelectContext; err = %s", err.Error()))
	}

	var banner models.Banner

	tr := b.txGetter.DefaultTrOrDB(ctx, b.db)
	err = tr.GetContext(ctx, &banner, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("BannersRepo.GetBannerPostgres.GetContext; err = %s", err.Error()))
		}
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("BannersRepo.GetBannerPostgres.SelectContext; err = %s", err.Error()))
	}

	return &banner, nil
}

func (b *BannersRepo) GetManyBannerPostgres(ctx context.Context, getManyPostgresBannerParams *GetManyPostgresBanner) (*[]models.Banner, error) {
	ctx, span := otel.Tracer("").Start(ctx, "BannersRepo.GetManyBannerPostgres")
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
		OrderBy(sql_queries.BannerIdColumnName).
		Offset(uint64(*getManyPostgresBannerParams.Offset))

	if getManyPostgresBannerParams.Limit != nil && *getManyPostgresBannerParams.Limit != 0 {
		queryBuilder = queryBuilder.Limit(uint64(*getManyPostgresBannerParams.Limit))
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("BannersRepo.GetManyBannerPostgres.Select; err = %s", err.Error()))
	}

	var banners []models.Banner

	tr := b.txGetter.DefaultTrOrDB(ctx, b.db)
	err = tr.SelectContext(ctx, &banners, query, args...)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("BannersRepo.GetManyBannerPostgres.SelectContext; err = %s", err.Error()))
	}

	return &banners, nil
}

func (b *BannersRepo) GetManyPossibleTagIds(ctx context.Context, bannerIds []models.BannerId, manyBanner *[]models.Banner) (*[]models.FullBanner, error) {
	ctx, span := otel.Tracer("").Start(ctx, "BannersRepo.GetPossibleTagIds")
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
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("BannersRepo.GetPossibleTagIds.Select; err = %s", err.Error()))
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
			return nil, err
		}

		if bannerId != bannerIds[currIndex] {
			currIndex++
			manyFullBanner[currIndex] = *(*manyBanner)[currIndex].ToFullBannerWithoutTagIds()
		}
		manyFullBanner[currIndex].TagIds = append(manyFullBanner[currIndex].TagIds, tagId)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &manyFullBanner, nil
}

func (b *BannersRepo) AddBannerPostgres(ctx context.Context, addPostgresBannerParams *AddPostgresBanner) (*GetInsertParams, error) {
	ctx, span := otel.Tracer("").Start(ctx, "BannersRepo.AddBannerPostgres")
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
		).
		Suffix("RETURNING banner_id,created_at,updated_at").
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("BannersRepo.AddBannerPostgres.Insert; err = %s", err.Error()))
	}

	var bannerId models.BannerId
	var createdAt time.Time
	var updatedAt time.Time

	tr := b.txGetter.DefaultTrOrDB(ctx, b.db)
	err = tr.QueryRowContext(ctx, query, args...).Scan(&bannerId, &createdAt, &updatedAt)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("BannersRepo.AddBannerPostgres.Scan; err = %s", err.Error()))
	}

	return &GetInsertParams{
		BannerId:  bannerId,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

func (b *BannersRepo) CheckExist(ctx context.Context, checkExistBannerParams *CheckExistBanner) (*[]ExistBanner, error) {
	ctx, span := otel.Tracer("").Start(ctx, "BannersRepo.CheckExist")
	defer span.End()

	query, args, err := sq.Select(sql_queries.TagIdColumnName, sql_queries.FeatureIdColumnName).
		From(fmt.Sprintf("%s b", sql_queries.BannersTableName)).
		InnerJoin(fmt.Sprintf("%s bxt ON bxt.banner_id = b.banner_id", sql_queries.BannersXTagsTableName)).
		Where(
			sq.And{
				sq.Eq{sql_queries.TagIdColumnName: checkExistBannerParams.TagId},
				sq.Eq{sql_queries.FeatureIdColumnName: checkExistBannerParams.FeatureId},
			},
		).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("BannersRepo.CheckExist.Scan; err = %s", err.Error()))
	}

	var existParams []ExistBanner

	tr := b.txGetter.DefaultTrOrDB(ctx, b.db)
	if err = tr.SelectContext(ctx, &existParams, query, args...); err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("BannersRepo.CheckExist.Scan; err = %s", err.Error()))
	}
	return &existParams, nil
}

func (b *BannersRepo) AddTags(ctx context.Context, addTagsPostgresParams *AddTagsPostgres) error {
	ctx, span := otel.Tracer("").Start(ctx, "BannersRepo.AddTags")
	defer span.End()

	sqlBuilder := sq.Insert(sql_queries.BannersXTagsTableName).Columns(sql_queries.InsertTagColumns...)

	for _, tagId := range addTagsPostgresParams.TagId {
		sqlBuilder = sqlBuilder.Values(addTagsPostgresParams.BannerId, tagId)
	}

	query, args, err := sqlBuilder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("BannersRepo.AddTags.Insert; err = %s", err.Error()))
	}

	tr := b.txGetter.DefaultTrOrDB(ctx, b.db)
	if _, err = tr.ExecContext(ctx, query, args...); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("BannersRepo.AddTags.ExecContext; err = %s", err.Error()))
	}

	return nil
}
