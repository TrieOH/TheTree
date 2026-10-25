package app

import (
	"net/http"
	"testing"

	spec "IdentityX"
	"IdentityX/internal/authz"
	"IdentityX/internal/setup"
	libauthz "lib/authz"

	"github.com/go-chi/chi/v5"
)

// pass-through stubs for the shared test router.
func mwJWT(next http.Handler) http.Handler     { return next }
func mwAnyAuth(next http.Handler) http.Handler { return next }

// testScopeCheckers returns the real scope-checker registry, so routers
// built in tests exercise the same scope enforcement production gets. The
// platform-only checker only reads the context identity, so a service with
// nil repos is a safe host for it here.
func testScopeCheckers() map[string]libauthz.ScopeChecker {
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
func newTestRouter(t *testing.T, h *server, primitives libauthz.Primitives, st *setup.State) *chi.Mux {
	t.Helper()
	resolver, err := libauthz.NewResolver(spec.OpenAPISpec, primitives, chainOptions(st))
	if err != nil {
		t.Fatalf("resolve auth chains: %v", err)
	}
	r := chi.NewRouter()
	mountStrict(r, h, resolver.Chains())
	return r
}

// testPrimitives are pass-through auth stubs with the real scope checkers.
func testPrimitives() libauthz.Primitives {
	return libauthz.Primitives{
		JWT:    mwJWT,
		APIKey: mwJWT,
		Any:    mwAnyAuth,
		Scopes: testScopeCheckers(),
	}
}
