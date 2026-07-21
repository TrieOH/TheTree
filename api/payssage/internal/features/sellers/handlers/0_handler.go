package handlers

import (
	"net/http"

	"payssage/internal/features/sellers/queries"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	queries *queries.Queries
}

func NewHandlers(queries *queries.Queries) *Handlers {
	return &Handlers{queries: queries}
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
