package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (h *Handlers) GetByID(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	id, err := req.Path("product_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	product, err := h.ops.GetProductByID(r.Context(), id)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, product)
}
