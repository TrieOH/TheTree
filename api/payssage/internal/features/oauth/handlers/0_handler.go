package handlers

import (
	"net/http"

	"payssage/internal/features/oauth/commands"

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
) {
	r.Group(func(r chi.Router) {
		r.With(jwtAuth).Post("/providers/{provider}/connect", h.Connect)
		r.Get("/providers/{provider}/callback", h.Callback)
	})
}
