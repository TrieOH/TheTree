package handlers

import (
	"IdentityX/internal/features/profile_schemas/commands"
	"IdentityX/internal/features/profile_schemas/queries"
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
		// project-scoped schema
		r.Get("/projects/{project_id}/profile-schema", h.GetSchema)
		r.Put("/projects/{project_id}/profile-schema", h.UpsertSchema)
		// platform-wide schema
		r.Get("/profile-schema", h.GetSchema)
		r.Put("/profile-schema", h.UpsertSchema)
	})
}
