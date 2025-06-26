package transportrest

import (
	"strconv"
	"time"

	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	jsoniter "github.com/json-iterator/go"
	jsoniterExtra "github.com/json-iterator/go/extra"
	"github.com/rs/zerolog/log"

	"gobase/config"
	"gobase/di/registry"
	middlewarerest "gobase/internal/pkg/middleware/rest"
)

type TransportModule struct {
	cfg        *config.MainConfig
	restRouter *registry.RESTRouter
}

func NewTransport(cfg *config.MainConfig, restRouter *registry.RESTRouter) registry.IApplicationTransportREST {
	return &TransportModule{
		cfg:        cfg,
		restRouter: restRouter,
	}
}

func (m *TransportModule) Run() (registry.CleanupFunc, error) {
	restConfig := m.cfg.Server.Rest
	readTimeoutSecond := restConfig.ReadTimeoutSecond
	writeTimeoutSecond := restConfig.WriteTimeoutSecond
	bodyLimitMB := restConfig.BodyLimitMB

	jsoniterExtra.SetNamingStrategy(jsoniterExtra.LowerCaseWithUnderscores)
	jsonHandler := jsoniter.ConfigCompatibleWithStandardLibrary

	srv := fiber.New(fiber.Config{
		ErrorHandler: middlewarerest.GetErrorMiddleware(),
		BodyLimit:    bodyLimitMB * 1024 * 1024,
		ReadTimeout:  time.Duration(readTimeoutSecond) * time.Second,
		WriteTimeout: time.Duration(writeTimeoutSecond) * time.Second,
		JSONEncoder:  jsonHandler.Marshal,
		JSONDecoder:  jsonHandler.Unmarshal,
	})

	// Liveness check api
	srv.Get("/health", func(fc *fiber.Ctx) error {
		return fc.JSON("OK")
	})

	srv.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))
	srv.Use(otelfiber.Middleware())
	srv.Use(logger.New())
	srv.Use(cors.New())
	srv.Use(helmet.New())

	cleanup := func() {
		err := srv.Shutdown()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to shutdown fiber server")
		}
	}

	err := srv.Listen(m.cfg.Server.Rest.ListenAddress + ":" + strconv.Itoa(m.cfg.Server.Rest.Port))
	if err != nil {
		return nil, err
	}

	return cleanup, err
}
