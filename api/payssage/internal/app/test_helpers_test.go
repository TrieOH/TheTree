package app

import (
	"payssage/internal/handlers"
	"payssage/internal/handlers/webhooks"
	"payssage/internal/openapi"

	"github.com/go-chi/chi/v5"
)

// newTestRouter mounts the strict server with the real middleware stack
// (raw-request capture + validation + auth dispatch + fun-envelope error
// handlers) on a fresh chi router plus harness routes.
func newTestRouter(h *handlers.Server, mw middlewares) *chi.Mux {
	r := chi.NewRouter()
	mountStrict(r, h, mw)
	return r
}

var _ = webhooks.RawRequestMiddleware
var _ openapi.MiddlewareFunc = webhooks.RawRequestMiddleware
