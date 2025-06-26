package middlewarerest

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func GetLimiterMiddleware() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        20,
		Expiration: 1 * time.Minute,
	})
}
