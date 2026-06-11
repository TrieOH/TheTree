package handlers

import (
	"Informd/internal/features/steps/commands"
	"Informd/internal/features/steps/queries"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	commands *commands.Command
	queries  *queries.Queries
}

func NewHandlers(
	commands *commands.Command,
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
	anyAuth func(http.Handler) http.Handler,
) {
	r.Get("/forms/{form_id}/steps", h.List)
	r.Get("/namespaces/{namespace_id}/forms/{form_id}/steps", h.ListNamespaced)
	r.Group(func(r chi.Router) {
		r.Use(anyAuth)
		r.Post("/forms/{form_id}/steps", h.CreateStep)
		r.Put("/forms/{form_id}/steps", h.BulkEditSteps)
		r.Post("/namespaces/{namespace_id}/forms/{form_id}/steps", h.CreateNamespacedStep)
		r.Put("/namespaces/{namespace_id}/forms/{form_id}/steps", h.BulkEditNamespacedSteps)
	})
}
