package handlers

import (
	"net/http"
	"payssage/internal/features/webhook_events"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	ops *webhook_events.Operations
}

func NewHandlers(ops *webhook_events.Operations) *Handlers {
	return &Handlers{ops: ops}
}

func RegisterRoutes(
	r *chi.Mux,
	h *Handlers,
	jwtAuth func(http.Handler) http.Handler,
) {
	r.Group(func(r chi.Router) {
		r.Use(jwtAuth)
		r.Get("/wallets/{wallet_id}/webhooks/events", h.ListByWallet)
		r.Get("/webhooks/events/{event_id}", h.GetByID)
	})
}
