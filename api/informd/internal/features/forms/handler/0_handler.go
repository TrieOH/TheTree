package handler

import (
	"Informd/internal/features/forms/commands"
	"Informd/internal/features/forms/queries"
	"Informd/models"
	"net/http"

	"github.com/MintzyG/fun/middlewares"
	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	commands *commands.CommandService
	queries  *queries.QueryService
}

func NewHandlers(
	commands *commands.CommandService,
	queries *queries.QueryService,
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
	r.Group(func(r chi.Router) {
		r.Use(anyAuth)
		r.Post("/forms", h.Create)
		r.With(middlewares.WithParams[models.BulkGetParams]()).Post("/forms/bulk", h.BulkGet)
		r.Post("/namespaces/{namespace_id}/forms", h.CreateInNamespace)
		r.Post("/forms/{form_id}/steps", h.Create)
	})
}
