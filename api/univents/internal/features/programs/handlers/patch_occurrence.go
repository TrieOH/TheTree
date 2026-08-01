package handlers

import (
	"net/http"
	"univents/models"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

func (h *Handlers) PatchOccurrence(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	id, err := req.Path("occurrence_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload models.PatchProgramOccurrenceRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	occurrence, err := h.ops.PatchOccurrence(r.Context(), payload.ToInput(id))
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, occurrence)
}
