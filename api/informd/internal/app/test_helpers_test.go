package app

import (
	"testing"

	spec "Informd"
	"Informd/internal/handlers"
	libauthz "lib/authz"

	"github.com/go-chi/chi/v5"
)

// newTestRouter mounts the strict server with the real middleware stack
// (validation + auth dispatch + fun-envelope error handlers) on a fresh
// chi router plus harness routes.
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
