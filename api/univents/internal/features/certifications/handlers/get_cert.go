package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (h *Handlers) GetCertByID(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	certID, err := req.Path("cert_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	cert, err := h.queries.GetByID(r.Context(), certID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, cert)
}
