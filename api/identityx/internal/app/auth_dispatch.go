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

// clientOnlyOps are the operations requiring a platform-level client
// identity (nil ProjectID). The list is validated against the spec at
// resolver construction (an unknown operationId fails startup), and the
// guard is applied by the resolver itself.
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
	"getPlatformProfileSchema", "upsertPlatformProfileSchema",
	"getProjectProfileSchema", "upsertProjectProfileSchema",
	"listOutdatedPlatformProfiles",
}

// authResolver derives every operation's chain from the spec's security
// blocks, keyed by generated-form operationId. The setup guard, the
// platform-client scope guard, and their op lists are validated against the
// spec at construction; a mismatch fails boot, never production.
func authResolver(mw middlewares, guard func(http.Handler) http.Handler) (*authz.Resolver, error) {
	return authz.NewResolver(spec.OpenAPISpec, authz.Registry{
		"bearerAuth":            mw.jwtAuth,
		"apiKeyAuth":            mw.apiKeyAuth,
		"apiKeyAuth+bearerAuth": mw.anyAuth,
	}, authz.Options{
		SetupGuard:     guard,
		SkipSetupGuard: []string{"getSetup", "postSetup"},
		ClientOnly:     mw.clientOnly,
		ClientOnlyOps:  clientOnlyOps,
	})
}

// resolveAuthChains composes the spec-derived chains. Fails at boot when
// the spec and the middleware registry disagree.
func resolveAuthChains(mw middlewares) (map[string][]func(http.Handler) http.Handler, error) {
	resolver, err := authResolver(mw, setupGuard())
	if err != nil {
		return nil, err
	}
	return resolver.Chains(), nil
}
