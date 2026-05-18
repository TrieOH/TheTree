package handlers

import (
	"Informd/internal/features/namespaces/commands"
	"Informd/internal/features/namespaces/queries"
	"Informd/models"
	"net/http"

	"github.com/MintzyG/fun/middlewares"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	commands *commands.CommandService
	queries  *queries.QueryService
}

func NewHandler(
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
		r.Get("/namespaces", h.ListNamespaces)
		r.Post("/namespaces", h.Create)
		r.Get("/namespaces/{namespace_id}/members", h.ListMembers)
		r.Post("/namespaces/{namespace_id}/members", h.AddMember)
		r.Delete("/namespaces/{namespace_id}/members", h.RemoveMember)
		r.With(middlewares.WithParams[models.BulkGetParams]()).Post("/namespaces/bulk", h.BulkGet)
	})
}
