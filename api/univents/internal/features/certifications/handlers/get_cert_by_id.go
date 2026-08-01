package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (h *Handlers) GetCertByID(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	id, err := req.Path("cert_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	cert, err := h.ops.GetCertByID(r.Context(), id)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, cert)
}
