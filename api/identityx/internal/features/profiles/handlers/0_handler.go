package handlers

import (
	"IdentityX/internal/features/profiles"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	ops *profiles.Operations
}

func New(ops *profiles.Operations) *Handlers {
	return &Handlers{
		ops: ops,
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
		// platform-scoped (actor has NULL project_id)
		r.Get("/actors/{actor_id}/profile", h.GetPlatformProfile)
		r.Put("/actors/{actor_id}/profile", h.UpsertPlatformProfile)
		// project-scoped
		r.Get("/projects/{project_id}/actors/{actor_id}/profile", h.GetProfile)
		r.Put("/projects/{project_id}/actors/{actor_id}/profile", h.UpsertProfile)
	})
}
