package handlers

import (
	"net/http"
	commands2 "univents/internal/features/editions/commands"
	"univents/internal/features/editions/queries"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	commands *commands2.Commands
	queries  *queries.Queries
}

func NewHandlers(
	commands *commands2.Commands,
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
	r.Route("/events/{event_id}/editions", func(r chi.Router) {
		r.Get("/", h.List)
		r.With(jwt).Post("/", h.Create)
		r.With(jwt).Get("/admin", h.ListAdmin)
		r.With(jwt).Route("/{edition_id}", func(r chi.Router) {
			r.Post("/announce", h.Announce)
			r.Post("/payments/connect", h.ConnectPaymentAccount)
			r.Post("/payments/disconnect", h.DisconnectPaymentAccount)
		})
	})
}
