package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (h *Handlers) ListByEdition(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	editionID, err := req.Path("edition_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	signatures, err := h.ops.ListByEdition(r.Context(), editionID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, signatures)
}
