package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (h *Handlers) GetByID(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	editionID, err := req.Path("edition_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	sigID, err := req.Path("sig_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	sig, err := h.queries.GetByID(r.Context(), sigID, editionID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, sig)
}
