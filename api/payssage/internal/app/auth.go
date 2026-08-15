package app

import (
	"context"
	"lib/telemetry"

	idx "sdk/identityx"

	"github.com/MintzyG/fun"

	fm "github.com/MintzyG/fun/middlewares"
	"go.uber.org/zap"
)

func (app *Payssage) setupAuthMiddlewares() *fm.Middleware[*idx.AccessClaims] {
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
		// identityx returns the fun envelope ({code, data, ...}); the
		// identity lives under `data` — SetResult(&ident) directly would
		// leave it zeroed (owner_id 000 in wallets) because the envelope
		// has no top-level `subject`.
		var envelope struct {
			Data idx.Identity `json:"data"`
		}
		resp, err := app.httpClient.R().
			WithContext(ctx).
			SetHeader("X-API-KEY", rawKey).
			SetResult(&envelope).
			Get("http://identityx:8080/auth/introspect")
		if err != nil {
			telemetry.Log().Error("error fetching identity", zap.Error(err))
			return ctx, err
		}
		if resp.IsStatusFailure() {
			telemetry.Log().Error("introspect failed", zap.Int("status", resp.StatusCode()))
			return ctx, fun.ErrForbidden("invalid api key")
		}
		ident := envelope.Data
		return idx.WithIdentity(ctx, &ident), nil
	}

	return fm.New[*idx.AccessClaims](keyFunc, jwtHook, apiKeyHook)
}
