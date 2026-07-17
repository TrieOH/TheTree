package handlers

import (
	"net/http"

	"payssage/internal/features/wallets/commands"
	"payssage/internal/features/wallets/queries"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	commands *commands.Commands
	queries  *queries.Queries
}

func NewHandlers(
	commands *commands.Commands,
	queries *queries.Queries,
) *Handlers {
	return &Handlers{
		commands: commands,
		queries:  queries,
	}
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
