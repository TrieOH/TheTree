package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (handler *Handlers) Delete(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	signatureID, err := req.Path("signature_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	err = handler.commands.Delete(r.Context(), signatureID)
	if fun.Bail(w, err) {
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
