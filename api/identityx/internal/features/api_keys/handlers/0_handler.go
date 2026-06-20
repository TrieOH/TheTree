package handlers

import (
	"IdentityX/internal/features/api_keys/commands"
	"net/http"

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
	jwtAuth func(http.Handler) http.Handler,
	clientOnly func(http.Handler) http.Handler,
) {
	r.Group(func(r chi.Router) {
		r.Use(jwtAuth, clientOnly)
		r.Post("/projects/{project_id}/api_keys", h.Create)
	})
}
