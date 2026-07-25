package app

import (
	"context"

	"lib/telemetry"
	idx "sdk/identityx"

	"github.com/MintzyG/fun"
	mws "github.com/MintzyG/fun/middlewares"
	"go.uber.org/zap"
)

func SetupAuthMiddlewares() *mws.Middleware[*idx.AccessClaims] {
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

	apiKeyHook := func(ctx context.Context, rawKey string) (context.Context, error) {
		telemetry.Log().Info("user tried to use api key",
			zap.String("message", "this service does not provide a public api"),
			zap.String("key", rawKey),
		)
		return ctx, fun.ErrForbidden("this service does not provide public access to the api")
	}

	return mws.New[*idx.AccessClaims](keyFunc, jwtHook, apiKeyHook)
}
