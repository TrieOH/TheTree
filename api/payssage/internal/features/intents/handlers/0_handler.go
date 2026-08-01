package handlers

import (
	"net/http"
	"payssage/internal/features/intents"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	ops *intents.Operations
}

func NewHandlers(ops *intents.Operations) *Handlers {
	return &Handlers{ops: ops}
}

func RegisterRoutes(
	r *chi.Mux,
	h *Handlers,
	jwtAuth func(http.Handler) http.Handler,
) {
	r.Group(func(r chi.Router) {
		r.Use(jwtAuth)
		r.Get("/intents", h.ListByProfile)
		r.Get("/intents/{intent_id}", h.GetByID)
		r.Post("/intents/{intent_id}/cancel", h.Cancel)
		r.Get("/wallets/{wallet_id}/intents", h.ListByWallet)
		r.Get("/organizations/{organization_id}/intents", h.ListByOrg)
		r.Post("/wallets/{wallet_id}/intents", h.Checkout)
		r.Post("/testmode/intents/create", h.HardCreate)
	})
}
