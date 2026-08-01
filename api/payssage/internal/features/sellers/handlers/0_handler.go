package handlers

import (
	"net/http"
	"payssage/internal/features/sellers"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	ops *sellers.Operations
}

func NewHandlers(ops *sellers.Operations) *Handlers {
	return &Handlers{ops: ops}
}

func RegisterRoutes(
	r *chi.Mux,
	h *Handlers,
	jwtAuth func(http.Handler) http.Handler,
) {
	r.Group(func(r chi.Router) {
		r.Use(jwtAuth)
		r.Get("/wallets/{wallet_id}/sellers", h.ListByWallet)
	})
}
