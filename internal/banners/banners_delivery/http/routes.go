package banners_http

import (
	"avito/assignment/internal/banners"
	"avito/assignment/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func MapBannersRoutes(group fiber.Router, h banners.Handlers, mw *middleware.MDWManager) {
	//group.Post()
}
