package constant

import (
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

var (
	DevHosts = []string{
		"localhost:8892",
	}
	Host   string
	Tracer *tracesdk.TracerProvider
)
