package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (h *Handlers) UnbindCollector(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	walletID, err := req.Path("wallet_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	if fun.Bail(w, h.ops.UnbindCollector(r.Context(), walletID)) {
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
