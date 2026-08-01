package app

import (
	"context"
	"net/http"
	"slices"

	"IdentityX/internal/handler"
	"lib/globals"

	"github.com/MintzyG/fun"
)

// authChains maps every operationId to the middleware chain it must run
// through, mirroring the parity-test matrix. The setup guard runs on every
// operation except the two /auth/setup routes, which manage the flag
// themselves.
func authChains(mw middlewares) map[string][]func(http.Handler) http.Handler {
	setupGuard := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !globals.SetupComplete() {
				fun.ServiceUnavailable("please setup IDX first on /auth/setup").Send(w)
				return
			}
			next.ServeHTTP(w, r)
		})
	}

	chains := map[string][]func(http.Handler) http.Handler{
		// authn — public except where noted
		"getSetup":         {},
		"postSetup":        {},
		"postRegister":     {},
		"postLogin":        {},
		"postRefresh":      {},
		"getOAuthConnect":  {},
		"getOAuthCallback": {},
		"getJWKS":          {},
		"postLogout":       {mw.jwtAuth},
		"getIntrospect":    {mw.anyAuth},
		// actors
		"listActors":      {mw.anyAuth, mw.clientOnly},
		"createActor":     {mw.anyAuth, mw.clientOnly},
		"getActor":        {mw.anyAuth, mw.clientOnly},
		"getActorByEmail": {mw.anyAuth, mw.clientOnly},
		// api_keys
		"createAPIKey": {mw.anyAuth, mw.clientOnly},
		// capabilities
		"listCapabilities": {mw.anyAuth},
		"createCapability": {mw.jwtAuth, mw.clientOnly},
		// organizations
		"listOrganizations":               {mw.jwtAuth, mw.clientOnly},
		"createOrganization":              {mw.jwtAuth, mw.clientOnly},
		"listOrganizationMembers":         {mw.jwtAuth, mw.clientOnly},
		"addOrganizationMember":           {mw.jwtAuth, mw.clientOnly},
		"removeOrganizationMember":        {mw.jwtAuth, mw.clientOnly},
		"listOrganizationProjects":        {mw.jwtAuth, mw.clientOnly},
		"createOrganizationProject":       {mw.jwtAuth, mw.clientOnly},
		"listOrganizationProjectActors":   {mw.jwtAuth, mw.clientOnly},
		"createOrganizationProjectActor":  {mw.jwtAuth, mw.clientOnly},
		"listOrganizationProjectMembers":  {mw.jwtAuth, mw.clientOnly},
		"addOrganizationProjectMember":    {mw.jwtAuth, mw.clientOnly},
		"removeOrganizationProjectMember": {mw.jwtAuth, mw.clientOnly},
		"getOrganizationProjectActor":     {mw.jwtAuth, mw.clientOnly},
		// projects
		"listProjects":        {mw.anyAuth, mw.clientOnly},
		"createProject":       {mw.anyAuth, mw.clientOnly},
		"listProjectMembers":  {mw.anyAuth, mw.clientOnly},
		"addProjectMember":    {mw.anyAuth, mw.clientOnly},
		"removeProjectMember": {mw.anyAuth, mw.clientOnly},
		// profiles
		"getPlatformProfile":    {mw.jwtAuth, mw.clientOnly},
		"upsertPlatformProfile": {mw.jwtAuth, mw.clientOnly},
		"getProjectProfile":     {mw.jwtAuth, mw.clientOnly},
		"upsertProjectProfile":  {mw.jwtAuth, mw.clientOnly},
		// profile_schemas
		"getPlatformProfileSchema":    {mw.jwtAuth, mw.clientOnly},
		"upsertPlatformProfileSchema": {mw.jwtAuth, mw.clientOnly},
		"getProjectProfileSchema":     {mw.jwtAuth, mw.clientOnly},
		"upsertProjectProfileSchema":  {mw.jwtAuth, mw.clientOnly},
	}

	for key, chain := range chains {
		if key == "getSetup" || key == "postSetup" {
			continue
		}
		chains[key] = append([]func(http.Handler) http.Handler{setupGuard}, chain...)
	}
	return chains
}

// authDispatch is the strict-server middleware that resolves the auth
// chain for each operation (by operationId) and runs it around the
// handler. Public operations (empty chain) pass through untouched.
func authDispatch(mw middlewares) handler.StrictMiddlewareFunc {
	chains := authChains(mw)
	return func(f handler.StrictHandlerFunc, operationID string) handler.StrictHandlerFunc {
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
				resp, ferr = f(ctx, w, r, request)
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
