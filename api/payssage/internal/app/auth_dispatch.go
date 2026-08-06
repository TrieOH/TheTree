package app

import (
	"net/http"

	"lib/authz"
	spec "payssage"
)

// authResolver derives every operation's chain from the spec's security
// blocks, keyed by generated-form operationId.
func authResolver(mw middlewares) (*authz.Resolver, error) {
	return authz.NewResolver(spec.OpenAPISpec, authz.Primitives{
		JWT:    mw.jwtAuth,
		APIKey: mw.apiKeyAuth,
		Any:    mw.anyAuth,
	}, authz.Options{})
}

// resolveAuthChains resolves the spec-derived chains. Fails at boot when
// the spec and the middleware primitives disagree.
func resolveAuthChains(mw middlewares) (map[string][]func(http.Handler) http.Handler, error) {
	resolver, err := authResolver(mw)
	if err != nil {
		return nil, err
	}
	return resolver.Chains(), nil
}
