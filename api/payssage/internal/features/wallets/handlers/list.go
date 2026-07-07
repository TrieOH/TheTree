package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

// List godoc
// @Summary Lists your wallets
// @Tags wallets
// @ID wallets_list
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Wallet "Wallet data"
// @Failure 401 {object} fun.Response "Unauthorized"
// @Failure 404 {object} fun.Response "Bad Request"
// @Failure 500 {object} fun.Response "Internal Server Error"
// @Failure 503 {object} fun.Response "Internal Server Error"
// @Router /wallets [get]
func (h *Handlers) List(w http.ResponseWriter, r *http.Request) {
	wallet, err := h.queries.List(r.Context())
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, wallet, http.StatusCreated)
}
