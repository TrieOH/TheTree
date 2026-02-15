package router

import (
	eventhttp "univents/internal/eventcore/interfaces/http"
	"univents/internal/interfaces/http/middleware"
	systemhttp "univents/internal/interfaces/http/system"

	"github.com/go-chi/chi/v5"
)

func registerRoutes(r *chi.Mux, deps *HTTPDeps) {
	registerSystemRoutes(r, deps.SystemHandler, deps.AuthMiddleware)
	registerEventsRoutes(r, deps.EventsHandler, deps.AuthMiddleware)
}

func registerSystemRoutes(
	r *chi.Mux,
	h *systemhttp.UniventsHandler,
	authMW *middleware.AuthMiddleware,
) {
	r.Group(func(r chi.Router) {
		r.Get("/health", h.Health)
		r.With(authMW.Auth()).
			Get("/protected/health", h.ProtectedHealth)
	})
}

func registerEventsRoutes(
	r *chi.Mux,
	h *eventhttp.EventsHandler,
	authMW *middleware.AuthMiddleware,
) {
	r.Group(func(r chi.Router) {
		r.Route("/events", func(r chi.Router) {
			r.Get("/", h.List)
			r.With(authMW.Auth()).
				Post("/", h.Create)
			r.With(authMW.Auth()).
				Post("/{event_id}/publish", h.Publish)
		})
	})
}
