package handlers

import (
	"IdentityX/internal/features/actors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	ops *actors.Operations
}

func NewHandlers(ops *actors.Operations) *Handlers {
	return &Handlers{
		ops: ops,
	}
}

func RegisterRoutes(
	r *chi.Mux,
	h *Handlers,
	anyAuth func(http.Handler) http.Handler,
	clientOnly func(http.Handler) http.Handler,
) {
	r.Group(func(r chi.Router) {
		r.Use(anyAuth, clientOnly)
		r.Get("/projects/{project_id}/actors/{actor_id}", h.GetByID)
		r.Get("/projects/{project_id}/actors/{actor_email}:by_email", h.GetByEmail)
		r.Post("/projects/{project_id}/actors", h.Create)
		r.Get("/projects/{project_id}/actors", h.List)
	})
}
