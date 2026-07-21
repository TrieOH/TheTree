package handlers

import (
	"net/http"

	"payssage/internal/features/webhook_deliveries/queries"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	queries *queries.Queries
}

func NewHandlers(
	queries *queries.Queries,
) *Handlers {
	return &Handlers{
		queries: queries,
	}
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
