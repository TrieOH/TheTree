package authz

import (
	"IdentityX/internal/shared/errx"
	"context"

	"github.com/MintzyG/fail/v3"
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
		return nil, fail.New(errx.AuthPrincipalNotInContext).RecordCtx(ctx)
	}

	p, ok := val.(*Principal)
	if !ok {
		return nil, fail.New(errx.AuthInvalidPrincipal).RecordCtx(ctx)
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
