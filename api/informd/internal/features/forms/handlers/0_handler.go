package handlers

import (
	"Informd/internal/features/forms/commands"
	"Informd/internal/features/forms/queries"
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
	anyAuth func(http.Handler) http.Handler,
) {
	r.Group(func(r chi.Router) {
		r.Use(anyAuth)
		r.Post("/forms", h.Create)
		r.Get("/forms", h.ListMine)
		r.Get("/forms/archived", h.ListMineArchived)
		r.Get("/forms/{form_id}/full", h.GetFull)
		r.Get("/forms/{form_id}/members", h.ListMembers)
		r.Post("/forms/{form_id}/members", h.AddMember)
		r.Delete("/forms/{form_id}/members", h.RemoveMember)
		r.Post("/forms/{form_id}/open", h.Open)
		r.Post("/forms/{form_id}/close", h.Close)
		r.Post("/forms/{form_id}/archive", h.Archive)
		r.Post("/forms/{form_id}/redraft", h.ReDraft)
		r.Get("/forms/{form_id}/responses/count", h.ResponseCount)
	})
}
