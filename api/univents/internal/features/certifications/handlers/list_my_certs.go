package handlers

import (
	"net/http"

	idx "sdk/identityx"

	"github.com/MintzyG/fun"
)

func (h *Handlers) ListMyCerts(w http.ResponseWriter, r *http.Request) {
	ident, err := idx.RequireIdentity(r.Context())
	if fun.Bail(w, err) {
		return
	}
	certs, err := h.ops.ListCertsByUser(r.Context(), ident.Sub.ID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, certs)
}
