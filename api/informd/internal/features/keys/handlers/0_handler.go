package handlers

import (
	"Informd/internal/features/keys/commands"
	"Informd/internal/features/keys/queries"
	_ "Informd/models"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	commands *commands.CommandService
	queries  *queries.QueryService
}

func NewHandlers(
	commands *commands.CommandService,
	queries *queries.QueryService,
) *Handler {
	return &Handler{
		commands: commands,
		queries:  queries,
	}
}

func RegisterRoutes(
	r *chi.Mux,
	h *Handler,
	jwt func(http.Handler) http.Handler,
) {
	r.Group(func(r chi.Router) {
		r.Use(jwt)
		r.Post("/api-keys", h.Create)
		r.Delete("/api-keys/{id}", h.Revoke)
	})
}
