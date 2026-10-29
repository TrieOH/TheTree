package app

import (
	"context"

	"IdentityX/internal/authz"
	"IdentityX/internal/features/api_keys"
	"IdentityX/internal/tokens"
	"IdentityX/models"
	libauthz "lib/authz"
	"lib/httpserver"
)

// authPrimitives builds the primitives IdentityX's spec-derived chains
// resolve against. IdentityX issues the tokens, so it verifies them itself
// instead of listing the kit's identityx module. The scope checkers come
// from the Access-check module: adding a scope is one spec x-scope line plus
// one entry in authz.ScopeCheckers.
func authPrimitives(tokensMgr *tokens.Manager, apiKeys *api_keys.Operations, authzSvc *authz.Service) libauthz.Primitives {
	p := authenticators(tokensMgr, apiKeys).Primitives()
	p.Scopes = authzSvc.ScopeCheckers()
	return p
}

// authenticators wires JWT and API-key authentication. Both cross the
// Token-lifecycle module and the api_keys feature's Authenticate seam —
// neither touches the repos directly, so verification logic lives where the
// concept lives and this factory stays a thin binding.
func authenticators(tokensMgr *tokens.Manager, apiKeys *api_keys.Operations) httpserver.Authenticators {
	jwtAuth := func(ctx context.Context, token string) (context.Context, error) {
		claims := &models.AccessClaims{}
		err := tokensMgr.Verify(ctx, token, claims)
		if err != nil {
			return nil, err
		}
		identity := &models.Identity{
			Sub: models.SubjectFromAccessSub(&claims.Sub),
			Cred: models.Credential{
				Type: models.TokenCredentialType,
			},
		}
		return models.WithIdentity(ctx, identity), nil
	}

	apiKeyAuth := func(ctx context.Context, rawKey string) (context.Context, error) {
		identity, err := apiKeys.Authenticate(ctx, rawKey)
		if err != nil {
			return nil, err
		}
		return models.WithIdentity(ctx, identity), nil
	}
	return httpserver.Authenticators{JWT: jwtAuth, APIKey: apiKeyAuth}
}
