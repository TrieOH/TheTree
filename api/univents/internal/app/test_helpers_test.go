package app

import (
	"testing"

	"lib/authz"
	spec "univents"
	"univents/internal/handlers"

	"github.com/go-chi/chi/v5"
)

// newTestRouter mounts the strict server with the real middleware stack
// (validation + auth dispatch + fun-envelope error handlers) on a fresh
// chi router plus harness routes.
func newTestRouter(t *testing.T, h *handlers.Server, primitives authz.Primitives) *chi.Mux {
	t.Helper()
	resolver, err := authz.NewResolver(spec.OpenAPISpec, primitives, authz.Options{})
	if err != nil {
		t.Fatalf("resolve auth chains: %v", err)
	}
	r := chi.NewRouter()
	mountStrict(r, h, resolver.Chains())
	return r
}
