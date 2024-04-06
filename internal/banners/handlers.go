package banners

import "github.com/gofiber/fiber/v2"

type Handlers interface {
	GetBanner() fiber.Handler
}
