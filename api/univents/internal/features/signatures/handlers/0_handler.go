package handlers

import (
	"net/http"
	"univents/internal/features/signatures"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	ops *signatures.Operations
}

func NewHandlers(ops *signatures.Operations) *Handlers {
	return &Handlers{ops: ops}
}

func RegisterRoutes(
	r *chi.Mux,
	h *Handlers,
	jwt func(http.Handler) http.Handler,
) {
	r.Get("/editions/{edition_id}/signatures", h.ListByEdition)
	r.Get("/signatures/{signature_id}", h.GetByID)
	r.With(jwt).Post("/editions/{edition_id}/signatures", h.Create)
	r.With(jwt).Delete("/signatures/{signature_id}", h.Delete)

	r.Get("/editions/{edition_id}/signature-requests", h.ListRequestsByEdition)
	r.Get("/signature-requests/{request_id}", h.GetRequestByID)
	r.With(jwt).Post("/editions/{edition_id}/signature-requests", h.CreateRequest)
	r.Post("/signature-requests/fulfill", h.FulfillRequest)
	r.Post("/signature-requests/deny", h.DenyRequest)
	r.With(jwt).Post("/signature-requests/{request_id}/cancel", h.CancelRequest)

	r.Post("/signatures/revoke", h.RevokeSignature)
}
