package handlers

import (
	"Informd/internal/features/steps"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	ops *steps.Operations
}

func NewHandlers(ops *steps.Operations) *Handlers {
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
		r.Post("/forms/{form_id}/steps", h.CreateStep)
		r.Put("/forms/{form_id}/steps", h.BulkEditSteps)
		r.Get("/forms/{form_id}/steps", h.List)
		// TODO: kill these duplicated namespaced routes — CheckForm already anchors via the form's namespace.
		r.Get("/namespaces/{namespace_id}/forms/{form_id}/steps", h.ListNamespaced)
		r.Post("/namespaces/{namespace_id}/forms/{form_id}/steps", h.CreateNamespacedStep)
		r.Put("/namespaces/{namespace_id}/forms/{form_id}/steps", h.BulkEditNamespacedSteps)
	})
}
