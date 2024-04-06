package traces

import (
	"avito/assignment/config"
	"avito/assignment/pkg/constant"
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"log"
)

func NewTracer(ctx context.Context, cfg config.Config) {
	exp, err := jaeger.New(
		jaeger.WithCollectorEndpoint(
			jaeger.WithEndpoint(cfg.OpenTelemetry.URL),
		),
	)
	if err != nil {
		log.Fatalf("Cannot create Jaeger exporter: %s", err.Error())
	}

	constant.Tracer = tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.OpenTelemetry.ServiceName),
		)),
	)
	otel.SetTracerProvider(constant.Tracer)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{}, propagation.Baggage{},
		),
	)
}
