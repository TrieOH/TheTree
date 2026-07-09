package handlers

import (
	"net/http"
	"univents/internal/features/events/commands"
	"univents/internal/features/events/queries"

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
	r.Route("/events", func(r chi.Router) {
		r.Get("/", h.List)
		r.With(jwt).Post("/", h.Create)
		r.With(jwt).Get("/own", h.ListOwn)
		r.With(jwt).Route("/{event_id}", func(r chi.Router) {
			r.Post("/publish", h.Publish)
			r.Post("/gallery", h.AddGalleryImage)
			r.Delete("/gallery", h.RemoveGalleryImage)
			r.Put("/logo", h.SetLogo)
			r.Delete("/logo", h.UnsetLogo)
			r.Put("/banner", h.SetBanner)
			r.Delete("/banner", h.UnsetBanner)
		})
	})
}
