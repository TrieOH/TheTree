package router

import (
	"univents/internal/commerce/interfaces/http/tickets"
	eventhttp "univents/internal/core/interfaces/http"
	activityhttp "univents/internal/core/interfaces/http/activities"
	editionhttp "univents/internal/core/interfaces/http/editions"
	"univents/internal/interfaces/http/middleware"
	systemhttp "univents/internal/interfaces/http/system"

	"github.com/go-chi/chi/v5"
)

func registerRoutes(r *chi.Mux, deps *HTTPDeps) {
	registerSystemRoutes(r, deps.SystemHandler, deps.AuthMiddleware)
	registerEventsRoutes(r, deps.EventsHandler, deps.AuthMiddleware)
	registerEditionsRoutes(r, deps.EditionsHandler, deps.AuthMiddleware)
	registerTicketsRoutes(r, deps.TicketsHandler, deps.AuthMiddleware)
	registerActivitiesRoutes(r, deps.ActivitiesHandler, deps.AuthMiddleware)
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
	r.Get("/events", h.ListEvents)
	r.Group(func(r chi.Router) {
		r.Use(authMW.Auth())
		r.Get("/events/own", h.ListOwnEvents)
		r.Post("/events", h.CreateEvent)
		r.Patch("/events/{event_id}", h.PatchEvent)
		r.Post("/events/{event_id}/publish", h.PublishEvent)
	})
}

func registerEditionsRoutes(
	r *chi.Mux,
	h *editionhttp.Handler,
	authMW *middleware.AuthMiddleware,
) {
	r.Get("/events/{event_id}/editions", h.List)
	r.Group(func(r chi.Router) {
		r.Use(authMW.Auth())
		r.Post("/events/{event_id}/editions", h.Create)
		r.Post("/events/{event_id}/editions/{edition_id}", h.Announce)
	})
}

func registerTicketsRoutes(
	r *chi.Mux,
	h *tickets.TicketsHandler,
	authMW *middleware.AuthMiddleware,
) {
	r.Group(func(r chi.Router) {
		r.Use(authMW.Auth())
		r.Post("/events/{event_id}/editions/{edition_id}/tickets", h.Create)
	})
}

func registerActivitiesRoutes(
	r *chi.Mux,
	h *activityhttp.Handler,
	authMW *middleware.AuthMiddleware,
) {
	r.Group(func(r chi.Router) {
		r.Use(authMW.Auth())
		r.Post("/events/{event_id}/editions/{edition_id}/activities", h.Create)
		r.Post("/events/{event_id}/editions/{edition_id}/activities/{activity_id}", h.Publish)
	})
}
