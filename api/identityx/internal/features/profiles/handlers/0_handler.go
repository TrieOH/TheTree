package handlers

import (
	"IdentityX/internal/features/profiles/commands"
	"IdentityX/internal/features/profiles/queries"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	queries  *queries.Queries
	commands *commands.Commands
}

func New(
	queries *queries.Queries,
	commands *commands.Commands,
) *Handlers {
	return &Handlers{
		queries:  queries,
		commands: commands,
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
