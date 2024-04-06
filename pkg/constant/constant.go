package constant

import (
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

// server constants
var (
	DevHosts = []string{
		"localhost:8892",
	}
	Host   string
	Tracer *tracesdk.TracerProvider
)

// role constants
var (
	AllRoles   = []string{UserToken, AdminToken}
	AdminRoles = []string{AdminToken}
	UserToken  = "user_token"
	AdminToken = "admin_token"
)
