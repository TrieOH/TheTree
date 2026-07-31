package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (handler *Handlers) ListCertsByEdition(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	editionID, err := req.Path("edition_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	certs, err := handler.queries.ListCertsByEdition(r.Context(), editionID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, certs)
}
