package handlers

import (
	"net/http"
	"payssage/internal/features/wallets"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	ops *wallets.Operations
}

func NewHandlers(ops *wallets.Operations) *Handlers {
	return &Handlers{ops: ops}
}

func RegisterRoutes(
	r *chi.Mux,
	h *Handlers,
	jwtAuth func(http.Handler) http.Handler,
) {
	r.Group(func(r chi.Router) {
		r.Use(jwtAuth)
		r.Post("/wallets", h.Create)
		r.Get("/wallets", h.List)
		r.Get("/wallets/{wallet_id}", h.GetByID)
		r.Patch("/wallets/{wallet_id}/fee", h.SetFeeBPS)
		r.Patch("/wallets/{wallet_id}/sandbox", h.SetSandboxState)
		r.Get("/organizations/{organization_id}/wallets", h.ListFromOrg)
		r.Post("/wallets/{wallet_id}/collector", h.BindCollector)
		r.Delete("/wallets/{wallet_id}/collector", h.UnbindCollector)
	})
}
