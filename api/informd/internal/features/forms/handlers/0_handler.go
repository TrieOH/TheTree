package handlers

import (
	"Informd/internal/features/forms/commands"
	"Informd/internal/features/forms/queries"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	commands *commands.CommandService
	queries  *queries.QueryService
}

func NewHandlers(
	commands *commands.CommandService,
	queries *queries.QueryService,
) *Handlers {
	return &Handlers{
		commands: commands,
		queries:  queries,
	}
}

func RegisterRoutes(
	r *chi.Mux,
	h *Handlers,
	anyAuth func(http.Handler) http.Handler,
) {
	r.Group(func(r chi.Router) {
		r.Use(anyAuth)
		r.Post("/forms", h.Create)
		r.Get("/forms", h.ListMine)
		r.Get("/forms/{form_id}/members", h.ListMembers)
		r.Post("/forms/{form_id}/members", h.AddMember)
		r.Delete("/forms/{form_id}/members", h.RemoveMember)
	})
}
