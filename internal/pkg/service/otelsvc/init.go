package otelsvc

import (
	"context"
	"errors"
	"time"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TracerProviderType string
type MeterProviderType string

const (
	TracerProviderTypeOTLP   TracerProviderType = "otlp"
	TracerProviderTypeZipkin TracerProviderType = "zipkin"

	MeterProviderTypeOTLP   MeterProviderType = "otlp"
	MeterProviderTypeStdout MeterProviderType = "stdout"
)

type Service interface {
	Init() (err error)
	Shutdown()
}

type ServiceModule struct {
	opts          ServiceOpts
	shutdownFuncs []func(context.Context) error
}

type ServiceOpts struct {
	Enabled            bool
	ServiceName        string
	TracerProvider     TracerProviderType
	OtlpEndpoint       string
	ZipkinEndpoint     string
	SampleRate         float64
	MeterProvider      MeterProviderType
	OtlpMetricEndpoint string
	MetricInterval     time.Duration
}

func NewService(opts ServiceOpts) Service {
	return &ServiceModule{
		opts: opts,
	}
}

func (s *ServiceModule) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var err error
	for _, fn := range s.shutdownFuncs {
		err = errors.Join(err, fn(ctx))
	}
	s.shutdownFuncs = nil

	if err != nil {
		log.Fatal().Err(err).Msg("failed to shutdown otel")
	}
}

// Init bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func (s *ServiceModule) Init() (err error) {
	if !s.opts.Enabled {
		return nil
	}

	// Set up propagator.
	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	// Set up trace provider.
	tracerProvider, err := newTracerProvider(s.opts)
	if err != nil {
		s.Shutdown()
		return
	}
	if tracerProvider != nil {
		s.shutdownFuncs = append(s.shutdownFuncs, tracerProvider.Shutdown)
		otel.SetTracerProvider(tracerProvider)
	}

	// Set up meter provider.
	meterProvider, err := newMeterProvider(s.opts)
	if err != nil {
		s.Shutdown()
		return
	}
	if meterProvider != nil {
		s.shutdownFuncs = append(s.shutdownFuncs, meterProvider.Shutdown)
		otel.SetMeterProvider(meterProvider)
	}

	return
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newTracerProvider(opts ServiceOpts) (*trace.TracerProvider, error) {
	if opts.TracerProvider == "" {
		return nil, nil
	}
	var exporter trace.SpanExporter
	var err error

	switch opts.TracerProvider {
	case TracerProviderTypeOTLP:
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		conn, err_conn := grpc.DialContext(ctx, opts.OtlpEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
		if err_conn != nil {
			return nil, err_conn
		}
		exporter, err = otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	case TracerProviderTypeZipkin:
		exporter, err = zipkin.New(opts.ZipkinEndpoint)
	default:
		exporter, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
	}

	if err != nil {
		return nil, err
	}

	var sampler trace.Sampler
	if opts.SampleRate > 0 && opts.SampleRate <= 1 {
		sampler = trace.TraceIDRatioBased(opts.SampleRate)
	} else {
		sampler = trace.AlwaysSample()
	}

	resource, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(opts.ServiceName),
		),
	)

	if err != nil {
		return nil, err
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter, trace.WithBatchTimeout(time.Second)),
		trace.WithSampler(sampler),
		trace.WithResource(resource),
	)

	return tracerProvider, nil
}

func newMeterProvider(opts ServiceOpts) (*metric.MeterProvider, error) {
	if opts.MeterProvider == "" {
		return nil, nil
	}

	var metricExporter metric.Exporter
	var err error

	switch opts.MeterProvider {
	case MeterProviderTypeOTLP:
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		conn, err_conn := grpc.DialContext(ctx, opts.OtlpMetricEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
		if err_conn != nil {
			return nil, err_conn
		}
		metricExporter, err = otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
	case MeterProviderTypeStdout:
		metricExporter, err = stdoutmetric.New()
	default:
		metricExporter, err = stdoutmetric.New()
	}

	if err != nil {
		return nil, err
	}

	interval := 3 * time.Second
	if opts.MetricInterval > 0 {
		interval = opts.MetricInterval
	}

	resource, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(opts.ServiceName),
		),
	)
	if err != nil {
		return nil, err
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithResource(resource),
		metric.WithReader(
			metric.NewPeriodicReader(
				metricExporter,
				metric.WithInterval(interval),
			),
		),
	)
	return meterProvider, nil
}
