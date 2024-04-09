package error_handler

import (
	"avito/assignment/pkg/constant"
	"avito/assignment/pkg/traces"
	"avito/assignment/pkg/utilities"
	"errors"
	"github.com/gofiber/fiber/v2"
)

func FiberErrorHandler(ctx *fiber.Ctx, err error) error {
	setStatusCode(ctx, err)
	traceId, ok := ctx.Locals(traces.TraceIDHeader).(string)
	if !ok {
		traceId = "impossible to get traceId"
	}
	if utilities.InStringSlice(constant.Host, constant.DevHosts) {
		return ctx.JSON(map[string]interface{}{
			"error":    err.Error(),
			"data":     nil,
			"trace-id": traceId,
		})
	}
	return ctx.JSON(map[string]interface{}{
		"data":     nil,
		"trace-id": traceId,
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
