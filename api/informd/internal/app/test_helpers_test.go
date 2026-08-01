package app

import (
	"Informd/internal/handlers"

	"github.com/go-chi/chi/v5"
)

// newTestRouter mounts the strict server with the real middleware stack
// (validation + auth dispatch + fun-envelope error handlers) on a fresh
// chi router plus harness routes.
func newTestRouter(h *handlers.Server, mw middlewares) *chi.Mux {
	r := chi.NewRouter()
	mountStrict(r, h, mw)
	return r
}
