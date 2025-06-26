package middlewarerest

import (
	"encoding/json"
	"strings"

	"github.com/gofiber/fiber/v2"

	modeldto "gobase/internal/model/dto"
	"gobase/internal/pkg/helper"
)

func GetResponseMiddleware() fiber.Handler {
	return func(fc *fiber.Ctx) error {
		err := fc.Next()

		contentType := string(fc.Response().Header.ContentType())
		if strings.Contains(contentType, "application/json") {
			baseRes := &helper.JSONResponse{}
			if err := json.Unmarshal(fc.Response().Body(), baseRes); err != nil {
				return err
			}
			resStatus := modeldto.ResponseStatusDto{
				Success:        true,
				ResponseTimeMs: measuresResponseTimeMs(fc),
			}
			baseRes.ResponseStatus = resStatus
			err := fc.JSON(baseRes)
			if err != nil {
				return err
			}
		}

		return err
	}
}
