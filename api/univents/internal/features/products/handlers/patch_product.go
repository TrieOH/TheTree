package handlers

import (
	"net/http"
	"univents/models"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

func (h *Handlers) PatchProduct(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	productID, err := req.Path("product_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload models.PatchProductRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	product, err := h.ops.PatchProduct(r.Context(), payload.ToInput(productID))
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, product)
}
