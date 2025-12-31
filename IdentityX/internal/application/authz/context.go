package authz

import (
	"context"

	"GoAuth/internal/apierr"
)

type ctxKey string

const principalKey ctxKey = "principal"

func WithPrincipal(ctx context.Context, p *Principal) context.Context {
	return context.WithValue(ctx, principalKey, p)
}

func RequirePrincipal(ctx context.Context) (*Principal, error) {
	val := ctx.Value(principalKey)
	if val == nil {
		return nil, apierr.ErrUnauthorized.
			WithMsg("authentication required").
			WithID(apierr.AuthMissingPrincipal)
	}

	p, ok := val.(*Principal)
	if !ok {
		return nil, apierr.ErrInternal.
			WithMsg("invalid principal in context").
			WithID(apierr.AuthMissingPrincipal)
	}

	return p, nil
}
