package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (h *Handlers) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	productID, err := req.Path("product_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	err = h.ops.DeleteProduct(r.Context(), productID)
	if fun.Bail(w, err) {
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
