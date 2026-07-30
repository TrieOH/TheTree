package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (handler *Handlers) GetCertByID(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	id, err := req.Path("cert_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	cert, err := handler.commands.GetCertByID(r.Context(), id)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, cert)
}
