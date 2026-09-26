package app

import (
	"net/http"
	"testing"

	libauthz "lib/authz"
	spec "payssage"
	"payssage/internal/handlers"

	"github.com/go-chi/chi/v5"
)

// pass-through stubs for the shared test router.
func mwAPIKey(next http.Handler) http.Handler { return next }
func mwAny(next http.Handler) http.Handler {
	return next
}

// newTestRouter mounts the strict server with the real middleware stack
// (raw-request capture + validation + auth dispatch + fun-envelope error
// handlers) on a fresh chi router plus harness routes.
func newTestRouter(t *testing.T, h *handlers.Server, primitives libauthz.Primitives) *chi.Mux {
	t.Helper()
	resolver, err := libauthz.NewResolver(spec.OpenAPISpec, primitives, libauthz.Options{})
	if err != nil {
		t.Fatalf("resolve auth chains: %v", err)
	}
	r := chi.NewRouter()
	mountStrict(r, h, resolver.Chains())
	return r
}
