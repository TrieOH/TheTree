package handlers

import (
	"net/http"

	"payssage/internal/features/intents/commands"
	"payssage/internal/features/intents/queries"

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
		r.Get("/intents", h.ListByProfile)
		r.Get("/intents/{intent_id}", h.GetByID)
		r.Get("/wallets/{wallet_id}/intents", h.ListByWallet)
		r.Get("/organizations/{organization_id}/intents", h.ListByOrg)
		r.Post("/wallets/{wallet_id}/intents", h.Checkout)
	})
}
