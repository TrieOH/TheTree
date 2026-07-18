package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (h *Handlers) GetByID(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	endpointID, err := req.Path("endpoint_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	endpoint, err := h.queries.GetByID(r.Context(), endpointID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, endpoint, http.StatusOK)
}
