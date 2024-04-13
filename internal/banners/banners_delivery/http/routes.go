package banners_http

import (
	"avito/assignment/internal/middleware"
	"avito/assignment/pkg/constant"
	"github.com/gofiber/fiber/v2"
)

func MapBannersRoutes(group fiber.Router, h Handlers, mw *middleware.MDWManager) {
	group.Get("/user_banner", mw.CheckAuthToken(constant.AllRoles), h.GetBanner())
	group.Get("/banner", mw.CheckAuthToken(constant.AdminRoles), h.GetManyBanner())
	group.Post("/banner", mw.CheckAuthToken(constant.AdminRoles), h.AddBanner())
	group.Patch("/banner/:banner_id", mw.CheckAuthToken(constant.AdminRoles), h.PatchBanner())
	group.Delete("/banner/:banner_id", mw.CheckAuthToken(constant.AdminRoles), h.DeleteBanner())
}
