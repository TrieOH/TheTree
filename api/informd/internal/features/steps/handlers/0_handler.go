package handlers

import (
	"Informd/internal/features/steps/commands"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	commands *commands.Command
}

func NewHandlers(
	commands *commands.Command,
) *Handlers {
	return &Handlers{
		commands: commands,
	}
}

func RegisterRoutes(
	r *chi.Mux,
	h *Handlers,
	anyAuth func(http.Handler) http.Handler,
) {
	r.Group(func(r chi.Router) {
		r.Use(anyAuth)
		r.Post("/forms/{form_id}/steps", h.CreateStep)
		r.Put("/forms/{form_id}/steps", h.BulkEditSteps)
		r.Post("/namespaces/{namespace_id}/forms/{form_id}/steps", h.CreateNamespacedStep)
		r.Put("/namespaces/{namespace_id}/forms/{form_id}/steps", h.BulkEditNamespacedSteps)
	})
}
