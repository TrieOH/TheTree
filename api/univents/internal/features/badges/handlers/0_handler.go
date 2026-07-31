package handlers

import (
	"net/http"
	"univents/internal/features/badges/commands"
	"univents/internal/features/badges/queries"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	cmds  *commands.Commands
	quers *queries.Queries
}

func NewHandler(cmds *commands.Commands, quers *queries.Queries) *Handler {
	return &Handler{
		cmds:  cmds,
		quers: quers,
	}
}

func RegisterRoutes(r *chi.Mux, h *Handler, jwt func(http.Handler) http.Handler) {
	r.Route("/events/{event_id}/editions/{edition_id}/badges", func(r chi.Router) {
		r.Use(jwt)
		r.Post("/", h.CreateTemplate)
		r.Get("/", h.ListTemplates)
		r.Get("/{template_id}", h.GetTemplate)
		r.Delete("/{template_id}", h.DeleteTemplate)
	})
}
