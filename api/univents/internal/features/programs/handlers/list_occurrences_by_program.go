package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (h *Handlers) ListOccurrencesByProgram(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	programID, err := req.Path("program_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	occurrences, err := h.ops.ListOccurrencesByProgram(r.Context(), programID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, occurrences)
}
