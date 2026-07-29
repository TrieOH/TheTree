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
		r.Get("/", h.ListPublic)
		r.Get("/{event_slug}:by-slug", h.GetBySlug)
		r.With(jwt).Post("/", h.Create)
		r.With(jwt).Get("/owned", h.ListOwned)
		r.With(jwt).Get("/joined", h.ListJoined)
		r.With(jwt).Route("/{event_id}", func(r chi.Router) {
			r.Patch("/", h.Patch)
			r.Post("/publish", h.Publish)
			r.Post("/discontinue", h.Discontinue)

			// Members
			r.Get("/members", h.ListMembers)
			r.Post("/members", h.AddMember)
			r.Delete("/members/{user_id}", h.RemoveMember)
		})
	})
}
