package tracing

import (
	"GoAuth/internal/application/authz"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func AnnotatePrincipal(span trace.Span, principal *authz.Principal) {
	if principal == nil {
		return
	}

	span.SetAttributes(
		attribute.String("user.id", principal.UserID.String()),
		attribute.String("user.session_id", principal.SessionID.String()),
		attribute.String("user.type", principal.UserType),
	)

	if principal.ProjectID != nil {
		span.SetAttributes(attribute.String("user.project_id", principal.ProjectID.String()))
	}
}
