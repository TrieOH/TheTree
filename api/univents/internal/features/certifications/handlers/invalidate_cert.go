package handlers

import (
	"net/http"
	"univents/models"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

func (handler *Handlers) InvalidateCert(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	id, err := req.Path("cert_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload models.InvalidCertReason
	if bind.BailInto(w, req, &payload) {
		return
	}
	err = handler.commands.InvalidateCert(r.Context(), id, &payload.Reason)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, nil, http.StatusNoContent)
}
