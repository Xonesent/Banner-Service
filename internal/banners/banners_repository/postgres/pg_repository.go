package banners_postgres

import (
	"avito/assignment/config"
	trmsqlx "github.com/avito-tech/go-transaction-manager/sqlx"
	"github.com/jmoiron/sqlx"
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
