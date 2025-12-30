package tracing

import (
	"GoAuth/internal/domain/auth"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func AnnotateAccessClaims(span trace.Span, claims *auth.AccessClaims) {
	span.SetAttributes(
		attribute.String("user.id", claims.Sub.ID.String()),
		attribute.String("user.session_id", claims.Sub.SessionID.String()),
		attribute.String("user.type", claims.Sub.UserType),
	)

	if claims.Sub.ProjectID != nil {
		span.SetAttributes(attribute.String("user.project_id", claims.Sub.ProjectID.String()))
	}
}
