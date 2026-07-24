package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (handler *Handlers) GetBySlug(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	eventSlug, err := req.Path("event_slug").StringRequired()
	if fun.Bail(w, err) {
		return
	}
	editionSlug, err := req.Path("edition_slug").StringRequired()
	if fun.Bail(w, err) {
		return
	}
	edition, err := handler.queries.GetByEventAndEditionSlug(r.Context(), eventSlug, editionSlug)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, edition)
}
