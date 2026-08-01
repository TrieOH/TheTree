package handlers

import (
	"net/http"
	"univents/internal/features/badges"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	ops *badges.Operations
}

func NewHandler(ops *badges.Operations) *Handlers {
	return &Handlers{ops: ops}
}

func RegisterRoutes(r *chi.Mux, h *Handlers, jwt func(http.Handler) http.Handler) {
	r.Route("/editions/{edition_id}/badges", func(r chi.Router) {
		r.Use(jwt)
		r.Post("/", h.CreateTemplate)
		r.Get("/", h.ListTemplates)
	})
	r.Route("/badges/{template_id}", func(r chi.Router) {
		r.Use(jwt)
		r.Get("/", h.GetTemplate)
		r.Delete("/", h.DeleteTemplate)
	})
}
