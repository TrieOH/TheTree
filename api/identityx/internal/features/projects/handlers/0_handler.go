package handlers

import (
	"IdentityX/internal/features/projects/commands"
	"IdentityX/internal/features/projects/queries"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	commands *commands.Commands
	queries  *queries.Queries
}

func NewHandlers(
	commands *commands.Commands,
	queries *queries.Queries,
) *Handlers {
	return &Handlers{
		commands: commands,
		queries:  queries,
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
		r.Get("/projects", h.List)
		r.Post("/projects", h.Create)
		r.Get("/projects/{project_id}/members", h.ListMembers)
		r.Post("/projects/{project_id}/members", h.AddMember)
		r.Delete("/projects/{project_id}/members", h.RemoveMember)
	})
}
