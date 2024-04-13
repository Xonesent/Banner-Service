package error_handler

import (
	"avito/assignment/pkg/constant"
	"avito/assignment/pkg/traces"
	"avito/assignment/pkg/utilities"
	"errors"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func FiberErrorHandler(ctx *fiber.Ctx, err error) error {
	setStatusCode(ctx, err)
	traceId, ok := ctx.Locals(traces.TraceIDHeader).(string)
	if !ok {
		traceId = "impossible to get traceId"
	}
	errorInfo := strings.Split(err.Error(), "|")
	if utilities.InStringSlice(constant.Host, constant.DevHosts) {
		return ctx.JSON(map[string]interface{}{
			"data":        nil,
			"trace-id":    traceId,
			"error_place": errorInfo[0],
			"error_value": errorInfo[1],
			"error_args":  errorInfo[2],
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
	}
	ctx.Status(statusCode)
}
