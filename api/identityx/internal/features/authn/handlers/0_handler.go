package handlers

import (
	"IdentityX/internal/features/authn/commands"

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
		r.Post("/auth/setup", h.Setup)
		r.Post("/auth/register", h.Register)
		r.Post("/auth/login", h.Login)
		r.Post("/auth/logout", h.Logout)
	})
}
