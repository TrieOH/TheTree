package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (handler *Handlers) ListByEdition(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	editionID, err := req.Path("edition_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	products, err := handler.queries.ListProductsByEdition(r.Context(), editionID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, products)
}
