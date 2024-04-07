package banners_postgres

import (
	"avito/assignment/config"
	"avito/assignment/internal/models/banner_models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	trmsqlx "github.com/avito-tech/go-transaction-manager/sqlx"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel"
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

func (b *BannersRepo) GetPossibleBannerIds(ctx context.Context, tagId int) ([]int, error) {
	ctx, span := otel.Tracer("").Start(ctx, "BannersRepo.GetPossibleBannerIds")
	defer span.End()

	query, args, err := sq.Select(BannerIdColumnName).
		From(BannersXTagsTableName).
		Where(
			sq.And{
				sq.Eq{TagIdColumnName: tagId},
			},
		).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var bannerIds []int

	tr := b.txGetter.DefaultTrOrDB(ctx, b.db)
	err = tr.SelectContext(
		ctx,
		&bannerIds,
		query,
		args...,
	)
	if err != nil {
		return nil, err
	}

	return bannerIds, nil
}

func (b *BannersRepo) GetPossibleTagIds(ctx context.Context, bannerId int) ([]int, error) {
	ctx, span := otel.Tracer("").Start(ctx, "BannersRepo.GetPossibleBannerIds")
	defer span.End()

	query, args, err := sq.Select(TagIdColumnName).
		From(BannersXTagsTableName).
		Where(
			sq.And{
				sq.Eq{BannerIdColumnName: bannerId},
			},
		).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var tagIds []int

	tr := b.txGetter.DefaultTrOrDB(ctx, b.db)
	err = tr.SelectContext(
		ctx,
		&tagIds,
		query,
		args...,
	)
	if err != nil {
		return nil, err
	}

	return tagIds, nil
}

func (b *BannersRepo) GetBanner(ctx context.Context, params banner_models.GetPostgresBanner) (*banner_models.BannerContent, error) {
	ctx, span := otel.Tracer("").Start(ctx, "BannersRepo.GetPossibleBannerIds")
	defer span.End()

	query, args, err := sq.Select(GetBannerColumns...).
		From(BannersTableName).
		Where(
			sq.And{
				sq.Eq{BannerIdColumnName: params.PossibleBannerIds},
				sq.Eq{FeatureIdColumnName: params.FeatureId},
			},
		).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var bannerIds banner_models.BannerContent

	tr := b.txGetter.DefaultTrOrDB(ctx, b.db)
	err = tr.GetContext(
		ctx,
		&bannerIds,
		query,
		args...,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fiber.ErrNotFound
		}
		return nil, err
	}

	return &bannerIds, nil
}

func (b *BannersRepo) SelectBanner(ctx context.Context, params banner_models.SelectPostgresBanner) (*[]banner_models.FullBannerContent, error) {
	ctx, span := otel.Tracer("").Start(ctx, "BannersRepo.GetPossibleBannerIds")
	defer span.End()

	var conditions sq.And
	if params.TagId != 0 {
		conditions = append(conditions, sq.Eq{BannerIdColumnName: params.PossibleBannerIds})
	}
	if params.FeatureId != 0 {
		conditions = append(conditions, sq.Eq{FeatureIdColumnName: params.FeatureId})
	}

	queryBuilder := sq.Select(SelectBannerColumns...).
		From(BannersTableName).
		Where(conditions).
		PlaceholderFormat(sq.Dollar).
		Offset(uint64(params.Offset))

	if params.Limit != 0 {
		queryBuilder = queryBuilder.Limit(uint64(params.Limit))
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var bannerContents []banner_models.FullBannerContent

	tr := b.txGetter.DefaultTrOrDB(ctx, b.db)
	err = tr.SelectContext(
		ctx,
		&bannerContents,
		query,
		args...,
	)
	if err != nil {
		return nil, err
	}

	return &bannerContents, nil
}

func (b *BannersRepo) AddBanner(ctx context.Context, params banner_models.AddBanner) (*int, error) {
	ctx, span := otel.Tracer("").Start(ctx, "BannersRepo.GetPossibleBannerIds")
	defer span.End()

	query, args, err := sq.Insert(BannersTableName).
		Columns(InsertBannerColumns...).
		Values(
			params.Content.Title,
			params.Content.Text,
			params.Content.Url,
			params.FeatureId,
			params.IsActive,
			time.Now(),
			time.Now(),
		).
		Suffix("RETURNING banner_id").
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	var bannerId int

	err = b.db.QueryRowContext(ctx, query, args...).Scan(&bannerId)
	if err != nil {
		return nil, err
	}

	return &bannerId, nil
}

func (b *BannersRepo) CheckExist(ctx context.Context, tagId int, featureId int) error {
	ctx, span := otel.Tracer("").Start(ctx, "BannersRepo.GetPossibleBannerIds")
	defer span.End()

	query, _, err := sq.Select("1").
		From(fmt.Sprintf("%s b", BannersTableName)).
		InnerJoin(fmt.Sprintf("%s bxt ON bxt.banner_id = b.banner_id", BannersXTagsTableName)).
		Where(
			sq.And{
				sq.Eq{TagIdColumnName: tagId},
				sq.Eq{FeatureIdColumnName: featureId},
			},
		).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return err
	}

	var exists bool

	tr := b.txGetter.DefaultTrOrDB(ctx, b.db)
	if err = tr.QueryRowxContext(ctx, query, tagId, featureId).Scan(&exists); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if exists {
		return errors.New(fmt.Sprintf("impossible to add, banner with tagId = %d and featureId = %d already exists", tagId, featureId))
	}
	return nil
}

func (b *BannersRepo) AddTags(ctx context.Context, bannerId int, tagId int) error {
	ctx, span := otel.Tracer("").Start(ctx, "BannersRepo.GetPossibleBannerIds")
	defer span.End()

	query, args, err := sq.Insert(BannersXTagsTableName).
		Columns(InsertTagColumns...).
		Values(
			bannerId,
			tagId,
		).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return err
	}

	tr := b.txGetter.DefaultTrOrDB(ctx, b.db)
	if _, err = tr.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}
