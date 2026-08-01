package handlers

import (
	"net/http"
	"univents/internal/features/editions"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	ops *editions.Operations
}

func NewHandlers(ops *editions.Operations) *Handlers {
	return &Handlers{ops: ops}
}

func RegisterRoutes(
	r *chi.Mux,
	h *Handlers,
	jwt func(http.Handler) http.Handler,
) {
	r.Route("/events/{event_slug}:by-slug/editions/{edition_slug}:by-slug", func(r chi.Router) {
		r.Get("/", h.GetBySlug)
	})

	r.Route("/events/{event_id}/editions", func(r chi.Router) {
		r.Get("/", h.ListPublic)
		r.Get("/active", h.GetActive)
		r.Get("/past", h.GetPast)
		r.Get("/upcoming", h.GetUpcoming)
		r.With(jwt).Post("/", h.Create)
		r.With(jwt).Get("/draft", h.ListDraft)
		r.With(jwt).Patch("/{edition_id}", h.Patch)
		r.With(jwt).Post("/{edition_id}/publish", h.Publish)
	})
}
