package error_handler

import (
	"avito/assignment/pkg/constant"
	"avito/assignment/pkg/utilities"
	"errors"
	"github.com/gofiber/fiber/v2"
)

func FiberErrorHandler(ctx *fiber.Ctx, err error) error {
	setStatusCode(ctx, err)
	if utilities.InStringSlice(constant.Host, constant.DevHosts) {
		return ctx.JSON(map[string]interface{}{
			"error": err.Error(),
			"data":  nil,
		})
	}
	return ctx.JSON(map[string]interface{}{
		"data": nil,
	})
}

func setStatusCode(ctx *fiber.Ctx, err error) {
	statusCode := 0
	var e *fiber.Error
	if errors.As(err, &e) {
		statusCode = e.Code
	} else {
		statusCode = fiber.StatusInternalServerError
	}
	ctx.Status(statusCode)
}
