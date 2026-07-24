package handlers

import (
	"net/http"
	"univents/internal/features/ticket_types/commands"
	"univents/internal/features/ticket_types/queries"

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
	jwt func(http.Handler) http.Handler,
) {
	r.Get("/editions/{edition_id}/ticket-types", h.List)
	r.Get("/ticket-types/{ticket_type_id}", h.GetByID)
	r.With(jwt).Patch("/ticket-types/{ticket_type_id}", h.Patch)
	r.With(jwt).Post("/editions/{edition_id}/ticket-types", h.Create)
}
