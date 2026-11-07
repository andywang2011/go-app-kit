package telemetry

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// Config controls OpenTelemetry export for the process.
type Config struct {
	Enabled  bool
	Endpoint string // OTLP HTTP endpoint (host:port or full URL)
	// Insecure allows plaintext (non-TLS) OTLP export. It defaults to false, so
	// exports use TLS; an http:// endpoint or a plaintext collector requires
	// Insecure: true.
	Insecure    bool
	ServiceName string
	SampleRatio float64
}

// ConfigDefault returns safe defaults (telemetry off).
func ConfigDefault() Config {
	return Config{
		Enabled:     false,
		Endpoint:    "localhost:4318",
		ServiceName: "",
		SampleRatio: 1.0,
	}
}

// LoadConfig reads app.otel.* from the active viper config.
func LoadConfig() Config {
	cfg := ConfigDefault()
	cfg.Enabled = viper.GetBool("app.otel.enabled")
	cfg.Endpoint = strings.TrimSpace(viper.GetString("app.otel.endpoint"))
	if cfg.Endpoint == "" {
		cfg.Endpoint = ConfigDefault().Endpoint
	}
	cfg.Insecure = viper.GetBool("app.otel.insecure")
	cfg.ServiceName = strings.TrimSpace(viper.GetString("app.otel.service_name"))
	if r := viper.GetFloat64("app.otel.sample_ratio"); r > 0 {
		cfg.SampleRatio = r
	}
	return cfg
}

// exporterOptions builds OTLP HTTP options. TLS is the default; plaintext is
// used only when cfg.Insecure is explicitly true.
func exporterOptions(cfg Config) ([]otlptracehttp.Option, error) {
	var opts []otlptracehttp.Option
	switch {
	case strings.HasPrefix(cfg.Endpoint, "http://"):
		if !cfg.Insecure {
			return nil, errors.New("telemetry: plaintext http:// endpoint requires Insecure: true")
		}
		opts = append(opts, otlptracehttp.WithEndpointURL(cfg.Endpoint))
	case strings.HasPrefix(cfg.Endpoint, "https://"):
		opts = append(opts, otlptracehttp.WithEndpointURL(cfg.Endpoint))
	default:
		opts = append(opts, otlptracehttp.WithEndpoint(cfg.Endpoint))
		if cfg.Insecure {
			opts = append(opts, otlptracehttp.WithInsecure())
		}
	}
	return opts, nil
}

// Setup configures the global OTel tracer provider when cfg.Enabled is true.
func Setup(ctx context.Context, cfg Config) (shutdown func(context.Context) error, err error) {
	noop := func(context.Context) error { return nil }
	if !cfg.Enabled {
		return noop, nil
	}
	if cfg.ServiceName == "" {
		return noop, errors.New("telemetry: service_name required when enabled")
	}

	opts, err := exporterOptions(cfg)
	if err != nil {
		return noop, err
	}

	exporter, err := otlptracehttp.New(ctx, opts...)
	if err != nil {
		return noop, fmt.Errorf("telemetry: otlp exporter: %w", err)
	}

	ratio := cfg.SampleRatio
	if ratio <= 0 {
		ratio = 1.0
	}
	if ratio > 1 {
		ratio = 1
	}

	// NewSchemaless avoids Schema URL conflicts with resource.Default(), which
	// tracks the SDK's current schema while semconv imports may lag.
	res, err := resource.Merge(
		resource.Default(),
		resource.NewSchemaless(semconv.ServiceName(cfg.ServiceName)),
	)
	if err != nil {
		return noop, fmt.Errorf("telemetry: resource: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(ratio)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tp.Shutdown, nil
}
