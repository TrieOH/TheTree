package app

import (
	"context"

	"IdentityX/internal/authz"
	"IdentityX/internal/services"
	"IdentityX/internal/tokens"
	"IdentityX/models"
	libauthz "lib/authz"

	mws "github.com/MintzyG/fun/middlewares"
)

// authPrimitives builds the primitives IdentityX's spec-derived chains
// resolve against. IdentityX issues the tokens, so it verifies them itself
// instead of listing the kit's identityx module. The scope checkers come
// from the Access-check module: adding a scope is one spec x-scope line plus
// one entry in authz.ScopeCheckers.
func authPrimitives(tokensMgr *tokens.Manager, ops *services.Operations, authzSvc *authz.Service) libauthz.Primitives {
	authMW := authMiddlewares(tokensMgr, ops)
	return libauthz.Primitives{
		JWT:    authMW.JWT(),
		APIKey: authMW.APIKey(),
		Any:    authMW.AnyAuth(),
		Scopes: authzSvc.ScopeCheckers(),
	}
}

// authMiddlewares wires the JWT and API-key authentication middlewares.
// Both cross the Token-lifecycle module and the api_keys feature's
// Authenticate seam — neither touches the repos directly, so verification
// logic lives where the concept lives and this factory stays a thin binding.
func authMiddlewares(tokensMgr *tokens.Manager, ops *services.Operations) *mws.Middleware[*models.AccessClaims] {
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
