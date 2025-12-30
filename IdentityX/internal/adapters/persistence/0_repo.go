package persistence

import "go.opentelemetry.io/otel"

var (
	GoAuthRepoTracer = otel.Tracer("goauth/repo")
)
