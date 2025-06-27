package provider

import (
	"github.com/rs/zerolog/log"

	"gobase/config"
	"gobase/di/registry"
	iconfig "gobase/internal/config"
	"gobase/internal/pkg/service/otelsvc"
)

func InitializeOtel(otelsvc otelsvc.Service) {
	err := otelsvc.Init()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize otel")
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
	otel otelsvc.Service,
) registry.InitializerFunc {
	return func() {
		InitializeStructConverter()
		InitializeOtel(otel)
	}
}
