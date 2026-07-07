package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

// GetByID godoc
// @Summary gets a wallet by ID
// @Tags wallets
// @ID wallets_getbyid
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.Wallet "Wallet data"
// @Failure 401 {object} fun.Response "Unauthorized"
// @Failure 404 {object} fun.Response "Bad Request"
// @Failure 500 {object} fun.Response "Internal Server Error"
// @Failure 503 {object} fun.Response "Internal Server Error"
// @Router /wallets/{wallet_id} [get]
func (h *Handlers) GetByID(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	walletID, err := req.Path("wallet_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	wallet, err := h.queries.GetByID(r.Context(), walletID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, wallet, http.StatusCreated)
}
