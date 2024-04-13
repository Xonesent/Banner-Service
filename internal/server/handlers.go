package server

import (
	banners_http "avito/assignment/internal/banners/banners_delivery/http"
	banners_postgres "avito/assignment/internal/banners/banners_repository"
	"avito/assignment/internal/banners/banners_usecase"
	"avito/assignment/internal/middleware"
	trmsqlx "github.com/avito-tech/go-transaction-manager/sqlx"
	"github.com/avito-tech/go-transaction-manager/trm/manager"
)

func (s *Server) MapHandlers() (err error) {
	bannersPGRepo := banners_postgres.NewBannerRepository(s.cfg, s.pgDB, trmsqlx.DefaultCtxGetter)
	bannersRedisRepo := banners_postgres.NewClientRedisRepository(s.redis, s.cfg)
	trManager := manager.Must(trmsqlx.NewDefaultFactory(s.pgDB))

	bannersUC := banners_usecase.NewBannersUC(s.cfg, trManager, bannersPGRepo, bannersRedisRepo)

	bannersHandlers := banners_http.NewUserHandler(bannersUC, s.cfg)
	mw := middleware.NewOfficiantMiddleware(s.cfg)

	bannersGroup := s.fiber.Group("")

	banners_http.MapBannersRoutes(bannersGroup, bannersHandlers, mw)

	return nil
}
