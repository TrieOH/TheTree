package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
	idx "sdk/identityx"
)

func (handler *Handlers) ListMyCerts(w http.ResponseWriter, r *http.Request) {
	ident, err := idx.RequireIdentity(r.Context())
	if fun.Bail(w, err) {
		return
	}
	certs, err := handler.queries.ListCertsByUser(r.Context(), ident.Sub.ID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, certs)
}
