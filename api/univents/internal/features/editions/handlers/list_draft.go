package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (handler *Handlers) ListDraft(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	eventID, err := req.Path("event_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	editions, err := handler.queries.ListDraft(r.Context(), eventID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, editions)
}
