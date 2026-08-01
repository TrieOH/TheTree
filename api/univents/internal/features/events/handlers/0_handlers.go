package handlers

import (
	"net/http"
	"univents/internal/features/events"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	ops *events.Operations
}

func NewHandlers(ops *events.Operations) *Handlers {
	return &Handlers{ops: ops}
}

func RegisterRoutes(
	r *chi.Mux,
	h *Handlers,
	jwt func(http.Handler) http.Handler,
) {
	r.Route("/events", func(r chi.Router) {
		r.Get("/", h.ListPublic)
		r.Get("/{event_slug}:by-slug", h.GetBySlug)
		r.With(jwt).Post("/", h.Create)
		r.With(jwt).Get("/owned", h.ListOwned)
		r.With(jwt).Get("/joined", h.ListJoined)
		r.With(jwt).Route("/{event_id}", func(r chi.Router) {
			r.Patch("/", h.Patch)
			r.Post("/publish", h.Publish)
			r.Post("/discontinue", h.Discontinue)

			// Members
			r.Get("/members", h.ListMembers)
			r.Post("/members", h.AddMember)
			r.Delete("/members/{user_id}", h.RemoveMember)
		})
	})
}
