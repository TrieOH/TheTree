package handlers

import (
	"net/http"
	"payssage/internal/features/webhook_endpoints"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	ops *webhook_endpoints.Operations
}

func NewHandlers(ops *webhook_endpoints.Operations) *Handlers {
	return &Handlers{ops: ops}
}

func RegisterRoutes(
	r *chi.Mux,
	h *Handlers,
	jwtAuth func(http.Handler) http.Handler,
) {
	r.Group(func(r chi.Router) {
		r.Use(jwtAuth)
		r.Post("/wallets/{wallet_id}/webhooks/endpoints", h.Create)
		r.Get("/wallets/{wallet_id}/webhooks/endpoints", h.ListByWallet)
		r.Get("/webhooks/endpoints/{endpoint_id}", h.GetByID)
		r.Delete("/webhooks/endpoints/{endpoint_id}", h.Delete)
	})
}
