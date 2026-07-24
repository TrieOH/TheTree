package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (handler *Handlers) DeleteVariant(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	variantID, err := req.Path("variant_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	err = handler.commands.DeleteVariant(r.Context(), variantID)
	if fun.Bail(w, err) {
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
