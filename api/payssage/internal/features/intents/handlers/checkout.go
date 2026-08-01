package handlers

import (
	"net/http"
	"payssage/models"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

func (h *Handlers) Checkout(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	walletID, err := req.Path("wallet_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload models.CreateIntentRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	intent, err := h.ops.Checkout(r.Context(), payload.ToInput(walletID))
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, intent, http.StatusCreated)
}
