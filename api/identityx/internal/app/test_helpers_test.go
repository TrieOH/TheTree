package app

import (
	"net/http"
	"testing"

	"IdentityX/internal/authz"
	"IdentityX/internal/handlers"
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

// newTestRouter mounts the strict server with the real middleware stack
// (validation + auth dispatch + fun-envelope error handlers) on a fresh
// chi router plus harness routes.
func newTestRouter(t *testing.T, h *handlers.Server, mw middlewares) *chi.Mux {
	t.Helper()
	chains, err := resolveAuthChains(mw)
	if err != nil {
		t.Fatalf("resolveAuthChains: %v", err)
	}
	r := chi.NewRouter()
	mountStrict(r, h, chains)
	return r
}
