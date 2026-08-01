package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (h *Handlers) GetOccurrenceByID(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	id, err := req.Path("occurrence_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	occurrence, err := h.ops.GetOccurrenceByID(r.Context(), id)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, occurrence)
}
