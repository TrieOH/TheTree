package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (handler *Handlers) DeleteOccurrence(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	id, err := req.Path("occurrence_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	occurrence, err := handler.commands.DeleteOccurrence(r.Context(), id)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, occurrence)
}
