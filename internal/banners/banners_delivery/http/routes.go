package banners_http

import (
	"avito/assignment/internal/banners"
	"avito/assignment/internal/middleware"
	"avito/assignment/pkg/constant"
	"github.com/gofiber/fiber/v2"
)

func MapBannersRoutes(group fiber.Router, h banners.Handlers, mw *middleware.MDWManager) {
	group.Get("/user_banner", mw.CheckAuthToken(constant.AllRoles), h.GetBanner())
	group.Get("/banner", mw.CheckAuthToken(constant.AdminRoles))
	group.Post("/banner", mw.CheckAuthToken(constant.AdminRoles))
	group.Patch("/banner/:banner_id", mw.CheckAuthToken(constant.AdminRoles))
	group.Delete("/banner/:banner_id", mw.CheckAuthToken(constant.AdminRoles))
}
