package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
	"github.com/google/uuid"
)

type bindCollectorRequest struct {
	CollectorID uuid.UUID `json:"collector_id"`
}

func (h *Handlers) BindCollector(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	walletID, err := req.Path("wallet_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload bindCollectorRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	if fun.Bail(w, h.ops.BindCollector(r.Context(), walletID, payload.CollectorID)) {
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
