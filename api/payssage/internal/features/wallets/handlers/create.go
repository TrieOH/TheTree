package handlers

import (
	"net/http"
	"payssage/models"

	"github.com/MintzyG/fun"
)

// Create godoc
// @Summary Create a wallet
// @Description Creates a wallet
// @Tags wallets
// @ID wallets_create
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateWalletRequest true "Wallet creation data"
// @Success 200 {object} models.Wallet "Wallet data"
// @Failure 401 {object} fun.Response "Unauthorized"
// @Failure 404 {object} fun.Response "Bad Request"
// @Failure 500 {object} fun.Response "Internal Server Error"
// @Failure 503 {object} fun.Response "Internal Server Error"
// @Router /wallets [post]
func (h *Handlers) Create(w http.ResponseWriter, r *http.Request) {
	var payload models.CreateWalletRequest
	if fun.BailInto(w, fun.From(r), &payload) {
		return
	}
	wallet, err := h.commands.Create(r.Context(), payload.ToInput())
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, wallet, http.StatusCreated)
}
