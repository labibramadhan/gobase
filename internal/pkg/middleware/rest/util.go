package middlewarerest

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func measuresResponseTimeMs(fc *fiber.Ctx) int64 {
	var (
		responseTimeMs int64 = 0
	)
	// validate request time from context if exist calculate the latency
	if startTime, ok := fc.Locals("req_start_time").(time.Time); ok {
		// measures the time elapsed from startTime to the present.
		responseTimeMs = time.Since(startTime).Milliseconds()
	}
	return responseTimeMs
}
