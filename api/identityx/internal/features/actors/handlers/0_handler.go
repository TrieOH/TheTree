package handlers

import (
	"IdentityX/internal/features/actors/queries"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	queries *queries.Queries
}

func NewHandlers(
	queries *queries.Queries,
) *Handlers {
	return &Handlers{
		queries: queries,
	}
}

func RegisterRoutes(
	r *chi.Mux,
	h *Handlers,
	jwtAuth func(http.Handler) http.Handler,
	clientOnly func(http.Handler) http.Handler,
) {
	r.Group(func(r chi.Router) {
		r.Use(jwtAuth, clientOnly)
		r.Get("/projects/{project_id}/actors/{actor_id}", h.GetByID)
	})
}
