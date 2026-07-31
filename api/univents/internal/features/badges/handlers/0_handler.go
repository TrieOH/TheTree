package handlers

import (
	"net/http"
	"univents/internal/features/badges/commands"
	"univents/internal/features/badges/queries"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	commands *commands.Commands
	queries  *queries.Queries
}

func NewHandler(commands *commands.Commands, queries *queries.Queries) *Handler {
	return &Handler{
		commands: commands,
		queries:  queries,
	}
}

func RegisterRoutes(r *chi.Mux, h *Handler, jwt func(http.Handler) http.Handler) {
	r.Route("/editions/{edition_id}/badges", func(r chi.Router) {
		r.Use(jwt)
		r.Post("/", h.CreateTemplate)
		r.Get("/", h.ListTemplates)
	})
	r.Route("/badges/{template_id}", func(r chi.Router) {
		r.Use(jwt)
		r.Get("/", h.GetTemplate)
		r.Delete("/", h.DeleteTemplate)
	})
}
