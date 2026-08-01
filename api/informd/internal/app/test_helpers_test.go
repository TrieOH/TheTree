package app

import (
	"net/http"
	"testing"

	"Informd/internal/handlers"

	"github.com/go-chi/chi/v5"
)

// pass-through stubs for the shared test router.
func mwJWT(next http.Handler) http.Handler { return next }

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
