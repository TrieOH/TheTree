package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

type denyRequestPayload struct {
	Reason *string `json:"reason"`
}

func (h *Handlers) DenyRequest(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	token, err := req.Query("token").StringRequired()
	if fun.Bail(w, err) {
		return
	}
	var payload denyRequestPayload
	if bind.BailInto(w, req, &payload) {
		return
	}
	err = h.ops.DenyRequest(r.Context(), token, payload.Reason)
	if fun.Bail(w, err) {
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
