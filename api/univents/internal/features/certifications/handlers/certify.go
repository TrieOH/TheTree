package handlers

import (
	"net/http"
	"univents/models"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

func (h *Handlers) Certify(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	userID, err := req.Path("user_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload models.CertifyRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	cert, err := h.commands.Certify(r.Context(), payload.ToInput(userID))
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, cert)
}
