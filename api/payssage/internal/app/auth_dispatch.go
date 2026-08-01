package app

import (
	"context"
	"net/http"
	"slices"
	"strings"

	"lib/authz"
	spec "payssage"
	"payssage/internal/openapi"
)

// authResolver derives every operation's chain from the spec's security
// blocks, keyed by spec-form operationId.
func authResolver(mw middlewares) (*authz.Resolver, error) {
	return authz.NewResolver(spec.OpenAPISpec, authz.Registry{
		"bearerAuth": mw.jwtAuth,
	}, authz.Options{})
}

// resolveAuthChains resolves the spec-derived chains. Fails at boot when
// the spec and the middleware registry disagree.
func resolveAuthChains(mw middlewares) (map[string][]func(http.Handler) http.Handler, error) {
	resolver, err := authResolver(mw)
	if err != nil {
		return nil, err
	}
	return resolver.Chains(), nil
}

// authDispatch is the strict-server middleware that resolves the auth
// chain for each operation (by operationId) and runs it around the
// handler. Public operations (empty chain) pass through untouched.
func authDispatch(chains map[string][]func(http.Handler) http.Handler) openapi.StrictMiddlewareFunc {
	return func(f openapi.StrictHandlerFunc, operationID string) openapi.StrictHandlerFunc {
		if operationID != "" {
			operationID = strings.ToLower(operationID[:1]) + operationID[1:]
		}
		chain := chains[operationID]
		if len(chain) == 0 {
			return f
		}
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (any, error) {
			var resp any
			var ferr error
			var called bool
			var next http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called = true
				resp, ferr = f(r.Context(), w, r, request)
			})
			for i := range slices.Backward(chain) {
				next = chain[i](next)
			}
			next.ServeHTTP(w, r)
			if !called {
				// An auth middleware rejected the request and already
				// wrote the response; nothing to return.
				return nil, nil
			}
			return resp, ferr
		}
	}
}
