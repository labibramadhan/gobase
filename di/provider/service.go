package provider

import (
	"time"

	"clodeo.tech/public/go-universe/pkg/localization"
	"github.com/google/wire"

	"gobase/config"
	"gobase/di/registry"
	"gobase/internal/pkg/service/otelsvc"
	"gobase/internal/pkg/service/structprocessor"
)

var ServiceSet = wire.NewSet(
	ProvideServiceStructProcessorService,
	ProvideServiceOtelService,
)

func ProvideServiceStructProcessorService(localizer localization.Localizer) structprocessor.StructProcessorService {
	return structprocessor.NewStructProcessorService(structprocessor.StructProcessorServiceModuleOpts{
		Localizer: localizer,
	})
}

func ProvideServiceOtelService(cfg *config.MainConfig) (otelsvc.Service, registry.CleanupFunc) {
	svc := otelsvc.NewService(otelsvc.ServiceOpts{
		Enabled:            cfg.Otel.Enabled,
		ServiceName:        cfg.General.ServiceName,
		TracerProvider:     otelsvc.TracerProviderType(cfg.Otel.TracerProvider),
		OtlpEndpoint:       cfg.Otel.OtlpEndpoint,
		ZipkinEndpoint:     cfg.Otel.ZipkinEndpoint,
		SampleRate:         cfg.Otel.SampleRate,
		MeterProvider:      otelsvc.MeterProviderType(cfg.Otel.MeterProvider),
		OtlpMetricEndpoint: cfg.Otel.OtlpMetricEndpoint,
		MetricInterval:     time.Duration(cfg.Otel.MetricIntervalMs) * time.Millisecond,
	})

	return svc, svc.Shutdown
}
