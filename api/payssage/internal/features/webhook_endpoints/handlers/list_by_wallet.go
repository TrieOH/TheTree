package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (h *Handlers) ListByWallet(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	walletID, err := req.Path("wallet_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	endpoints, err := h.queries.ListByWallet(r.Context(), walletID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, endpoints, http.StatusOK)
}
