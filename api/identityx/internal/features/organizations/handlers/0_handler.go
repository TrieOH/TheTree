package handlers

import (
	"IdentityX/internal/features/organizations/commands"
	"IdentityX/internal/features/organizations/queries"
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
		r.Get("/organizations", h.ListOrgs)
		r.Post("/organizations", h.Create)
		r.Get("/organizations/{organization_id}/members", h.ListMembers)
		r.Post("/organizations/{organization_id}/members", h.AddMember)
		r.Delete("/organizations/{organization_id}/members", h.RemoveMember)
	})
}
