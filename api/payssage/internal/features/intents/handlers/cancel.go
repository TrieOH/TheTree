package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (h *Handlers) Cancel(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	intentID, err := req.Path("intent_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	intent, err := h.ops.Cancel(r.Context(), intentID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, intent)
}
