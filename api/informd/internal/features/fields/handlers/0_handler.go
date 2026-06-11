package handlers

import (
	"Informd/internal/features/fields/commands"
	"Informd/internal/features/fields/queries"
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
	r.Group(func(r chi.Router) {
		r.Use(anyAuth)
		r.Post("/forms/{form_id}/steps/{step_id}/fields", h.CreateField)
		r.Put("/forms/{form_id}/steps/{step_id}/fields", h.BulkEditFields)
		r.Get("/forms/{form_id}/steps/{step_id}/fields", h.List)
		r.Get("/forms/{form_id}/steps/{step_id}/fields/{field_id}/select", h.GetSelectConfig)
		r.Delete("/forms/{form_id}/steps/{step_id}/fields/{field_id}", h.DeleteField)
		r.Put("/forms/{form_id}/steps/{step_id}/fields/{field_id}/select", h.EditSelectConfig)

		r.Post("/namespaces/{namespace_id}/forms/{form_id}/steps/{step_id}/fields", h.CreateNamespacedField)
		r.Put("/namespaces/{namespace_id}/forms/{form_id}/steps/{step_id}/fields", h.BulkEditNamespacedFields)
		r.Get("/namespaces/{namespace_id}/forms/{form_id}/steps/{step_id}/fields", h.ListNamespaced)
		r.Get("/namespaces/{namespace_id}/forms/{form_id}/steps/{step_id}/fields/{field_id}/select", h.GetSelectConfigNamespaced)
		r.Delete("/namespaces/{namespace_id}/forms/{form_id}/steps/{step_id}/fields/{field_id}", h.DeleteNamespacedField)
		r.Put("/namespaces/{namespace_id}/forms/{form_id}/steps/{step_id}/fields/{field_id}/select", h.EditSelectConfigNamespaced)
	})
}
