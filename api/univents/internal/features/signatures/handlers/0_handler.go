package handlers

import (
	"net/http"
	"univents/internal/features/signatures/commands"
	"univents/internal/features/signatures/queries"

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
) {
	r.Group(func(r chi.Router) {
		r.Use(jwtAuth)
		r.Post("/events/{event_id}/editions/{edition_id}/signatures", h.Add)
		r.Get("/events/{event_id}/editions/{edition_id}/signatures", h.List)
		r.Delete("/events/{event_id}/editions/{edition_id}/signatures/{sig_id}", h.Remove)
		r.Get("/events/{event_id}/editions/{edition_id}/signatures/{sig_id}", h.GetByID)
	})
}
