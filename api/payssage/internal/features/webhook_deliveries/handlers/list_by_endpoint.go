package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (h *Handlers) ListByEndpoint(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	endpointID, err := req.Path("endpoint_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	deliveries, err := h.queries.ListByEndpoint(r.Context(), endpointID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, deliveries, http.StatusOK)
}
