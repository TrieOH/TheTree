package app

import (
	"net/http"

	spec "IdentityX"
	"lib/authz"
	"lib/globals"

	"github.com/MintzyG/fun"
)

// setupGuard returns the middleware that gates every operation (except the
// two /auth/setup routes, which manage the flag themselves) until setup has
// completed.
func setupGuard() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !globals.SetupComplete() {
				fun.ServiceUnavailable("please setup IDX first on /auth/setup").Send(w)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// authResolver derives every operation's chain from the spec's security
// blocks, keyed by generated-form operationId. The setup guard and its op
// list are validated against the spec at construction; a mismatch fails
// boot, never production. Platform-vs-project scope is not a chain concern:
// it is enforced per handler via handlers.RequireClientOnly (platform-only
// operations) and handlers.RequireProjectClientOnly (project-scoped
// operations, not used yet).
func authResolver(mw middlewares, guard func(http.Handler) http.Handler) (*authz.Resolver, error) {
	return authz.NewResolver(spec.OpenAPISpec, authz.Primitives{
		JWT:    mw.jwtAuth,
		APIKey: mw.apiKeyAuth,
		Any:    mw.anyAuth,
	}, authz.Options{
		SetupGuard:     guard,
		SkipSetupGuard: []string{"getSetup", "postSetup"},
	})
}

// resolveAuthChains composes the spec-derived chains. Fails at boot when
// the spec and the middleware primitives disagree.
func resolveAuthChains(mw middlewares) (map[string][]func(http.Handler) http.Handler, error) {
	resolver, err := authResolver(mw, setupGuard())
	if err != nil {
		return nil, err
	}
	return resolver.Chains(), nil
}
