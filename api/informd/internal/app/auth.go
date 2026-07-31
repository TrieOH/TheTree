package app

import (
	"context"

	idx "sdk/identityx"

	"github.com/MintzyG/fun"
	fm "github.com/MintzyG/fun/middlewares"
)

func (app *Informd) setupAuthMiddlewares() *fm.Middleware[*idx.AccessClaims] {
	keyFunc := func(ctx context.Context, tokenStr string) (*idx.AccessClaims, error) {
		return app.idxClient.Tokens.VerifyAccessToken(ctx, tokenStr)
	}

	jwtHook := func(ctx context.Context, claims *idx.AccessClaims) (context.Context, error) {
		return idx.WithIdentity(ctx, &idx.Identity{
			Sub: idx.Subject{
				ID:           claims.Sub.ID,
				ProjectID:    claims.Sub.ProjectID,
				Email:        claims.Sub.Email,
				Type:         claims.Sub.Type,
				Capabilities: claims.Sub.Capabilities,
				Metadata:     claims.Sub.Metadata,
			},
			Cred: idx.Credential{
				Type: "token",
			},
		}), nil
	}

	apiKeyHook := func(_ context.Context, _ string) (context.Context, error) {
		return nil, fun.ErrNotImplemented("api keys are not yet supported")
	}

	return fm.New[*idx.AccessClaims](keyFunc, jwtHook, apiKeyHook)
}
