package handlers

import (
	"net/http"
	"univents/internal/features/activities/commands"
	"univents/internal/features/activities/queries"

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
	r.Route("/events/{event_id}/editions/{edition_id}/activities", func(r chi.Router) {
		r.Get("/", h.List)
		r.With(jwt).Post("/", h.Create)
		r.With(jwt).Get("/admin", h.ListAdmin)
		r.With(jwt).Route("/{activity_id}", func(r chi.Router) {
			r.Post("/publish", h.Publish)
			r.Post("/complete", h.Complete)
			r.Post("/register", h.Register)
			r.Post("/unregister", h.Unregister)
			r.Get("/records", h.ListRecords)
			r.Post("/records/{record_id}", h.MarkAttendance)
		})
	})
}
