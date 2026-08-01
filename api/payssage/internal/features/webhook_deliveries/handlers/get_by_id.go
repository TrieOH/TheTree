package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (h *Handlers) GetByID(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	deliveryID, err := req.Path("delivery_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	delivery, err := h.ops.GetByID(r.Context(), deliveryID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, delivery, http.StatusOK)
}
