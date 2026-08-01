package handlers

import (
	"net/http"
	"payssage/models"

	"github.com/MintzyG/fun"
)

// SetSandboxState godoc
// @Summary Sets the sandbox state of a wallet
// @Tags wallets
// @ID wallets_setsandboxstate
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.SetSandboxRequest true "Sandbox state data"
// @Success 200 {object} fun.Response "OK"
// @Failure 401 {object} fun.Response "Unauthorized"
// @Failure 404 {object} fun.Response "Bad Request"
// @Failure 500 {object} fun.Response "Internal Server Error"
// @Failure 503 {object} fun.Response "Internal Server Error"
// @Router /wallets/{wallet_id}/sandbox [patch]
func (h *Handlers) SetSandboxState(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	walletID, err := req.Path("wallet_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload models.SetSandboxRequest
	if fun.BailInto(w, fun.From(r), &payload) {
		return
	}
	err = h.ops.SetSandbox(r.Context(), payload.ToInput(walletID))
	if fun.Bail(w, err) {
		return
	}
	fun.OK().Send(w)
}
