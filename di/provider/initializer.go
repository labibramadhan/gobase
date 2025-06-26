package provider

import (
	env "clodeo.tech/public/go-universe/pkg/env"
	"clodeo.tech/public/go-universe/pkg/tracer"
	"github.com/rs/zerolog/log"

	"gobase/config"
	"gobase/di/registry"
	iconfig "gobase/internal/config"
)

func InitializeTracer(cfg *config.MainConfig) {
	if cfg.Tracer.Enabled {
		tracerConfig := &tracer.TracerConfig{
			Provider:    cfg.Tracer.Provider,
			Environment: env.GetEnvironmentName(),
			ServiceName: cfg.ServiceName,
		}

		if cfg.Tracer.Provider == "jaeger" {
			tracerConfig.JaegerCollectorURL = cfg.Tracer.Jaeger.CollectorUrl
		}

		err := tracer.Init(tracerConfig)

		if err != nil {
			log.Fatal().Err(err).Msg("failed to initialize tracer")
		}
	}
}

func InitializeStructConverter() {
	err := iconfig.InitStructConverter()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize struct converter")
	}
}

func Initializer(
	cfg *config.MainConfig,
) registry.InitializerFunc {
	return func() {
		InitializeTracer(cfg)
		InitializeStructConverter()
	}
}
