package banners_postgres

import (
	"avito/assignment/config"
	"avito/assignment/internal/models/banner_models"
	"context"
	"database/sql"
	"errors"
	sq "github.com/Masterminds/squirrel"
	trmsqlx "github.com/avito-tech/go-transaction-manager/sqlx"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel"
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
