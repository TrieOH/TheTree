package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (h *Handlers) GetByVendorCode(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	editionID, err := req.Path("edition_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	vendorCode, err := req.Path("vendor_code").StringRequired()
	if fun.Bail(w, err) {
		return
	}
	product, err := h.ops.GetProductByVendorCode(r.Context(), editionID, vendorCode)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, product)
}
