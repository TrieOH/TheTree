package auth

import (
	"GoAuth/internal/domain/authz"
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type ctxKey string

const (
	principalKey ctxKey = "principal"
)

func WithPrincipal(ctx context.Context, p *authz.Principal) context.Context {
	return context.WithValue(ctx, principalKey, p)
}

func RequirePrincipal(ctx context.Context) (*authz.Principal, error) {
	val := ctx.Value(principalKey)
	if val == nil {
		return nil, authz.ErrMissingPrincipal{}
	}

	p, ok := val.(*authz.Principal)
	if !ok {
		return nil, authz.ErrPrincipalMissingInContext{}
	}

	return p, nil
}

func RequirePrincipalAndAnnotate(ctx context.Context, span trace.Span) (*authz.Principal, error) {
	var principal *authz.Principal
	principal, err := RequirePrincipal(ctx)
	if err != nil {
		return nil, err
	}

	AnnotatePrincipal(span, principal)

	return principal, nil
}

// AnnotatePrincipal annotates a span with the principal's information.
func AnnotatePrincipal(span trace.Span, principal *authz.Principal) {
	if span == nil || principal == nil {
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
