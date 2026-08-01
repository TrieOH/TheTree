package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (h *Handlers) Delete(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	signatureID, err := req.Path("signature_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	err = h.ops.Delete(r.Context(), signatureID)
	if fun.Bail(w, err) {
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
