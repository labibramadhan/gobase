package middlewarerest

import (
	"strconv"

	pkgErr "clodeo.tech/public/go-universe/pkg/err"
	"github.com/gofiber/fiber/v2"

	modeldto "gobase/internal/model/dto"
	"gobase/internal/pkg/helper"
)

func GetErrorMiddleware() fiber.ErrorHandler {
	return func(fc *fiber.Ctx, err error) error {
		responseTimeMs := measuresResponseTimeMs(fc)

		baseRes := helper.JSONResponse{}
		cusErr := pkgErr.GetError(err)
		message := cusErr.Error()
		if cusErr.Message != "" {
			message = cusErr.Message
		}
		resStatus := modeldto.ResponseStatusDto{
			Success:        false,
			ResponseTimeMs: responseTimeMs,
			ErrorMessage:   message,
			ErrorCode:      strconv.FormatInt(int64(cusErr.HTTPCode), 10),
		}
		baseRes.ResponseStatus = resStatus
		fc.Response().SetStatusCode(cusErr.HTTPCode)
		return fc.JSON(baseRes)
	}
}
