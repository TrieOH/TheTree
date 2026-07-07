package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

// ListFromOrg godoc
// @Summary lists wallets from an org
// @Tags wallets
// @ID wallets_listfromorg
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Wallet "Wallet data"
// @Failure 401 {object} fun.Response "Unauthorized"
// @Failure 404 {object} fun.Response "Bad Request"
// @Failure 500 {object} fun.Response "Internal Server Error"
// @Failure 503 {object} fun.Response "Internal Server Error"
// @Router /organizations/{organization_id}/wallets [get]
func (h *Handlers) ListFromOrg(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	orgID, err := req.Path("organization_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	wallet, err := h.queries.ListFromOrg(r.Context(), orgID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, wallet, http.StatusCreated)
}
