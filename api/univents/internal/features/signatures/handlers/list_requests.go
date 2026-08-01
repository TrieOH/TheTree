package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (h *Handlers) ListRequestsByEdition(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	editionID, err := req.Path("edition_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	requests, err := h.ops.ListRequestsByEdition(r.Context(), editionID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, requests)
}
