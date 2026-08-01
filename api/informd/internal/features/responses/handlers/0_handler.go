package handlers

import (
	"Informd/internal/features/responses"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	ops *responses.Operations
}

func NewHandlers(ops *responses.Operations) *Handlers {
	return &Handlers{
		ops: ops,
	}
}

func RegisterRoutes(
	r *chi.Mux,
	h *Handlers,
) {
	r.Group(func(r chi.Router) {
		r.Post("/forms/{form_id}/responses", h.Submit)
	})
}
