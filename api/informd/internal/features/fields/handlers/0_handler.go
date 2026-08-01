package handlers

import (
	"Informd/internal/features/fields"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	ops *fields.Operations
}

func NewHandlers(ops *fields.Operations) *Handlers {
	return &Handlers{
		ops: ops,
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

		// TODO: kill these duplicated namespaced routes — CheckForm already anchors via the form's namespace.
		r.Post("/namespaces/{namespace_id}/forms/{form_id}/steps/{step_id}/fields", h.CreateNamespacedField)
		r.Put("/namespaces/{namespace_id}/forms/{form_id}/steps/{step_id}/fields", h.BulkEditNamespacedFields)
		r.Get("/namespaces/{namespace_id}/forms/{form_id}/steps/{step_id}/fields", h.ListNamespaced)
		r.Get("/namespaces/{namespace_id}/forms/{form_id}/steps/{step_id}/fields/{field_id}/select", h.GetSelectConfigNamespaced)
		r.Delete("/namespaces/{namespace_id}/forms/{form_id}/steps/{step_id}/fields/{field_id}", h.DeleteNamespacedField)
		r.Put("/namespaces/{namespace_id}/forms/{form_id}/steps/{step_id}/fields/{field_id}/select", h.EditSelectConfigNamespaced)
	})
}
