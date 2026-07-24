package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (handler *Handlers) GetBySlug(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	slug, err := req.Path("event_slug").StringRequired()
	if fun.Bail(w, err) {
		return
	}
	event, err := handler.queries.GetBySlug(r.Context(), slug)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, event)
}
