package app

import (
	"context"

	"IdentityX/internal/services"
	"IdentityX/internal/tokens"
	"IdentityX/models"

	mws "github.com/MintzyG/fun/middlewares"
)

// SetupAuthMiddlewares wires the JWT and API-key authentication
// middlewares. Both cross the Token-lifecycle module and the api_keys
// feature's Authenticate seam — neither touches the repos directly, so
// verification logic lives where the concept lives and this factory stays
// a thin binding.
func (app *IdentityX) SetupAuthMiddlewares(
	tokensMgr *tokens.Manager,
	ops *services.Operations,
) *mws.Middleware[*models.AccessClaims] {
	keyFunc := func(ctx context.Context, tokenStr string) (*models.AccessClaims, error) {
		claims := &models.AccessClaims{}
		err := tokensMgr.Verify(ctx, tokenStr, claims)
		if err != nil {
			return nil, err
		}
		return claims, nil
	}

	jwtHook := func(ctx context.Context, claims *models.AccessClaims) (context.Context, error) {
		identity := &models.Identity{
			Sub: models.SubjectFromAccessSub(&claims.Sub),
			Cred: models.Credential{
				Type: models.TokenCredentialType,
			},
		}
		return models.WithIdentity(ctx, identity), nil
	}

	apiKeyHook := func(ctx context.Context, rawKey string) (context.Context, error) {
		identity, err := ops.APIKeys.Authenticate(ctx, rawKey)
		if err != nil {
			return nil, err
		}
		return models.WithIdentity(ctx, identity), nil
	}
	return mws.New[*models.AccessClaims](keyFunc, jwtHook, apiKeyHook)
}
