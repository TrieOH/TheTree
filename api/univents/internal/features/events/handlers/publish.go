package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (h *Handlers) Publish(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	eventID, err := req.Path("event_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	err = h.ops.Publish(r.Context(), eventID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, nil, http.StatusNoContent)
}
