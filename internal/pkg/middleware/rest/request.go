package middlewarerest

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetRequestMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		startTime := time.Now()

		c.Locals("req_start_time", startTime)

		return c.Next()
	}
}
