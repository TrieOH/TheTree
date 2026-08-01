package handlers

import (
	"net/http"
	"univents/internal/features/ticket_types"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	ops *ticket_types.Operations
}

func NewHandlers(ops *ticket_types.Operations) *Handlers {
	return &Handlers{ops: ops}
}

func RegisterRoutes(
	r *chi.Mux,
	h *Handlers,
	jwt func(http.Handler) http.Handler,
) {
	r.Get("/editions/{edition_id}/ticket-types", h.List)
	r.Get("/ticket-types/{ticket_type_id}", h.GetByID)
	r.With(jwt).Patch("/ticket-types/{ticket_type_id}", h.Patch)
	r.With(jwt).Post("/editions/{edition_id}/ticket-types", h.Create)
}
