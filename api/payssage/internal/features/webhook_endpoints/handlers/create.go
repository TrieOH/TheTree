package handlers

import (
	"net/http"
	"payssage/models"

	"github.com/MintzyG/fun"
)

func (h *Handlers) Create(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	walletID, err := req.Path("wallet_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload models.CreateWebhookEndpointRequest
	if fun.BailInto(w, req, &payload) {
		return
	}
	endpoint, err := h.ops.Create(r.Context(), payload.ToInput(walletID))
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, endpoint, http.StatusCreated)
}
