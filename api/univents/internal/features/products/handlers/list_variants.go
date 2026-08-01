package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (h *Handlers) ListVariants(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	productID, err := req.Path("product_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	variants, err := h.ops.ListVariantsByProduct(r.Context(), productID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, variants)
}
