package handlers

import (
	"net/http"
	"payssage/internal/features/webhook_deliveries"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	ops *webhook_deliveries.Operations
}

func NewHandlers(ops *webhook_deliveries.Operations) *Handlers {
	return &Handlers{ops: ops}
}

func RegisterRoutes(
	r *chi.Mux,
	h *Handlers,
	jwtAuth func(http.Handler) http.Handler,
) {
	r.Group(func(r chi.Router) {
		r.Use(jwtAuth)
		r.Get("/webhooks/endpoints/{endpoint_id}/deliveries", h.ListByEndpoint)
		r.Get("/webhooks/deliveries/{delivery_id}", h.GetByID)
	})
}
