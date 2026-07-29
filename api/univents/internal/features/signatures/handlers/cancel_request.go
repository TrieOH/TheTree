package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

type cancelRequestPayload struct {
	Reason *string `json:"reason"`
}

func (handler *Handlers) CancelRequest(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	requestID, err := req.Path("request_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload cancelRequestPayload
	if bind.BailInto(w, req, &payload) {
		return
	}
	err = handler.commands.CancelRequest(r.Context(), requestID, payload.Reason)
	if fun.Bail(w, err) {
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
