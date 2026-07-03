package telemetry

import (
	"context"
	"demo-rrweb/internal/config"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

// Struct untuk menampung Trace dan Log
type TelemetryShutdown struct {
	Trace func(context.Context) error
	Log   func(context.Context) error
}

// InitTracer mengonfigurasi mesin OTel dengan membaca data dari .env
func InitTelemetry() (*TelemetryShutdown, error) {
	ctx := context.Background()

	// 1. Ambil nilai dari .env, sediakan fallback jika kosong
	endpoint := config.GetEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
	serviceName := config.GetEnv("OTEL_SERVICE_NAME", "rrweb-golang-api")

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(serviceName),
	)

	// Set up Tracing
	traceExporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(endpoint),
	)
	if err != nil {
		return nil, err
	}

	// Konfigurasi Provider trace
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)

	// Set up Logging
	logExporter, err := otlploggrpc.New(ctx,
		otlploggrpc.WithInsecure(),
		otlploggrpc.WithEndpoint(endpoint),
	)
	if err != nil {
		return nil, err
	}

	// Konfigurasi Provider log
	lp := log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(logExporter)),
		log.WithResource(res),
	)
	global.SetLoggerProvider(lp)

	return &TelemetryShutdown{
		Trace: tp.Shutdown,
		Log:   lp.Shutdown,
	}, nil
}
