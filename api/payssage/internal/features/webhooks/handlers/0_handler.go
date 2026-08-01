package handlers

import (
	"payssage/internal/features/webhooks"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	ops *webhooks.Operations
}

func NewHandlers(ops *webhooks.Operations) *Handlers {
	return &Handlers{ops: ops}
}

func RegisterRoutes(
	r *chi.Mux,
	h *Handlers,
) {
	r.Group(func(r chi.Router) {
		r.Post("/webhooks/{provider}", h.Receive)
	})
}
