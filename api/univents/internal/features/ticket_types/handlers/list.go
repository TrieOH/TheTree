package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (handler *Handlers) List(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	editionID, err := req.Path("edition_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	ticketTypes, err := handler.queries.ListByEdition(r.Context(), editionID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, ticketTypes)
}
