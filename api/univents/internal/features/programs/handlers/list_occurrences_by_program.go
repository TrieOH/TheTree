package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (handler *Handlers) ListOccurrencesByProgram(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	programID, err := req.Path("program_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	occurrences, err := handler.queries.ListOccurrencesByProgram(r.Context(), programID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, occurrences)
}
