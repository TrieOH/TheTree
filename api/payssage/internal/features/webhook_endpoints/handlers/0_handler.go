package handlers

import (
	"net/http"

	"payssage/internal/features/webhook_endpoints/commands"
	"payssage/internal/features/webhook_endpoints/queries"

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
		r.Post("/wallets/{wallet_id}/webhooks/endpoints", h.Create)
		r.Get("/wallets/{wallet_id}/webhooks/endpoints", h.ListByWallet)
		r.Get("/webhooks/endpoints/{endpoint_id}", h.GetByID)
		r.Delete("/webhooks/endpoints/{endpoint_id}", h.Delete)
	})
}
