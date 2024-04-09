package middleware

import (
	"avito/assignment/config"
	"avito/assignment/pkg/traces"
	"avito/assignment/pkg/utilities"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type MDWManager struct {
	cfg *config.Config
}

func NewOfficiantMiddleware(cfg *config.Config) *MDWManager {
	return &MDWManager{cfg: cfg}
}

func (m *MDWManager) CheckAuthToken(restrictions []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		_, span := traces.StartFiberTrace(c, "MDWManager.CheckAuthToken")
		defer span.End()

		token := c.Get("token")
		if token == "" {
			return fiber.NewError(fiber.StatusUnauthorized, fmt.Sprintf("MDWManager.CheckAuthToken.NilToken"))
		} else if !utilities.InStringSlice(token, restrictions) {
			return fiber.NewError(fiber.StatusForbidden, fmt.Sprintf("MDWManager.CheckAuthToken.Forbidden"))
		}

		c.Locals("token", token)

		return c.Next()
	}
}
