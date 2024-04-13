package errlst

import (
	"github.com/gofiber/fiber/v2"
)

var (
	HttpErrUnauthorized   = fiber.NewError(fiber.StatusUnauthorized, fiber.ErrUnauthorized.Error())
	HttpErrInvalidRequest = fiber.NewError(fiber.StatusBadRequest, fiber.ErrBadRequest.Error())
	HttpErrNotFound       = fiber.NewError(fiber.StatusNotFound, fiber.ErrNotFound.Error())
	HttpErrForbidden      = fiber.NewError(fiber.StatusForbidden, fiber.ErrForbidden.Error())
	HttpServerError       = fiber.NewError(fiber.StatusInternalServerError, fiber.ErrInternalServerError.Error())
)
