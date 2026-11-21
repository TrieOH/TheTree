package app

import (
	"net/http"
	"testing"

	spec "IdentityX"
	"IdentityX/internal/app/errmap"
	"IdentityX/internal/authz"
	"IdentityX/internal/openapi"
	"IdentityX/internal/problems"
	"IdentityX/internal/setup"
	"lib/httpserver"
	"lib/openapispec"

	"github.com/go-chi/chi/v5"
)

// pass-through stubs for the shared test router.
func mwJWT(next http.Handler) http.Handler     { return next }
func mwAnyAuth(next http.Handler) http.Handler { return next }

// testScopeCheckers returns the real scope-checker registry, so routers
// built in tests exercise the same scope enforcement production gets. The
// platform-only checker only reads the context identity, so a service with
// nil repos is a safe host for it here.
func testScopeCheckers() map[string]httpserver.ScopeChecker {
	return authz.New(nil, nil, nil).ScopeCheckers()
}

// completedSetup is a setup state that has already completed.
func completedSetup() *setup.State {
	st := setup.NewState(nil)
	st.MarkComplete()
	return st
}

// newTestRouter mounts the strict server with the real middleware stack
// (validation + auth dispatch + setup guard + fun-envelope error handlers)
// on a fresh chi router.
func newTestRouter(t *testing.T, h *server, primitives httpserver.Primitives, st *setup.State) *chi.Mux {
	t.Helper()
	ops, err := openapispec.Operations(spec.OpenAPISpec)
	if err != nil {
		t.Fatalf("read the spec: %v", err)
	}
	chains, err := httpserver.Chains(ops, primitives)
	if err != nil {
		t.Fatalf("resolve auth chains: %v", err)
	}
	err = guardSetup(chains, st)
	if err != nil {
		t.Fatalf("guard setup: %v", err)
	}
	// Boot installs the problems before serving; the router under test
	// answers errors through the same mapper.
	httpserver.UseProblems(problems.Mapper(errmap.Problem), "http://identityx.test")
	httpserver.UseOperationProblems(ops, true)
	r := chi.NewRouter()
	openapi.Mount(r, h, chains)
	return r
}

// testPrimitives are pass-through auth stubs with the real scope checkers.
func testPrimitives() httpserver.Primitives {
	return httpserver.Primitives{
		JWT:    mwJWT,
		APIKey: mwJWT,
		Any:    mwAnyAuth,
		Scopes: testScopeCheckers(),
	}
}
