package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (h *Handlers) ListCertsByEdition(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	editionID, err := req.Path("edition_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	certs, err := h.ops.ListCertsByEdition(r.Context(), editionID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, certs)
}
