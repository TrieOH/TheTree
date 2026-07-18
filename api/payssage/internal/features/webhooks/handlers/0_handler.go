package handlers

import (
	"payssage/internal/features/webhooks/commands"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	commands *commands.Commands
}

func NewHandlers(
	commands *commands.Commands,
) *Handlers {
	return &Handlers{
		commands: commands,
	}
}

func RegisterRoutes(
	r *chi.Mux,
	h *Handlers,
) {
	r.Group(func(r chi.Router) {
		r.Post("/webhooks/{provider}", h.Receive)
	})
}
