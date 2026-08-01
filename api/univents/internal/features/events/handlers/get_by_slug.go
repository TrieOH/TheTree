package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (h *Handlers) GetBySlug(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	slug, err := req.Path("event_slug").StringRequired()
	if fun.Bail(w, err) {
		return
	}
	event, err := h.ops.GetBySlug(r.Context(), slug)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, event)
}
