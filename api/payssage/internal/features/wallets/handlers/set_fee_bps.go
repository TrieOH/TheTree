package handlers

import (
	"net/http"
	"payssage/models"

	"github.com/MintzyG/fun"
)

// SetFeeBPS godoc
// @Summary Sets the fee of a wallet
// @Description Sets the fee of a wallet using Basis Point format
// @Tags wallets
// @ID wallets_setfeebps
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.SetFeeBPSRequest true "Wallet fee data"
// @Success 200 {object} fun.Response "OK"
// @Failure 401 {object} fun.Response "Unauthorized"
// @Failure 404 {object} fun.Response "Bad Request"
// @Failure 500 {object} fun.Response "Internal Server Error"
// @Failure 503 {object} fun.Response "Internal Server Error"
// @Router /wallets/{wallet_id}/fee [patch]
func (h *Handlers) SetFeeBPS(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	walletID, err := req.Path("wallet_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload models.SetFeeBPSRequest
	if fun.BailInto(w, fun.From(r), &payload) {
		return
	}
	err = h.ops.SetFeeBPS(r.Context(), payload.ToInput(walletID))
	if fun.Bail(w, err) {
		return
	}
	fun.OK().Send(w)
}
