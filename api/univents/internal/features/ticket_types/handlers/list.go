package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (h *Handlers) List(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	editionID, err := req.Path("edition_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	ticketTypes, err := h.ops.ListByEdition(r.Context(), editionID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, ticketTypes)
}
