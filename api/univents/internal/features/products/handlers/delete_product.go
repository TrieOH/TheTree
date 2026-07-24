package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (handler *Handlers) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	productID, err := req.Path("product_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	err = handler.commands.DeleteProduct(r.Context(), productID)
	if fun.Bail(w, err) {
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
