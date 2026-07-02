package telemetry

import (
	"context"
	"demo-rrweb/internal/config"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

// InitTracer mengonfigurasi mesin OTel dengan membaca data dari .env
func InitTracer() (*sdktrace.TracerProvider, error) {
	ctx := context.Background()

	// 1. Ambil nilai dari .env, sediakan nilai cadangan (fallback) jika kosong
	endpoint := config.GetEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
	serviceName := config.GetEnv("OTEL_SERVICE_NAME", "rrweb-golang-api")

	// 2. Gunakan variabel endpoint dari .env
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(endpoint),
	)
	if err != nil {
		return nil, err
	}

	// 3. Konfigurasi Provider dengan nama servis dari .env
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
		)),
	)

	otel.SetTracerProvider(tp)
	return tp, nil
}
