package transportrest

import (
	"context"
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
	srv        *fiber.App
	cfg        *config.MainConfig
	restRouter *registry.RESTRouter
}

func NewTransport(cfg *config.MainConfig, restRouter *registry.RESTRouter) (registry.IApplicationTransportREST, registry.CleanupFunc) {
	transportModule := &TransportModule{
		cfg:        cfg,
		restRouter: restRouter,
	}

	return transportModule, transportModule.Cleanup
}

func (m *TransportModule) Cleanup() {
	if m.srv == nil {
		return
	}
	err := m.srv.Shutdown()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to shutdown fiber server")
	}
}

func (m *TransportModule) Run(ctx context.Context) error {
	restConfig := m.cfg.Server.Rest
	readTimeoutSecond := restConfig.ReadTimeoutSecond
	writeTimeoutSecond := restConfig.WriteTimeoutSecond
	bodyLimitMB := restConfig.BodyLimitMB

	jsoniterExtra.SetNamingStrategy(jsoniterExtra.LowerCaseWithUnderscores)
	jsonHandler := jsoniter.ConfigCompatibleWithStandardLibrary

	m.srv = fiber.New(fiber.Config{
		ErrorHandler: middlewarerest.GetErrorMiddleware(),
		BodyLimit:    bodyLimitMB * 1024 * 1024,
		ReadTimeout:  time.Duration(readTimeoutSecond) * time.Second,
		WriteTimeout: time.Duration(writeTimeoutSecond) * time.Second,
		JSONEncoder:  jsonHandler.Marshal,
		JSONDecoder:  jsonHandler.Unmarshal,
	})

	// Liveness check api
	m.srv.Get("/health", func(fc *fiber.Ctx) error {
		return fc.JSON("OK")
	})

	m.srv.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))
	m.srv.Use(otelfiber.Middleware())
	m.srv.Use(logger.New())
	m.srv.Use(cors.New())
	m.srv.Use(helmet.New())

	err := m.srv.Listen(m.cfg.Server.Rest.ListenAddress + ":" + strconv.Itoa(m.cfg.Server.Rest.Port))
	if err != nil {
		return err
	}

	return err
}
