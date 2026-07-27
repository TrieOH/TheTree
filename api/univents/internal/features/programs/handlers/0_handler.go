package handlers

import (
	"net/http"
	"univents/internal/features/programs/commands"
	"univents/internal/features/programs/queries"

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
	r.Get("/editions/{edition_id}/programs", h.ListProgramsByEdition)
	r.With(jwt).Post("/editions/{edition_id}/programs", h.CreateProgram)
	r.Get("/programs/{program_id}", h.GetProgramByID)
	r.With(jwt).Patch("/programs/{program_id}", h.PatchProgram)
	r.With(jwt).Delete("/programs/{program_id}", h.DeleteProgram)

	r.Get("/programs/{program_id}/occurrences", h.ListOccurrencesByProgram)
	r.With(jwt).Post("/programs/{program_id}/occurrences", h.CreateOccurrence)
	r.Get("/editions/{edition_id}/occurrences", h.ListOccurrencesByEdition)
	r.Get("/occurrences/{occurrence_id}", h.GetOccurrenceByID)
	r.With(jwt).Patch("/occurrences/{occurrence_id}", h.PatchOccurrence)
	r.With(jwt).Delete("/occurrences/{occurrence_id}", h.DeleteOccurrence)
}
