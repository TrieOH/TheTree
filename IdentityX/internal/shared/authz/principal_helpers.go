package authz

import (
	"context"

	"github.com/MintzyG/fun"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type ctxKey string

const (
	principalKey ctxKey = "principal"
)

func WithPrincipal(ctx context.Context, p *Principal) context.Context {
	return context.WithValue(ctx, principalKey, p)
}

func RequirePrincipal(ctx context.Context) (*Principal, error) {
	val := ctx.Value(principalKey)
	if val == nil {
		return nil, fun.ErrUnauthorized("missing principal in context")
	}
	p, ok := val.(*Principal)
	if !ok {
		return nil, fun.Errf("invalid principal type: %T", val).Unauthorized()
	}
	return p, nil
}

func RequirePrincipalAndAnnotate(ctx context.Context, span trace.Span) (*Principal, error) {
	var principal *Principal
	principal, err := RequirePrincipal(ctx)
	if err != nil {
		return nil, err
	}
	AnnotatePrincipal(span, principal)
	return principal, nil
}

// AnnotatePrincipal annotates a span with the principal's information.
func AnnotatePrincipal(span trace.Span, principal *Principal) {
	if span == nil || principal == nil {
		return
	}
	span.SetAttributes(attribute.String("user.id", principal.UserID.String()))
	if principal.ProjectID != nil {
		span.SetAttributes(attribute.String("user.project_id", principal.ProjectID.String()))
	}
}
