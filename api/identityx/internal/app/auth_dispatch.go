package app

import (
	"context"
	"net/http"
	"slices"
	"strings"

	spec "IdentityX"
	"IdentityX/internal/openapi"
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

// clientOnlyOps are the operations requiring a platform-level client
// identity (nil ProjectID). The scope requirement lives here until the
// authorize module absorbs it as a per-operation requirement.
var clientOnlyOps = []string{
	"listActors", "createActor", "getActor", "getActorByEmail",
	"createAPIKey",
	"createCapability",
	"listOrganizations", "createOrganization", "listOrganizationMembers",
	"addOrganizationMember", "removeOrganizationMember",
	"listOrganizationProjects", "createOrganizationProject",
	"listOrganizationProjectActors", "createOrganizationProjectActor",
	"listOrganizationProjectMembers", "addOrganizationProjectMember",
	"removeOrganizationProjectMember", "getOrganizationProjectActor",
	"listProjects", "createProject", "listProjectMembers",
	"addProjectMember", "removeProjectMember",
	"getPlatformProfile", "upsertPlatformProfile",
	"getProjectProfile", "upsertProjectProfile",
	"getPlatformProfileSchema", "upsertPlatformProfileSchema",
	"getProjectProfileSchema", "upsertProjectProfileSchema",
}

// authResolver derives every operation's chain from the spec's security
// blocks, keyed by spec-form operationId.
func authResolver(mw middlewares, guard func(http.Handler) http.Handler) (*authz.Resolver, error) {
	return authz.NewResolver(spec.OpenAPISpec, authz.Registry{
		"bearerAuth":            mw.jwtAuth,
		"apiKeyAuth":            mw.apiKeyAuth,
		"apiKeyAuth+bearerAuth": mw.anyAuth,
	}, authz.Options{
		SetupGuard:     guard,
		SkipSetupGuard: []string{"getSetup", "postSetup"},
	})
}

// resolveAuthChains composes the spec-derived chains with the platform
// client scope guard. Fails at boot when the spec and the middleware
// registry disagree.
func resolveAuthChains(mw middlewares) (map[string][]func(http.Handler) http.Handler, error) {
	resolver, err := authResolver(mw, setupGuard())
	if err != nil {
		return nil, err
	}
	chains := resolver.Chains()
	for _, op := range clientOnlyOps {
		chains[op] = append(chains[op], mw.clientOnly)
	}
	return chains, nil
}

// authDispatch is the strict-server middleware that resolves the auth
// chain for each operation (by operationId) and runs it around the
// openapi. Public operations (empty chain) pass through untouched.
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
